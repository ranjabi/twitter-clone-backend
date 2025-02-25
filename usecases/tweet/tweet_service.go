package tweet

import (
	"net/http"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/types"
	"twitter-clone-backend/usecases/user"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	tweetRepository TweetRepository
	userRepository  user.Repository
}

func NewService(tweetRepository TweetRepository, userRepository user.Repository) Service {
	return Service{tweetRepository: tweetRepository, userRepository: userRepository}
}

func (s *Service) Create(tweet models.Tweet) (*models.Tweet, error) {
	newTweet, err := s.tweetRepository.Create(tweet)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to create tweet"}
	}

	err = s.userRepository.DeleteUserRecentTweetsCache(newTweet.UserId)
	if err != nil {
		return nil, err
	}

	return newTweet, nil
}

func (s *Service) FindById(id int) (*types.TweetWithUser, error) {
	tweet, err := s.tweetRepository.FindByIdV2(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &models.AppError{Err: err, Message: errmsg.TWEET_NOT_FOUND, Code: http.StatusNotFound}
		}
		return nil, err
	}

	return tweet, nil
}

func (s *Service) UpdateTweet(tweet models.Tweet) (*models.Tweet, error) {
	isTweetExist, err := s.tweetRepository.IsTweetExistById(tweet.Id)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return nil, &models.AppError{Err: err, Message: errmsg.TWEET_NOT_FOUND, Code: http.StatusNotFound}
	}

	newTweet, err := s.tweetRepository.UpdateTweet(tweet)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to update tweet"}
	}

	err = s.userRepository.DeleteUserRecentTweetsCache(newTweet.UserId)
	if err != nil {
		// TODO: what is the output if it goes into this?
		return nil, err
	}

	return newTweet, nil
}

func (s *Service) DeleteTweet(userId int, tweetId int) error {
	if err := s.tweetRepository.DeleteTweet(tweetId); err != nil {
		if err == pgx.ErrNoRows {
			return &models.AppError{Err: err, Message: errmsg.TWEET_NOT_FOUND, Code: http.StatusNotFound}
		}
		return &models.AppError{Err: err, Message: "Failed to delete tweet"}
	}

	err := s.userRepository.DeleteUserRecentTweetsCache(userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) LikeTweet(userId int, tweetId int) (int, error) {
	isTweetExist, err := s.tweetRepository.IsTweetExistById(tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return 0, &models.AppError{Err: err, Message: errmsg.TWEET_NOT_FOUND, Code: http.StatusNotFound}
	}

	isTweetLiked, err := s.tweetRepository.IsTweetLiked(userId, tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if isTweetLiked {
		likeCount, err := s.tweetRepository.GetTweetLikeCountById(tweetId)
		if err != nil {
			return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
		}
		return likeCount, nil
	}

	likeCount, err := s.tweetRepository.LikeTweet(userId, tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to like tweet"}
	}
	err = s.userRepository.DeleteUserRecentTweetsCache(userId)
	if err != nil {
		return 0, nil
	}

	return likeCount, nil
}

func (s *Service) UnlikeTweet(userId int, tweetId int) (int, error) {
	isTweetExist, err := s.tweetRepository.IsTweetExistById(tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return 0, &models.AppError{Err: err, Message: errmsg.TWEET_NOT_FOUND, Code: http.StatusNotFound}
	}

	isTweetLiked, err := s.tweetRepository.IsTweetLiked(userId, tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetLiked {
		likeCount, err := s.tweetRepository.GetTweetLikeCountById(tweetId)
		if err != nil {
			return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
		}
		return likeCount, nil
	}

	likeCount, err := s.tweetRepository.UnlikeTweet(userId, tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to unlike tweet"}
	}
	err = s.userRepository.DeleteUserRecentTweetsCache(userId)
	if err != nil {
		return 0, nil
	}

	return likeCount, nil
}
