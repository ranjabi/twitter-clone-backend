package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"twitter-clone-backend/middleware"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"
)

func Login(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			payload := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: 400}
			}

			var isUserExist bool
			query := `SELECT EXISTS (SELECT 1 FROM users WHERE email=@email)`
			args := pgx.NamedArgs{
				"email": payload.Email,
			}
			err := db.QueryRow(ctx, query, args).Scan(&isUserExist)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to check user account", Code: 500}
			}

			if !isUserExist {
				res, err := json.Marshal(models.ErrorResponseMessage{Message: "User not found. Please create an account"})
				if err != nil {
					return &models.AppError{Error: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(res)
				return nil
			}

			query = `SELECT id, username, password FROM users WHERE email=@email`

			var userId string
			var username string
			var hashedPassword string
			err = db.QueryRow(ctx, query, args).Scan(&userId, &username, &hashedPassword)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to get user credential", Code: 500}
			}

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(payload.Password))
			if err != nil {
				return &models.AppError{Error: err, Message: "Email/password is wrong", Code: 500}
			}

			claims := jwt.MapClaims{
				"userId":   userId,
				"username": username,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString([]byte(utils.JWT_SIGNATURE_KEY))
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to sign token", Code: 500}
			}

			type loginResponse struct {
				UserId   string `json:"userId"`
				Username string `json:"username"`
				Token    string `json:"token"`
			}

			userInfo := loginResponse{
				UserId:   userId,
				Username: username,
				Token:    signedToken,
			}

			res, err := json.Marshal(models.SuccessResponse[loginResponse]{Message: "Login success", Data: userInfo})
			if err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: 500}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		default:
			return &models.AppError{Error: nil, Message: utils.ErrMsgMethodNotAllowed, Code: 400}
		}

		return nil
	}
}
