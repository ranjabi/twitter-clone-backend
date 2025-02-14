package auth

import (
	"context"
	"errors"
	"net/http"
	"twitter-clone-backend/config"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/usecases/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	ctx            context.Context
	cfg            *config.Config
	userRepository user.Repository
}

func NewService(ctx context.Context, cfg *config.Config, userRepository user.Repository) Service {
	return Service{ctx, cfg, userRepository}
}

func (s Service) CreateUser(user models.User) (*models.User, error) {
	isUserExist, err := s.userRepository.IsUserExistByEmail(user.Email)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check user account"}
	}
	if isUserExist {
		return nil, &models.AppError{Err: err, Message: errmsg.EMAIL_ALREADY_EXIST, Code: http.StatusConflict}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to hash password"}
	}

	user.Password = string(hashedPassword)
	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to create account"}
	}

	return newUser, nil
}

func (s Service) Login(email string, password string) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.AppError{Err: err, Message: "User is not exist", Code: http.StatusNotFound}
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, &models.AppError{Err: err, Message: errmsg.WRONG_CREDENTIAL, Code: http.StatusUnauthorized}
	}

	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"fullName": user.FullName,
		"email":    user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.JwtSecret))
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to sign token"}
	}

	user.Token = signedToken

	return user, nil
}
