package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"twitter-clone-backend/constants"
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

func (s *Service) GetRecentTweets(userId int, page int) ([]models.Tweet, error) {
	lastTenTweets, err := s.userRepository.GetRecentTweets(userId, page)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to get recent tweets"}
	}

	return lastTenTweets, nil
}

func (s Service) GetUserByUsernameWithRecentTweets(username string, followerId int, page int) (*models.User, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	isFollowed, err := s.userRepository.IsFollowed(followerId, user.Id)
	if err != nil {
		return nil, err
	}
	user.IsFollowed = isFollowed
	fmt.Println("user:", user)
	/*
		id, username, email, ... -> identified by $.
		recentTweets -> identified by $.recentTweets
	*/
	userCacheStr, err := s.userRepository.GetUserCache(user.Id)
	if err != nil {
		return nil, err
	}

	if userCacheStr != "" {
		// $ ada
		var userCache []models.User
		err = json.Unmarshal([]byte(userCacheStr), &userCache)
		if err != nil {
			return nil, err
		}
		user = &userCache[0]
		userRecentTweetsCacheStr, err := s.userRepository.GetUserRecentTweetsCache(user.Id)
		if err != nil {
			return nil, err
		}

		if page == 1 {
			if userRecentTweetsCacheStr != "" {
				// $.recentTweets ada
				var userRecentTweetsCache [][]models.Tweet
				err = json.Unmarshal([]byte(userRecentTweetsCacheStr), &userRecentTweetsCache)
				if err != nil {
					return nil, err
				}
				user.RecentTweets = userRecentTweetsCache[0]
				utils.CacheLog("HIT UserCache UserRecentTweetsCache")
			} else {
				lastTenTweets, err := s.GetRecentTweets(user.Id, page)
				if err != nil {
					return nil, err
				}
				user.RecentTweets = lastTenTweets
				user.RecentTweetsLength = len(lastTenTweets)
				if len(lastTenTweets) < 10 {
					user.NextPageId = nil
				} else {
					nextPageId := page + 1
					user.NextPageId = &nextPageId
				}

				utils.CacheLog("HIT UserRecentTweetsCache")
				_, err = s.userRepository.SetUserRecentTweetsCache(user, lastTenTweets)
				if err != nil {
					return nil, err
				}
			}
		} else if page > 1 {
			lastTenTweets, err := s.GetRecentTweets(user.Id, page)
			if err != nil {
				return nil, err
			}
			user.RecentTweets = lastTenTweets
			user.RecentTweetsLength = len(lastTenTweets)
			if len(lastTenTweets) < 10 {
				user.NextPageId = nil
			} else {
				nextPageId := page + 1
				user.NextPageId = &nextPageId
			}
		}
	} else {
		// $ gaada $.recentTweets gaada
		lastTenTweets, err := s.GetRecentTweets(user.Id, page)
		if err != nil {
			return nil, err
		}
		user.RecentTweets = lastTenTweets
		user.RecentTweetsLength = len(lastTenTweets)
		if len(lastTenTweets) < 10 {
			user.NextPageId = nil
		} else {
			nextPageId := page + 1
			user.NextPageId = &nextPageId
		}

		if page == 1 {
			utils.CacheLog("MISS UserCache UserRecentTweetsCache")
			userWithoutRecentTweets := user
			temp := user.RecentTweets
			userWithoutRecentTweets.RecentTweets = nil
			_, err = s.userRepository.SetUserCache(userWithoutRecentTweets)
			if err != nil {
				return nil, err
			}
			_, err = s.userRepository.SetUserRecentTweetsCache(user, lastTenTweets)
			if err != nil {
				return nil, err
			}
			user.RecentTweets = temp
		}
	}

	var tweetsId []int
	for _, tweet := range user.RecentTweets {
		tweetsId = append(tweetsId, tweet.Id)
	}

	tweetsInteractions, err := s.userRepository.GetTweetsInteractions(user.Id, tweetsId)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to get recent tweets interactions"}
	}
	for i := range tweetsId {
		for j := range tweetsInteractions {
			if user.RecentTweets[i].Id == tweetsInteractions[j].TweetId {
				user.RecentTweets[i].IsLiked = tweetsInteractions[j].IsLiked
			}
		}
	}

	return user, nil
}

func (s Service) GetFeed(id int, email string, page int) (*models.Feed, error) {
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
		return nil, &models.AppError{Err: err, Message: constants.EMAIL_ALREADY_EXIST_MSG, Code: http.StatusConflict}
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
		return nil, &models.AppError{Err: err, Message: constants.USER_NOT_FOUND_MSG, Code: http.StatusNotFound}
	}

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to get user credential"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, &models.AppError{Err: err, Message: constants.WRONG_CREDENTIAL_MSG}
	}

	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"fullName": user.FullName,
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
