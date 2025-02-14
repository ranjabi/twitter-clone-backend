package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(user models.User) (*models.User, error)
	Login(email string, password string) (*models.User, error)
}

type Handler struct {
	userService Service
	authService AuthService
	validate    *validator.Validate
}

func NewHandler(userService Service, authService AuthService, validate *validator.Validate) Handler {
	return Handler{userService, authService, validate}
}

func (h Handler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		FullName string `json:"fullName" validate:"required"`
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_PARSE_REQUEST_BODY, Code: http.StatusBadRequest}
	}

	if err := h.validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &models.AppError{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusBadRequest}
		}
	}

	// karena manggil layer di dalam, maka pakai message dan error dari layer dalam
	newUser, err := h.authService.Register(models.User{
		FullName: payload.FullName,
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return utils.HandleErr(err)
	}

	res, err := json.Marshal(models.SuccessResponse{Message: "Account created successfully. Please login", Data: newUser})
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleLoginUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_PARSE_REQUEST_BODY, Code: http.StatusInternalServerError}
	}

	err := h.validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &models.AppError{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		return utils.HandleErr(err)
	}

	res, err := json.Marshal(models.SuccessResponse{Message: "Login success", Data: user})
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) *models.AppError {
	username := r.PathValue("username")
	queryParams := r.URL.Query()
	page := queryParams.Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	followerId := userInfo["id"].(float64)
	user, err := h.userService.GetProfileByUsernameWithRecentTweetsForFollower(username, int(followerId), pageInt)
	if err != nil {
		return utils.HandleErr(err)
	}

	res, err := json.Marshal(models.SuccessResponse{Data: user})
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleFollowOtherUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	followerId := userInfo["id"].(float64)

	followingIdStr := r.PathValue("id")
	followingId, err := strconv.Atoi(followingIdStr)
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	err = h.userService.FollowOtherUser(int(followerId), followingId)
	if err != nil {
		return utils.HandleErr(err)
	}

	res, err := json.Marshal(models.SuccessResponseMessage{Message: "User has been followed"})
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleUnfollowOtherUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	followerId := userInfo["id"].(float64)

	followingIdStr := r.PathValue("id")
	followingId, err := strconv.Atoi(followingIdStr)
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	err = h.userService.UnfollowOtherUser(int(followerId), followingId)
	if err != nil {
		return utils.HandleErr(err)
	}

	res, err := json.Marshal(models.SuccessResponseMessage{Message: "User has been unfollowed"})
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}

func (h Handler) HandleGetFeed(w http.ResponseWriter, r *http.Request) *models.AppError {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	email := userInfo["email"].(string)

	query := r.URL.Query()
	pageStr := query.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	feed, err := h.userService.GetFeed(id, email, page)
	if err != nil {
		return utils.HandleErr(err)
	}

	res, err := json.Marshal(models.SuccessResponse{Data: feed})
	if err != nil {
		return &models.AppError{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}
