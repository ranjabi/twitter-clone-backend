package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"twitter-clone-backend/models"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/utils"
)

func Register(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body) // request body is a stream
			payload := struct {
				Username	string	`json:"username"`
				Email		string	`json:"email"`
				Password	string	`json:"password"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: 400}
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to hash password", Code: 500}
			}

			// TODO: payload validation

			// assume everything is valid, continue to below
			// insert to db
			query := `INSERT INTO users (username, email, password) VALUES (LOWER(@username), LOWER(@email), @password) RETURNING username, email`
			args := pgx.NamedArgs{
				"username": payload.Username,
				"email": payload.Email,
				"password": string(hashedPassword),
			}

			type newUserResponse struct {
				Username	string	`json:"username"`
				Email		string	`json:"email"`
			}

			newUser := newUserResponse{}

			err = db.QueryRow(ctx, query, args).Scan(&newUser.Username, &newUser.Email)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to create new user", Code: 500}
			}

			res, err := json.Marshal(models.SuccessResponse[newUserResponse]{Message: "Account created successfully", Data: newUser}) // write to a string
			if err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: 500}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(res))
		default:
			return &models.AppError{Error: nil, Message: utils.ErrMsgMethodNotAllowed, Code: 400} // will error in serveHTTP if caught
		}

		return nil
	}
}