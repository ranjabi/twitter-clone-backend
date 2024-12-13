package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type appError struct {
	Error 	error
	Message	string
	Code	int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

// The ServeHTTP method calls the appHandler function and displays the returned error
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		println(e.Error.Error())
		http.Error(w, e.Message, e.Code)
	}
}

func healthCheckHandler(db *pgxpool.Pool, ctx context.Context) appHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		fmt.Fprintln(w, "Server OK")

		var test string
		err := db.QueryRow(ctx, "select 'OK'").Scan(&test)
		if err != nil {
			fmt.Fprintf(w, "Database NOT OK: %v\n", err)
		} else {
			fmt.Fprintf(w, "Database %v\n", test)
		}

		return nil
	}
}

func registerHandler(db *pgxpool.Pool, ctx context.Context) appHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			payload := struct {
				Username	string	`json:"username"`
				Email		string	`json:"email"`
				Password	string	`json:"password"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &appError{err, "Fail to decode JSON request payload", 500}
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
			if err != nil {
				return &appError{err, "Fail to hash password", 500}
			}

			// TODO: payload validation

			// assume everything is valid, continue to below
			// insert to db
			query := `INSERT INTO users (username, email, password) VALUES (@username, @email, @password) RETURNING username, email`
			args := pgx.NamedArgs{
				"username": payload.Username,
				"email": payload.Email,
				"password": string(hashedPassword),
			}

			type userResponse struct {
				Username	string	`json:"username"`
				Email		string	`json:"email"`
			}

			newUser := userResponse{}

			err = db.QueryRow(ctx, query, args).Scan(&newUser.Username, &newUser.Email)
			if err != nil {
				return &appError{err, "Fail to insert user credential", 500}
			}

			res, err := json.Marshal(newUser)
			if err != nil {
				return &appError{err, "Fail to encode JSON", 500}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(res))
		default:
			return &appError{nil, "Only accept POST request", 400} // will error in serveHTTP if caught
		}

		return nil
	}
}

func loginHandler(db *pgxpool.Pool, ctx context.Context) appHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			payload := struct {
				Email	string `json:"email"`
				Password	string `json:"password"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &appError{err, "Fail to decode JSON request payload", 500}
			}
			fmt.Printf("Decoded payload: %+v\n", payload)

			// TODO: check if user with email exist

			query := `SELECT password FROM users WHERE email=@email`
			args := pgx.NamedArgs{
				"email": payload.Email,
			}

			var hashedPassword string
			err := db.QueryRow(ctx, query, args).Scan(&hashedPassword)
			if err != nil {
				return &appError{err, "Fail to get password", 500}
			}

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
			if err != nil {
				return &appError{err, "Fail to compare password", 500}
			}

			response := SuccessResponseMessage{
				Message: "Login success",
			}
			res, err := json.Marshal(response)
			if err != nil {
				return &appError{err, "Fail to encode JSON", 500}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(res))
		default:
			return &appError{nil, "Only accept POST request", 400}
		}

		return nil
	}
}

func main() {
	ctx := context.Background()

	databaseUrl := "postgres://postgres:123456@localhost:5432/postgres"
	db, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}
	defer db.Close()

	http.Handle("/health-check", appHandler(healthCheckHandler(db, ctx)))
	http.Handle("/register", appHandler(registerHandler(db, ctx)))
	http.Handle("/login", appHandler(loginHandler(db, ctx)))

	fmt.Println("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Error starting server: ", err)
    }
}