package user

import (
	"net/http"
	"twitter-clone-backend/model"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateUser(user model.User) (*model.User, error) {
	isUserExist, err := s.repository.IsUserExistByEmail(user.Email)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check user account"}
	}
	if isUserExist {
		return nil, &models.AppError{Err: err, Message: "Email is already used", Code: http.StatusConflict}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to hash password"}
	}

	user.Password = string(hashedPassword)
	newUser, err := s.repository.CreateUser(user)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to create account"}
	}

	return newUser, nil
}

func (s Service) CheckUserCredential(email string, password string) (*model.User, error) {
	isUserExist, err := s.repository.IsUserExistByEmail(email)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check user account"}
	}

	if !isUserExist {
		return nil, &models.AppError{Err: err, Message: "User not found. Please create an account", Code: http.StatusNotFound}
	}

	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to get user credential"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Email/password is wrong"}
	}

	claims := jwt.MapClaims{
		"userId":   user.Id,
		"username": user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(utils.JWT_SIGNATURE_KEY))
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to sign token"}
	}

	user.Token = signedToken

	return user, nil
}

func (s Service) FollowOtherUser(followersId int, followingId int) error {
	if err := s.repository.FollowOtherUser(followersId, followingId); err != nil {
		return &models.AppError{Err: err, Message: "Failed to follow", Code: http.StatusConflict}
	}

	return nil
}