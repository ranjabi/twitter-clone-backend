package tweet

import (
	"twitter-clone-backend/model"
	"twitter-clone-backend/models"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateTweet(tweet model.Tweet) (*model.Tweet, *models.ServiceError) {
	isTweetExist, err := s.repository.IsTweetExistById(tweet.Id)
	if err != nil {
		return nil, &models.ServiceError{Err: err, Message: "Failed to check tweet"}
	}
	if isTweetExist {
		return nil, &models.ServiceError{Err: err, Message: "Email is already used"}
	}

    newTweet, err := s.repository.CreateTweet(tweet)
    if err != nil {
        return nil, &models.ServiceError{Err: err, Message: "Failed to create tweet"}
    }

    return newTweet, nil
}

func (s Service) UpdateTweet(tweet model.Tweet) (*model.Tweet, *models.ServiceError) {
    newTweet, err := s.repository.UpdateTweet(tweet)
    if err != nil {
        return nil, &models.ServiceError{Err: err, Message: "Failed to update tweet"}
    }

    return newTweet, nil
}
