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

func (c Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) *models.AppError {
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		return &models.AppError{Error: err, Message: utils.ErrMsgFailedToParseRequestBody, Code: http.StatusInternalServerError}
	}

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return &models.AppError{Error: nil, Message: fmt.Sprintf("Validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()), Code: http.StatusInternalServerError}
		}
	}

	newUser, err := c.service.CreateUser(model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return &models.AppError{Error: err, Message: "Failed to create account", Code: http.StatusInternalServerError}
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
		http.Error(w, utils.ErrMsgFailedToSerializeResponseBody, http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

	return nil
}
