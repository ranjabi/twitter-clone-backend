package user

import (
	"encoding/json"
	"net/http"
	"twitter-clone-backend/model"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (c Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, utils.ErrMsgFailedToParseRequestBody, http.StatusInternalServerError)
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
