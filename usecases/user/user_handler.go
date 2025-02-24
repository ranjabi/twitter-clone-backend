package user

import (
	"encoding/json"
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
