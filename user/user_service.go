package user

import (
	"context"
	"encoding/json"
	"net/http"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	ctx            context.Context
	userRepository UserRepository
}

func NewService(ctx context.Context, userRepository UserRepository) Service {
	return Service{ctx: ctx, userRepository: userRepository}
}

func (s Service) GetUserById(id int) (*models.User, error) {
	user, err := s.userRepository.GetUserById(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &models.AppError{Err: err, Message: "User not found", Code: http.StatusNotFound}
		}
		return nil, &models.AppError{Err: err, Message: "Failed to get user"}
	}

	return user, nil
}

func (s *Service) GetLastTenTweets(userId int) ([]models.Tweet, error) {
	lastTenTweets, err := s.userRepository.GetLastTenTweets(userId)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to get 10 recent tweets"}
	}

	var tweetsId []int
	for _, tweet := range lastTenTweets {
		tweetsId = append(tweetsId, tweet.Id)
	}

	lastTenTweetsInteractions, err := s.userRepository.GetLastTenTweetsInteractions(userId, tweetsId)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to get 10 recent tweets interactions"}
	}

	for i := range lastTenTweets {
		for j := range lastTenTweetsInteractions {
			if lastTenTweets[i].Id == lastTenTweetsInteractions[j].TweetId {
				lastTenTweets[i].IsLiked = lastTenTweetsInteractions[j].IsLiked
			}
		}
	}

	return lastTenTweets, nil
}

func (s Service) GetUserByUsernameWithRecentTweets(username string) (*models.User, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	/*
		id, username, email, ... -> identified by $.
		recentTweets -> identified by $.recentTweets
	*/
	userCacheStr, err := s.userRepository.GetUserCache(1)
	if err != nil {
		return nil, err
	}
	if userCacheStr != "" {
		// $ ada
		userRecentTweetsCache, err := s.userRepository.GetUserRecentTweetsCache(user.Id)
		if err != nil {
			return nil, err
		}
		var userCache []models.User
		err = json.Unmarshal([]byte(userCacheStr), &userCache)
		if err != nil {
			return nil, err
		}
		user = &userCache[0]

		if userRecentTweetsCache != "[]" {
			// $.recentTweets ada
			utils.CacheLog("HIT UserCache UserRecentTweetsCache")
			return user, nil
		}

		// $.recentTweets gaada
		lastTenTweets, err := s.GetLastTenTweets(user.Id)
		if err != nil {
			return nil, err
		}
		user.RecentTweets = lastTenTweets

		utils.CacheLog("HIT UserRecentTweetsCache")
		_, err = s.userRepository.SetUserRecentTweetsCache(user, lastTenTweets)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	// $ gaada $.recentTweets gaada
	lastTenTweets, err := s.GetLastTenTweets(user.Id)
	if err != nil {
		return nil, err
	}
	user.RecentTweets = lastTenTweets

	utils.CacheLog("MISS UserCache UserRecentTweetsCache")
	_, err = s.userRepository.SetUserCache(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s Service) GetFeed(id int, email string, page int) ([]models.Tweet, error) {
	isUserExist, err := s.userRepository.IsUserExistByEmail(email)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check user account"}
	}
	if !isUserExist {
		return nil, &models.AppError{Err: err, Message: "User not found", Code: http.StatusNotFound}
	}

	feed, err := s.userRepository.GetFeed(id, page)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (s Service) CreateUser(user models.User) (*models.User, error) {
	isUserExist, err := s.userRepository.IsUserExistByEmail(user.Email)
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
	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to create account"}
	}

	return newUser, nil
}

func (s Service) CheckUserCredential(email string, password string) (*models.User, error) {
	isUserExist, err := s.userRepository.IsUserExistByEmail(email)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check user account"}
	}
	if !isUserExist {
		return nil, &models.AppError{Err: err, Message: "User not found. Please create an account", Code: http.StatusNotFound}
	}

	user, err := s.userRepository.GetUserByEmail(email)
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
		"email":    user.Email,
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
	if err := s.userRepository.FollowOtherUser(followerId, followingId); err != nil {
		return &models.AppError{Err: err, Message: "Failed to follow", Code: http.StatusConflict}
	}

	return nil
}

func (s Service) UnfollowOtherUser(followerId int, followingId int) error {
	if err := s.userRepository.UnfollowOtherUser(followerId, followingId); err != nil {
		return &models.AppError{Err: err, Message: "Failed to unfollow", Code: http.StatusConflict}
	}

	return nil
}
