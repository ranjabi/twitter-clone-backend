package tweet

import (
	"net/http"
	"twitter-clone-backend/models"
	"twitter-clone-backend/user"
)

type Service struct {
	tweetRepository TweetRepository
	userRepository  user.UserRepository
}

func NewService(tweetRepository TweetRepository, userRepository user.UserRepository) Service {
	return Service{tweetRepository: tweetRepository, userRepository: userRepository}
}

func (s *Service) CreateTweet(tweet models.Tweet) (*models.Tweet, error) {
	newTweet, err := s.tweetRepository.CreateTweet(tweet)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to create tweet"}
	}

	err = s.userRepository.DeleteUserRecentTweetsCache(newTweet.UserId)
	if err != nil {
		return nil, err
	}

	return newTweet, nil
}

func (s *Service) UpdateTweet(tweet models.Tweet) (*models.Tweet, error) {
	isTweetExist, err := s.tweetRepository.IsTweetExistById(tweet.Id)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return nil, &models.AppError{Err: err, Message: "Tweet not found", Code: http.StatusNotFound}
	}

	newTweet, err := s.tweetRepository.UpdateTweet(tweet)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to update tweet"}
	}

	return newTweet, nil
}

func (s *Service) DeleteTweet(tweet models.Tweet) error {
	isTweetExist, err := s.tweetRepository.IsTweetExistById(tweet.Id)
	if err != nil {
		return &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return &models.AppError{Err: err, Message: "Tweet not found", Code: http.StatusNotFound}
	}

	if err := s.tweetRepository.DeleteTweet(tweet.Id); err != nil {
		return &models.AppError{Err: err, Message: "Failed to delete tweet"}
	}

	return nil
}

func (s *Service) LikeTweet(userId int, tweetId int) (int, error) {
	isTweetExist, err := s.tweetRepository.IsTweetExistById(tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return 0, &models.AppError{Err: err, Message: "Tweet not found", Code: http.StatusNotFound}
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

	return likeCount, nil
}

func (s *Service) UnlikeTweet(userId int, tweetId int) (int, error) {
	isTweetExist, err := s.tweetRepository.IsTweetExistById(tweetId)
	if err != nil {
		return 0, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return 0, &models.AppError{Err: err, Message: "Tweet not found", Code: http.StatusNotFound}
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

	return likeCount, nil
}
