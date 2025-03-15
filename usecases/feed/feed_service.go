package feed

import (
	"context"
	"fmt"

	"github.com/ranjabi/twitter-clone-backend/types"
)

type Service struct {
	ctx            context.Context
	feedRepository Repository
}

func NewService(ctx context.Context, feedRepository Repository) Service {
	return Service{
		ctx,
		feedRepository,
	}
}

func (s *Service) GetFeed(userId int, page int) (*types.Feed, error) {
	// x check if key exist in cache
	// TODO if not exist, query from db []<tweetId, userId> and set cache
	// populate with necessary tweet and user data

	isCacheExist, err := s.feedRepository.IsCacheExist(userId)
	if err != nil {
		return nil, err
	}

	if isCacheExist {
		feed, err := s.feedRepository.GetFeedFromCache(userId, page)
		if err != nil {
			return nil, err
		}
		return feed, nil
	}
	fmt.Println("Feed not exist in cache")

	return nil, nil
}
