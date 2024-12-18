package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
)

var validate *validator.Validate

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: http.StatusBadRequest}
	}

	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &models.AppError{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusBadRequest}
		}
	}

	// karena manggil layer di dalama, maka pakai message dan error dari layer dalam
	newUser, err := h.service.CreateUser(models.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if e, ok := err.(*models.AppError); ok {
		return e
	}

	newUserResponse := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		Username: newUser.Username,
		Email:    newUser.Email,
	}
	res, err := json.Marshal(models.SuccessResponse{Message: "Account created successfully", Data: newUserResponse})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleLoginUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
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

	user, err := h.service.CheckUserCredential(payload.Email, payload.Password)
	if e, ok := err.(*models.AppError); ok {
		return e
	}

	userResponse := struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	}{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Token:    user.Token,
	}
	res, err := json.Marshal(models.SuccessResponse{Message: "Login success", Data: userResponse})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleGetUserProfile(w http.ResponseWriter, r *http.Request) *models.AppError {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToParsePathValue}
	}

	user, err := h.service.GetUserById(id)
	if e, ok := err.(*models.AppError); ok {
		return e
	}

	userProfileResponse := struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		FollowerCount int `json:"followerCount"`
		FollowingCount int `json:"followingCount"`
	}{
		Id: user.Id, 
		Username: user.Username,
		FollowerCount: user.FollowerCount, 
		FollowingCount: user.FollowingCount,
	}

	res, err := json.Marshal(models.SuccessResponse{Data: userProfileResponse})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleFollowOtherUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	userId := userInfo["userId"].(float64)

	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		FollowingId int `json:"followingId" validate:"required"`
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

	err = h.service.FollowOtherUser(int(userId), payload.FollowingId)
	if e, ok := err.(*models.AppError); ok {
		return e
	}

	res, err := json.Marshal(models.SuccessResponseMessage{Message: "User has been followed"})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleUnfollowOtherUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	userId := userInfo["userId"].(float64)

	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		FollowingId int `json:"followingId" validate:"required"`
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

	err = h.service.UnfollowOtherUser(int(userId), payload.FollowingId)
	if e, ok := err.(*models.AppError); ok {
		return e
	}

	res, err := json.Marshal(models.SuccessResponseMessage{Message: "User has been unfollowed"})
	if err != nil {
		return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}
