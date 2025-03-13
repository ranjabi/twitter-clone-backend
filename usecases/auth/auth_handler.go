package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ranjabi/twitter-clone-backend/app"
	"github.com/ranjabi/twitter-clone-backend/errmsg"
	"github.com/ranjabi/twitter-clone-backend/response"
	"github.com/ranjabi/twitter-clone-backend/types"
	"github.com/ranjabi/twitter-clone-backend/utils"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	authService Service
	validate    *validator.Validate
}

func NewHandler(authService Service, validate *validator.Validate) Handler {
	return Handler{authService, validate}
}

func (h Handler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) *app.Error {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		FullName string `json:"fullName" validate:"required"`
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_REQUEST_BODY, Code: http.StatusBadRequest}
	}

	if err := h.validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &app.Error{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusBadRequest}
		}
	}

	// karena manggil layer di dalam, maka pakai message dan error dari layer dalam
	_, err := h.authService.Register(types.User{
		FullName: payload.FullName,
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return utils.HandleErr(err)
	}

	if err := response.Json(w, http.StatusCreated, response.Payload{Message: "Account created successfully. Please login"}); err != nil {
		return err
	}

	return nil
}

func (h Handler) HandleLoginUser(w http.ResponseWriter, r *http.Request) *app.Error {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_PARSE_REQUEST_BODY, Code: http.StatusInternalServerError}
	}

	err := h.validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &app.Error{Err: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		return utils.HandleErr(err)
	}

	userResponse := types.LoginResponse{
		Id:       user.Id,
		Username: user.Username,
		FullName: user.FullName,
		Token:    user.Token,
	}

	if err := response.Json(w, http.StatusOK, response.Payload{Message: "Login success", Data: userResponse}); err != nil {
		return err
	}

	return nil
}
