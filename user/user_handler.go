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

func (c Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, utils.ErrMsgFailedToParseRequestBody, http.StatusInternalServerError)
		return
	}

	err := validate.Struct(payload)
	if err != nil {
		fmt.Println("Validation error start:")
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Validation for '%s' failed on the '%s' tag\n", err.Field(), err.Tag())
		}
		fmt.Println("Validation error end:")

		return
	}

	newUser, err := c.service.CreateUser(model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	type NewUserResponse struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	newUserResponse := NewUserResponse{
		Username: newUser.Username,
		Email:    newUser.Email,
	}

	res, err := json.Marshal(models.SuccessResponse[NewUserResponse]{Message: "Account created successfully", Data: newUserResponse})
	if err != nil {
		http.Error(w, utils.ErrMsgFailedToSerializeResponseBody, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
