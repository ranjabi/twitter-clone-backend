package it

import (
	"testing"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestTweetCreate_Ok(t *testing.T) {
	err := ResetAndSeed()
	assert.NoError(t, err)

	testTweet := models.Tweet{
		Content: faker.Sentence(),
		UserId:  validUser.Id,
	}
	newTweet, err := tweetService.CreateTweet(testTweet)

	assert.NoError(t, err)
	assert.NotNil(t, newTweet)
	assert.Equal(t, testTweet.Content, newTweet.Content)
	assert.Equal(t, testTweet.UserId, newTweet.UserId)
}

func TestTweetUpdate_Ok(t *testing.T) {
	err := ResetAndSeed()
	assert.NoError(t, err)

	tweet, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tweet)

	tweet.Content = "Updated"
	updatedTweet, err := tweetService.UpdateTweet(*tweet)

	assert.NoError(t, err)
	assert.NotNil(t, updatedTweet)
	assert.Equal(t, "Updated", updatedTweet.Content)
	assert.Equal(t, validTweet.UserId, updatedTweet.UserId)
}

func TestTweetUpdate_NotFound(t *testing.T) {
	err := ResetAndSeed()
	assert.NoError(t, err)

	_, err = tweetService.UpdateTweet(notExistTweet)
	assert.EqualError(t, err, errmsg.TWEET_NOT_FOUND)
}

func TestTweetDelete_Ok(t *testing.T) {
	err := ResetAndSeed()
	assert.NoError(t, err)

	err = tweetService.DeleteTweet(validUser.Id, validTweet.Id)
	assert.NoError(t, err)
}

func TestTweetDelete_NotFound(t *testing.T) {
	err := ResetAndSeed()
	assert.NoError(t, err)

	err = tweetService.DeleteTweet(validUser.Id, notExistTweet.Id)
	assert.EqualError(t, err, errmsg.TWEET_NOT_FOUND)
}

func TestTweetLike_Ok(t *testing.T) {
	// TODO assert isLiked
	err := ResetAndSeed()
	assert.NoError(t, err)

	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.Equal(t, 0, tweetBefore.LikeCount)

	likeCount, err := tweetService.LikeTweet(validUser.Id, validTweet.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, likeCount)

	tweetAfter, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tweetAfter)
	assert.Equal(t, tweetBefore.LikeCount+1, tweetAfter.LikeCount)
}

func TestTweetLike_AlreadyLiked(t *testing.T) {
	err := ResetAndSeed()
	assert.NoError(t, err)

	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tweetBefore)
	assert.Equal(t, 0, tweetBefore.LikeCount)

	likeCount, err := tweetService.LikeTweet(validUser.Id, validTweet.Id)
	assert.NoError(t, err)
	assert.NotZero(t, likeCount)
	assert.Equal(t, 1, likeCount)

	likeCount, err = tweetService.LikeTweet(validUser.Id, validTweet.Id)
	assert.NoError(t, err)
	assert.NotZero(t, likeCount)
	assert.Equal(t, 1, likeCount)

	tweetAfter, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tweetAfter)
	assert.Equal(t, 1, tweetAfter.LikeCount)
}

func TestTweetUnlike_Ok(t *testing.T) {
	// TODO assert isLiked
	err := ResetAndSeed()
	assert.NoError(t, err)

	likeCount, err := tweetRepository.LikeTweet(validUser.Id, validTweet.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, likeCount)

	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tweetBefore)

	likeCount, err = tweetService.UnlikeTweet(validUser.Id, validTweet.Id)
	assert.NoError(t, err)
	assert.Equal(t, 0, likeCount)

	tweetAfter, err := tweetRepository.FindById(validTweet.Id)

	assert.NoError(t, err)
	assert.NotNil(t, tweetAfter)

	assert.Equal(t, tweetBefore.LikeCount-1, tweetAfter.LikeCount)
}

func TestTweetUnlike_AlreadyNotLiked(t *testing.T) {
	err := ResetAndSeed()
	assert.NoError(t, err)

	tweetBefore, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tweetBefore)
	assert.Equal(t, 0, tweetBefore.LikeCount)

	likeCount, err := tweetService.UnlikeTweet(validUser.Id, validTweet.Id)
	assert.NoError(t, err)
	assert.Equal(t, 0, likeCount)

	tweetAfter, err := tweetRepository.FindById(validTweet.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tweetAfter)
	assert.Equal(t, 0, tweetAfter.LikeCount)
}
