package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"twitter-clone-backend/model"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (c Handler) HandleUserRegister(w http.ResponseWriter, r *http.Request) *models.AppError {
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Username string `json:"username" validate:"required"`
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

	// karena manggil layer di dalama, maka pakai message dan error dari layer dalam
	newUser, err := c.service.CreateUser(model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if serviceErr, ok := err.(*models.ServiceError); ok {
		if serviceErr != nil {
			// TODO: BUG, THIS SHOULDN'T NEEDED.
			// err AUTO CASTED TO (*models.ServiceError) even though it's typed nil
			return &models.AppError{
				Err:     err,
				Message: serviceErr.Message,
				Code:    http.StatusInternalServerError,
			}
		}
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

func (c Handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	user, err := c.service.CheckUserCredential(payload.Email, payload.Password)
	// PREVIOUSLY:
	// if err != nil {
	// 	return &models.AppError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	// }
	if serviceErr, ok := err.(*models.ServiceError); ok {
		if serviceErr.Err != nil {
			// TODO: BUG, THIS SHOULDN'T NEEDED.
			// err AUTO CASTED TO (*models.ServiceError) even though it's typed nil
			return &models.AppError{
				Err:     err,
				Message: serviceErr.Message,
				Code:    http.StatusInternalServerError,
			}
		} else {
			// business logic error, set the code
			return &models.AppError{
				Err:     nil,
				Message: serviceErr.Message,
				Code:    http.StatusUnauthorized,
			}
		}
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
