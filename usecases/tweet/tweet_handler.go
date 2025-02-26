package tweet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"twitter-clone-backend/app"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/response"
	"twitter-clone-backend/types"
	"twitter-clone-backend/utils"

	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var validate *validator.Validate

type Handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service) Handler {
	validate = validator.New(validator.WithRequiredStructEnabled())
	return Handler{service: service, validate: validate}
}

func (h Handler) HandleCreateTweet(w http.ResponseWriter, r *http.Request) *app.Error {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	userId := userInfo["id"].(float64)

	payload := struct {
		Content string `json:"content" validate:"required"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_REQUEST_BODY, Code: http.StatusInternalServerError}
	}

	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &app.Error{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	newTweet, err := h.service.Create(types.Tweet{
		Content: payload.Content,
		User: types.User{
			Id: int(userId),
		},
	})
	if err != nil {
		return utils.HandleErr(err)
	}

	newTweetResponse := types.Tweet{
		Id:        newTweet.Id,
		Content:   newTweet.Content,
		CreatedAt: newTweet.CreatedAt,
		User:      newTweet.User,
	}

	if err := response.Json(w, http.StatusCreated, response.Payload{Data: newTweetResponse}); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleGetTweet(w http.ResponseWriter, r *http.Request) *app.Error {
	tweetIdStr := r.PathValue("id")
	tweetId, err := strconv.Atoi(tweetIdStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	tweet, err := h.service.FindById(tweetId)
	if err != nil {
		return &app.Error{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	tweetWithUserResponse := types.TweetWithUserResponse{
		Id:        tweet.Id,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt,
		LikeCount: tweet.LikeCount,
		User: types.UserResponse{
			Id:       tweet.User.Id,
			Username: tweet.User.Username,
			Email:    tweet.User.Email,
		},
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Data: tweetWithUserResponse}); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleUpdateTweet(w http.ResponseWriter, r *http.Request) *app.Error {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		TweetId int    `json:"tweetId"`
		Content string `json:"content"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_REQUEST_BODY, Code: http.StatusInternalServerError}
	}

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// TODO: make variable for fmt.Sprintf("Validation
			return &app.Error{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	if err != nil {
		return &app.Error{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	// not propagate because if db error we can't track it since repo send
	// errors.New() insead of its error
	newTweet, err := h.service.UpdateTweet(models.Tweet{
		Id:      payload.TweetId,
		Content: payload.Content,
	})
	if err != nil {
		return utils.HandleErr(err)
	}

	newTweetResponse := struct {
		Id         int              `json:"id"`
		Content    string           `json:"content"`
		CreatedAt  time.Time        `json:"createdAt"`
		ModifiedAt pgtype.Timestamp `json:"modifiedAt"`
		UserId     int              `json:"userId"`
	}{
		Id:         newTweet.Id,
		Content:    newTweet.Content,
		CreatedAt:  newTweet.CreatedAt,
		ModifiedAt: newTweet.ModifiedAt,
		UserId:     newTweet.UserId,
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Data: newTweetResponse}); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleDeleteTweet(w http.ResponseWriter, r *http.Request) *app.Error {
	tweetIdStr := r.PathValue("id")
	tweetId, err := strconv.Atoi(tweetIdStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	userId := userInfo["id"].(float64)

	err = h.service.DeleteTweet(int(userId), tweetId)
	if err != nil {
		return utils.HandleErr(err)
	}

	if err := response.Json(w, http.StatusOK, nil); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleLikeTweet(w http.ResponseWriter, r *http.Request) *app.Error {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	userId := userInfo["id"].(float64)

	likeCount, err := h.service.LikeTweet(int(userId), idInt)
	if err != nil {
		return utils.HandleErr(err)
	}

	likeTweetResponse := struct {
		Id        int `json:"id"`
		LikeCount int `json:"likeCount"`
	}{
		Id:        idInt,
		LikeCount: likeCount,
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Data: likeTweetResponse}); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleUnlikeTweet(w http.ResponseWriter, r *http.Request) *app.Error {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	userId := userInfo["id"].(float64)

	likeCount, err := h.service.UnlikeTweet(int(userId), idInt)
	if err != nil {
		return utils.HandleErr(err)
	}

	likeTweetResponse := struct {
		Id        int `json:"id"`
		LikeCount int `json:"likeCount"`
	}{
		Id:        idInt,
		LikeCount: likeCount,
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Data: likeTweetResponse}); err != nil {
		return err
	}

	return nil
}
