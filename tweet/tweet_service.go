package tweet

import (
	"net/http"
	"twitter-clone-backend/models"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateTweet(tweet models.Tweet) (*models.Tweet, error) {
	newTweet, err := s.repository.CreateTweet(tweet)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to create tweet"}
	}

	return newTweet, nil
}

func (s Service) UpdateTweet(tweet models.Tweet) (*models.Tweet, error) {
	isTweetExist, err := s.repository.IsTweetExistById(tweet.Id)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return nil, &models.AppError{Err: err, Message: "Tweet not found", Code: http.StatusNotFound}
	}

	newTweet, err := s.repository.UpdateTweet(tweet)
	if err != nil {
		return nil, &models.AppError{Err: err, Message: "Failed to update tweet"}
	}

	return newTweet, nil
}

func (s Service) DeleteTweet(tweet models.Tweet) error {
	isTweetExist, err := s.repository.IsTweetExistById(tweet.Id)
	if err != nil {
		return &models.AppError{Err: err, Message: "Failed to check tweet"}
	}
	if !isTweetExist {
		return &models.AppError{Err: err, Message: "Tweet not found", Code: http.StatusNotFound}
	}

	if err := s.repository.DeleteTweet(tweet.Id); err != nil {
		return &models.AppError{Err: err, Message: "Failed to delete tweet"}
	}

	return nil
}
