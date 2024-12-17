package tweet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"twitter-clone-backend/model"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
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

func (h Handler) HandleTweetCreate(w http.ResponseWriter, r *http.Request) *models.AppError {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	userId := userInfo["userId"].(string)

	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Content string `json:"content"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: http.StatusInternalServerError}
	}

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &models.AppError{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return &models.AppError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	// not propagate because if db error we can't track it since repo send
	// errors.New() insead of its error
	newTweet, err := h.service.CreateTweet(model.Tweet{
		Content: payload.Content,
		UserId:  userIdInt,
	})

	if serviceErr, ok := err.(*models.ServiceError); ok {
		return &models.AppError{
			Err:   err,
			Message: serviceErr.Message,
			Code:    http.StatusInternalServerError,
		}
	}

	newTweetResponse := struct {
		Id        int       `json:"id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"createdAt"`
		UserId    int       `json:"userId"`
	}{
		Id:        newTweet.Id,
		Content:   newTweet.Content,
		CreatedAt: newTweet.CreatedAt,
		UserId:    newTweet.UserId,
	}
	res, err := json.Marshal(models.SuccessResponse{Message: "Tweet created successfully", Data: newTweetResponse})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleUpdateTweet(w http.ResponseWriter, r *http.Request) *models.AppError {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		TweetId int    `json:"tweetId"`
		Content string `json:"content"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: http.StatusInternalServerError}
	}

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// todo: make variable for fmt.Sprintf("Validation
			return &models.AppError{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	if err != nil {
		return &models.AppError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	// not propagate because if db error we can't track it since repo send
	// errors.New() insead of its error
	newTweet, err := h.service.UpdateTweet(model.Tweet{
		Id:      payload.TweetId,
		Content: payload.Content,
	})
	if serviceErr, ok := err.(*models.ServiceError); ok {
		return &models.AppError{
			Err:   err,
			Message: serviceErr.Message,
			Code:    http.StatusInternalServerError,
		}
	}

	newTweetResponse := struct {
		Id         int       `json:"id"`
		Content    string    `json:"content"`
		CreatedAt  time.Time `json:"createdAt"`
		ModifiedAt time.Time `json:"modifiedAt"`
		UserId     int       `json:"userId"`
	}{
		Id:         newTweet.Id,
		Content:    newTweet.Content,
		CreatedAt:  newTweet.CreatedAt,
		ModifiedAt: newTweet.ModifiedAt,
		UserId:     newTweet.UserId,
	}
	res, err := json.Marshal(models.SuccessResponse{Message: "Tweet updated successfully", Data: newTweetResponse})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleDeleteTweet(w http.ResponseWriter, r *http.Request) *models.AppError {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		TweetId int    `json:"tweetId"`
		Content string `json:"content"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: http.StatusInternalServerError}
	}

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// todo: make variable for fmt.Sprintf("Validation
			return &models.AppError{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	if err != nil {
		return &models.AppError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	newTweet, err := h.service.UpdateTweet(model.Tweet{
		Id:      payload.TweetId,
		Content: payload.Content,
	})
	if serviceErr, ok := err.(*models.ServiceError); ok {
		return &models.AppError{
			Err:   err,
			Message: serviceErr.Message,
			Code:    http.StatusInternalServerError,
		}
	}
	
	newTweetResponse := struct {
		Id         int       `json:"id"`
		Content    string    `json:"content"`
		CreatedAt  time.Time `json:"createdAt"`
		ModifiedAt time.Time `json:"modifiedAt"`
		UserId     int       `json:"userId"`
	}{
		Id:         newTweet.Id,
		Content:    newTweet.Content,
		CreatedAt:  newTweet.CreatedAt,
		ModifiedAt: newTweet.ModifiedAt,
		UserId:     newTweet.UserId,
	}
	res, err := json.Marshal(models.SuccessResponse{Message: "Tweet updated successfully", Data: newTweetResponse})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}
