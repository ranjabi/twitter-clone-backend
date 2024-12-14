package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"
)

func Tweet(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		switch r.Method {
		case "POST":
			userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
			userId := userInfo["userId"]

			decoder := json.NewDecoder(r.Body)
			payload := struct {
				Content	string	`json:"content"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: 400}
			}
		
			query := `INSERT INTO tweets (content, user_id)  VALUES (@content, @user_id) RETURNING id, content, created_at`
			args := pgx.NamedArgs{
				"content": payload.Content,
				"user_id": userId,
			}
	
			type newTweetResponse struct {
				Id			string		`json:"id"`
				Content		string		`json:"content"`
				CreatedAt	time.Time	`json:"createdAt"`
			}
			newTweet := newTweetResponse{}
	
			err := db.QueryRow(ctx, query, args).Scan(&newTweet.Id, &newTweet.Content, &newTweet.CreatedAt)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to create tweet", Code: 500}
			}
	
			res, err := json.Marshal(models.SuccessResponse[newTweetResponse]{Message: "Tweet created successfully", Data: newTweet})
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