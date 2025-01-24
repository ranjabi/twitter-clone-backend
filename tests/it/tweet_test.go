package it

import (
	"testing"
	"twitter-clone-backend/models"
	"twitter-clone-backend/usecases/tweet"
	"twitter-clone-backend/usecases/user"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestTweetCreate_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	tweetRepository := tweet.NewRepository(ctx, pgConn, rdConn)
	tweetService := tweet.NewService(tweetRepository, userRepository)

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
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	tweetRepository := tweet.NewRepository(ctx, pgConn, rdConn)
	tweetService := tweet.NewService(tweetRepository, userRepository)

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
