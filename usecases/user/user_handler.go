package user

import (
	"net/http"
	"strconv"
	"twitter-clone-backend/app"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/response"
	"twitter-clone-backend/utils"

	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	userService Service
	validate    *validator.Validate
}

func NewHandler(userService Service, validate *validator.Validate) Handler {
	return Handler{userService, validate}
}

func (h Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) *app.Error {
	username := r.PathValue("username")
	queryParams := r.URL.Query()
	page := queryParams.Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	followerId := userInfo["id"].(float64)
	user, err := h.userService.GetProfileByUsernameWithRecentTweetsForFollower(username, int(followerId), pageInt)
	if err != nil {
		return utils.HandleErr(err)
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Data: user}); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleFollowOtherUser(w http.ResponseWriter, r *http.Request) *app.Error {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	followerId := userInfo["id"].(float64)

	followingIdStr := r.PathValue("id")
	followingId, err := strconv.Atoi(followingIdStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	err = h.userService.FollowOtherUser(int(followerId), followingId)
	if err != nil {
		return utils.HandleErr(err)
	}

	if err := response.Json(w, http.StatusOK, nil); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleUnfollowOtherUser(w http.ResponseWriter, r *http.Request) *app.Error {
	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	followerId := userInfo["id"].(float64)

	followingIdStr := r.PathValue("id")
	followingId, err := strconv.Atoi(followingIdStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	err = h.userService.UnfollowOtherUser(int(followerId), followingId)
	if err != nil {
		return utils.HandleErr(err)
	}

	if err := response.Json(w, http.StatusOK, nil); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleGetFeed(w http.ResponseWriter, r *http.Request) *app.Error {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	userInfo := r.Context().Value(utils.UserInfoKey).(jwt.MapClaims)
	email := userInfo["email"].(string)

	query := r.URL.Query()
	pageStr := query.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_PATH_VALUE}
	}

	feed, err := h.userService.GetFeed(id, email, page)
	if err != nil {
		return utils.HandleErr(err)
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Data: feed}); err != nil {
		return err
	}

	return nil
}
