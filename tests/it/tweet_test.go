package it

import (
	"twitter-clone-backend/types"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
)

func (s *TestSuite) TestTweetCreate_Ok() {
	testTweet := types.Tweet{
		Content: faker.Sentence(),
		User: types.User{
			Id: s.validUser.Id,
		},
	}

	newTweet, err := s.tweetService.Create(testTweet)

	s.NoError(err)
	s.NotNil(newTweet)
	s.Equal(testTweet.Content, newTweet.Content)
	s.Equal(testTweet.User.Id, newTweet.User.Id)
}

func (s *TestSuite) TestTweetUpdate_Ok() {
	tweet, err := s.tweetRepository.FindByIdV2(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweet)

	newTweetContent := faker.Sentence()
	tweet.Content = newTweetContent
	updatedTweet, err := s.tweetService.UpdateTweet(*tweet)

	s.NoError(err)
	s.NotNil(updatedTweet)
	s.Equal(newTweetContent, updatedTweet.Content)
}

// TODO
// func (s *TestSuite) TestTweetUpdate_NotFound() {
// 	_, err := s.tweetService.UpdateTweet(s.notExistTweet)
// 	s.EqualError(err, errmsg.TWEET_NOT_FOUND)
// }

func (s *TestSuite) TestTweetDelete_Ok() {
	err := s.tweetService.DeleteTweet(s.validUser.Id, s.validTweet.Id)
	s.NoError(err)

	_, err = s.tweetRepository.FindById(s.validTweet.Id)
	s.EqualError(err, pgx.ErrNoRows.Error())
}

func (s *TestSuite) TestTweetDelete_NotFound() {
	err := s.tweetService.DeleteTweet(s.validUser.Id, s.notExistTweet.Id)
	s.Nil(err)
}

func (s *TestSuite) TestTweetLike_Ok() {
	// Before
	tweetBefore, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.Equal(0, tweetBefore.LikeCount)

	// Test
	likeCount, err := s.tweetService.LikeTweet(s.validUser.Id, s.validTweet.Id)
	s.NoError(err)
	s.Equal(1, likeCount)

	// After
	tweetAfter, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweetAfter)
	s.Equal(tweetBefore.LikeCount+1, tweetAfter.LikeCount)
}

func (s *TestSuite) TestTweetLike_AlreadyLiked() {
	// Before
	tweetBefore, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweetBefore)
	s.Equal(0, tweetBefore.LikeCount)

	likeCount, err := s.tweetService.LikeTweet(s.validUser.Id, s.validTweet.Id)
	s.NoError(err)
	s.Equal(1, likeCount)

	// Test
	likeCount, err = s.tweetService.LikeTweet(s.validUser.Id, s.validTweet.Id)
	s.NoError(err)
	s.Equal(1, likeCount)

	// After
	tweetAfter, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweetAfter)
	s.Equal(1, tweetAfter.LikeCount)
}

func (s *TestSuite) TestTweetUnlike_Ok() {
	// Before
	likeCount, err := s.tweetRepository.LikeTweet(s.validUser.Id, s.validTweet.Id)
	s.NoError(err)
	s.Equal(1, likeCount)

	tweetBefore, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweetBefore)

	// Test
	likeCount, err = s.tweetService.UnlikeTweet(s.validUser.Id, s.validTweet.Id)
	s.NoError(err)
	s.Equal(0, likeCount)

	// After
	tweetAfter, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweetAfter)

	s.Equal(tweetBefore.LikeCount-1, tweetAfter.LikeCount)
}

func (s *TestSuite) TestTweetUnlike_AlreadyNotLiked() {
	// Before
	tweetBefore, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweetBefore)
	s.Equal(0, tweetBefore.LikeCount)

	// Test
	likeCount, err := s.tweetService.UnlikeTweet(s.validUser.Id, s.validTweet.Id)
	s.NoError(err)
	s.Equal(0, likeCount)

	// After
	tweetAfter, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweetAfter)
	s.Equal(0, tweetAfter.LikeCount)
}
