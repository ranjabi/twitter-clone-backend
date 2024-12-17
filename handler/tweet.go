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
		decoder := json.NewDecoder(r.Body)
		switch r.Method {
		case "POST":
			userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
			userId := userInfo["userId"]

			payload := struct {
				Content string `json:"content"`
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
				Id        string    `json:"tweetId"`
				Content   string    `json:"content"`
				CreatedAt time.Time `json:"createdAt"`
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
			w.Write(res)

		case "PUT":
			payload := struct {
				TweetId int    `json:"tweetId"`
				Content string `json:"content"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: 400}
			}

			var isTweetExist bool
			query := `SELECT EXISTS (SELECT 1 FROM tweets WHERE id=@id)`
			args := pgx.NamedArgs{
				"id": payload.TweetId,
			}
			err := db.QueryRow(ctx, query, args).Scan(&isTweetExist)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to check user", Code: 500}
			}

			if !isTweetExist {
				res, err := json.Marshal(models.ErrorResponseMessage{Message: "Tweet not found"})
				if err != nil {
					return &models.AppError{Error: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: 404}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write(res)

				return nil
			}

			query = `UPDATE tweets SET content=@content, modified_at=@modifiedAt WHERE id=@tweetId RETURNING id, content, modified_at`
			args = pgx.NamedArgs{
				"tweetId":    payload.TweetId,
				"content":    payload.Content,
				"modifiedAt": time.Now(),
			}

			type updatedTweetResponse struct {
				Id         string    `json:"tweetId"`
				Content    string    `json:"content"`
				ModifiedAt time.Time `json:"modifiedAt"`
			}
			updatedTweet := updatedTweetResponse{}

			err = db.QueryRow(ctx, query, args).Scan(&updatedTweet.Id, &updatedTweet.Content, &updatedTweet.ModifiedAt)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to update tweet", Code: 500}
			}

			res, err := json.Marshal(models.SuccessResponse[updatedTweetResponse]{Message: "Tweet updated successfully", Data: updatedTweet})
			if err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: 500}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(res)

		case "DELETE":
			payload := struct {
				TweetId int `json:"tweetId"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: 400}
			}

			var isTweetExist bool
			query := `SELECT EXISTS (SELECT 1 FROM tweets WHERE id=@id)`
			args := pgx.NamedArgs{
				"id": payload.TweetId,
			}
			err := db.QueryRow(ctx, query, args).Scan(&isTweetExist)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to check user", Code: 500}
			}

			if !isTweetExist {
				res, err := json.Marshal(models.ErrorResponseMessage{Message: "Tweet not found"})
				if err != nil {
					return &models.AppError{Error: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusNotFound}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write(res)
				return nil
			}

			query = `DELETE FROM tweets WHERE id=@id`
			args = pgx.NamedArgs{
				"id": payload.TweetId,
			}

			_, err = db.Exec(ctx, query, args)
			if err != nil {
				return &models.AppError{Error: err, Message: "Failed to delete tweet", Code: 500}
			}

			res, err := json.Marshal(models.SuccessResponseMessage{Message: "Tweet deleted successfully"})
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
