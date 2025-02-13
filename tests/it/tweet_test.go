package it

import (
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"

	"github.com/go-faker/faker/v4"
)

func (s *TestSuite) TestTweetCreate_Ok() {
	testTweet := models.Tweet{
		Content: faker.Sentence(),
		UserId:  s.validUser.Id,
	}

	newTweet, err := s.tweetService.CreateTweet(testTweet)

	s.NoError(err)
	s.NotNil(newTweet)
	s.Equal(testTweet.Content, newTweet.Content)
	s.Equal(testTweet.UserId, newTweet.UserId)
}

func (s *TestSuite) TestTweetUpdate_Ok() {
	tweet, err := s.tweetRepository.FindById(s.validTweet.Id)
	s.NoError(err)
	s.NotNil(tweet)

	newTweetContent := faker.Sentence()
	tweet.Content = newTweetContent
	updatedTweet, err := s.tweetService.UpdateTweet(*tweet)

	s.NoError(err)
	s.NotNil(updatedTweet)
	s.Equal(newTweetContent, updatedTweet.Content)
}

func (s *TestSuite) TestTweetUpdate_NotFound() {
	_, err := s.tweetService.UpdateTweet(s.notExistTweet)

	s.EqualError(err, errmsg.TWEET_NOT_FOUND)
}

func (s *TestSuite) TestTweetDelete_Ok() {
	err := s.tweetService.DeleteTweet(s.validUser.Id, s.validTweet.Id)

	s.NoError(err)

	// TODO
	// _, err = s.tweetRepository.FindById(s.validTweet.Id)
	// s.EqualError(err, errmsg.TWEET_NOT_FOUND)
}

// func (s *TestSuite) TestTweetDelete_NotFound() {
// 	err := ResetAndSeed()
// 	assert.NoError(t, err)

// 	err = tweetService.DeleteTweet(validUser.Id, notExistTweet.Id)
// 	assert.EqualError(t, err, errmsg.TWEET_NOT_FOUND)
// }

// func (s *TestSuite) TestTweetLike_Ok() {
// 	// TODO assert isLiked
// 	err := ResetAndSeed()
// 	assert.NoError(t, err)

// 	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 0, tweetBefore.LikeCount)

// 	likeCount, err := tweetService.LikeTweet(validUser.Id, validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, likeCount)

// 	tweetAfter, err := tweetRepository.FindById(validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, tweetAfter)
// 	assert.Equal(t, tweetBefore.LikeCount+1, tweetAfter.LikeCount)
// }

// func (s *TestSuite) TestTweetLike_AlreadyLiked() {
// 	err := ResetAndSeed()
// 	assert.NoError(t, err)

// 	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, tweetBefore)
// 	assert.Equal(t, 0, tweetBefore.LikeCount)

// 	likeCount, err := tweetService.LikeTweet(validUser.Id, validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotZero(t, likeCount)
// 	assert.Equal(t, 1, likeCount)

// 	likeCount, err = tweetService.LikeTweet(validUser.Id, validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotZero(t, likeCount)
// 	assert.Equal(t, 1, likeCount)

// 	tweetAfter, err := tweetRepository.FindById(validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, tweetAfter)
// 	assert.Equal(t, 1, tweetAfter.LikeCount)
// }

// func (s *TestSuite) TestTweetUnlike_Ok() {
// 	// TODO assert isLiked
// 	err := ResetAndSeed()
// 	assert.NoError(t, err)

// 	likeCount, err := tweetRepository.LikeTweet(validUser.Id, validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, likeCount)

// 	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, tweetBefore)

// 	likeCount, err = tweetService.UnlikeTweet(validUser.Id, validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 0, likeCount)

// 	tweetAfter, err := tweetRepository.FindById(validTweet.Id)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, tweetAfter)

// 	assert.Equal(t, tweetBefore.LikeCount-1, tweetAfter.LikeCount)
// }

// func (s *TestSuite) TestTweetUnlike_AlreadyNotLiked() {
// 	err := ResetAndSeed()
// 	assert.NoError(t, err)

// 	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, tweetBefore)
// 	assert.Equal(t, 0, tweetBefore.LikeCount)

// 	likeCount, err := tweetService.UnlikeTweet(validUser.Id, validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 0, likeCount)

// 	tweetAfter, err := tweetRepository.FindById(validTweet.Id)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, tweetAfter)
// 	assert.Equal(t, 0, tweetAfter.LikeCount)
// }
