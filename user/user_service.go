package user

import (
	"twitter-clone-backend/model"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateUser(user model.User) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		// http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return nil, err
	}

	user.Password = string(hashedPassword)
	newUser, err := s.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
