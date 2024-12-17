package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Follow(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		decoder := json.NewDecoder(r.Body)

		switch r.Method {
		case "POST":
			userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
			userId := userInfo["userId"]

			payload := struct {
				FolloweeId int `json:"followeeId"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Err: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: http.StatusBadRequest}
			}

			query := `INSERT INTO follows (followers_id, following_id) VALUES (@followers_id, @following_id)`
			args := pgx.NamedArgs{
				"followers_id": userId,
				"following_id": payload.FolloweeId,
			}

			_, err := db.Exec(ctx, query, args)
			if err != nil {
				return &models.AppError{Err: err, Message: "Failed to follow", Code: http.StatusInternalServerError}
			}

			res, err := json.Marshal(models.SuccessResponseMessage{Message: "User has been followed"})
			if err != nil {
				return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(res)

		default:
			return &models.AppError{Err: nil, Message: utils.ErrMsgMethodNotAllowed, Code: http.StatusMethodNotAllowed}
		}

		return nil
	}
}

func Unfollow(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		decoder := json.NewDecoder(r.Body)

		switch r.Method {
		case "POST":
			userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
			userId := userInfo["userId"]

			payload := struct {
				FolloweeId int `json:"followeeId"`
			}{}
			if err := decoder.Decode(&payload); err != nil {
				return &models.AppError{Err: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: http.StatusBadRequest}
			}

			query := `DELETE FROM follows WHERE followers_id=@followers_id and following_id=@following_id`
			args := pgx.NamedArgs{
				"followers_id": userId,
				"following_id": payload.FolloweeId,
			}

			_, err := db.Exec(ctx, query, args)
			if err != nil {
				return &models.AppError{Err: err, Message: "Failed to unfollow", Code: http.StatusInternalServerError}
			}

			res, err := json.Marshal(models.SuccessResponseMessage{Message: "User has been unfollowed"})
			if err != nil {
				return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(res)

		default:
			return &models.AppError{Err: nil, Message: utils.ErrMsgMethodNotAllowed, Code: http.StatusMethodNotAllowed}
		}

		return nil
	}
}
