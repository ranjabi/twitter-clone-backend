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

func Login(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			payload := struct {
				Email	string `json:"email"`
				Password	string `json:"password"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: 400}
			}

			// TODO: check if user with email exist

			query := `SELECT password FROM users WHERE email=@email`
			args := pgx.NamedArgs{
				"email": payload.Email,
			}

			var hashedPassword string
			err := db.QueryRow(ctx, query, args).Scan(&hashedPassword)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to get password", Code: 500}
			}

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
			if err != nil {
				return &models.AppError{Error: err, Message: "Email/password is wrong", Code: 500}
			}

			response := models.SuccessResponseMessage{
				Message: "Login success",
			}
			res, err := json.Marshal(response)
			if err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: 500}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(res))
		default:
			return &models.AppError{Error: nil, Message: utils.ErrMsgMethodNotAllowed, Code: 400}
		}

		return nil
	}
}