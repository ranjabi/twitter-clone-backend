package user

import (
	"net/http"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) GetUserById(id int) (*models.User, error) {
	user, err := s.repository.GetUserById(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &models.AppError{Err: err, Message: "User not found", Code: http.StatusNotFound}
		}
		return nil, &models.AppError{Err: err, Message: "Failed to get user"}
	}

	return user, nil
}

func (s *Service) GetLastTenTweets(userId int) ([]models.Tweet, error) {
	lastTenTweets, err := s.repository.GetLastTenTweets(userId)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to get 10 recent tweets"}
	}

	return lastTenTweets, nil
}

func (s Service) GetUserByIdWithRecentTweets(id int) (*models.User, error) {
	user, err := s.GetUserById(id)
	if err != nil {
		return nil, err
	}
	lastTenTweets, err := s.GetLastTenTweets(id)
	if err != nil {
		return nil, err
	}
	user.RecentTweets = lastTenTweets

	return user, nil
}

func (s Service) CreateUser(user models.User) (*models.User, error) {
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

func (s Service) CheckUserCredential(email string, password string) (*models.User, error) {
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

func (s Service) FollowOtherUser(followerId int, followingId int) error {
	if err := s.repository.FollowOtherUser(followerId, followingId); err != nil {
		return &models.AppError{Err: err, Message: "Failed to follow", Code: http.StatusConflict}
	}

	return nil
}

func (s Service) UnfollowOtherUser(followerId int, followingId int) error {
	if err := s.repository.UnfollowOtherUser(followerId, followingId); err != nil {
		return &models.AppError{Err: err, Message: "Failed to unfollow", Code: http.StatusConflict}
	}

	return nil
}
