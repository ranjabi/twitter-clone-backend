package handler

import (
	"context"
	"encoding/json"
	"net/http"
	
	"golang.org/x/crypto/bcrypt"
	
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	jwt "github.com/golang-jwt/jwt/v5"
	
	"twitter-clone-backend/models"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/utils"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("secret")

type appClaims struct {
	jwt.RegisteredClaims
	userId		string
	username	string
}

func Login(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			payload := struct {
				Email		string	`json:"email"`
				Password	string	`json:"password"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: 400}
			}

			// TODO: check if user with email exist

			query := `SELECT id, username, password FROM users WHERE email=@email`
			args := pgx.NamedArgs{
				"email": payload.Email,
			}

			var userId string
			var username string
			var hashedPassword string
			err := db.QueryRow(ctx, query, args).Scan(&userId, &username, &hashedPassword)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to get password", Code: 500}
			}

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
			if err != nil {
				return &models.AppError{Error: err, Message: "Email/password is wrong", Code: 500}
			}

			claims := appClaims{
				userId: userId,
				username: username,
			}
			token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)
			signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to sign token", Code: 500}
			}

			type loginResponse struct {
				UserId   	string	`json:"userId"`
    			Username	string	`json:"username"`
				Token		string	`json:"token"`
			}

			resData := loginResponse{
				UserId: userId,
				Username: username,
				Token: signedToken,
			}

			res, err := json.Marshal(models.SuccessResponse[loginResponse]{Message: "Login success", Data: resData})
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