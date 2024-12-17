package tweet

import (
    "twitter-clone-backend/model"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateTweet(tweet model.Tweet) (*model.Tweet, error) {
    newTweet, err := s.repository.CreateTweet(tweet)
    if err != nil {
        return nil, err
    }

    return newTweet, nil
}
