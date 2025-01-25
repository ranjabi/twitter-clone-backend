package it

import (
	"fmt"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/usecases/tweet"
	"twitter-clone-backend/usecases/user"

	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // for pgx sql driver
	"github.com/stretchr/testify/assert"
)

func TestUserFollow_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	followingUserBefore, err := userService.GetUserById(validUser.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followingUserBefore)
	followedUserBefore, err := userService.GetUserById(validUser2.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followedUserBefore)

	err = userService.FollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)

	followingUserAfter, err := userService.GetUserById(validUser.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followingUserAfter)
	followedUserAfter, err := userService.GetUserById(validUser2.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followedUserAfter)

	assert.Equal(t, followedUserAfter.FollowerCount, followedUserBefore.FollowerCount+1)
	assert.Equal(t, followingUserAfter.FollowingCount, followingUserBefore.FollowingCount+1)

	err = userService.UnfollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)
}

func TestUserFollow_AlreadyFollowed(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	followingUserBefore, err := userService.GetUserById(validUser.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followingUserBefore)
	followedUserBefore, err := userService.GetUserById(validUser2.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followedUserBefore)

	err = userService.FollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)
	err = userService.FollowOtherUser(validUser.Id, validUser2.Id)
	assert.EqualError(t, err, errmsg.ALREADY_FOLLOWED)

	err = userService.UnfollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)
}

func TestUserFollow_FolloweeNotExist(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	err := userService.FollowOtherUser(validUser.Id, notExistUser.Id)
	assert.EqualError(t, err, errmsg.USER_NOT_FOUND)
}

func TestUserUnfollow_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	err := userService.FollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)

	followingUserBefore, err := userService.GetUserById(validUser.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followingUserBefore)
	followedUserBefore, err := userService.GetUserById(validUser2.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followedUserBefore)

	err = userService.UnfollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)

	followingUserAfter, err := userService.GetUserById(validUser.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followingUserAfter)
	followedUserAfter, err := userService.GetUserById(validUser2.Id)
	assert.NoError(t, err)
	assert.NotNil(t, followedUserAfter)

	assert.Equal(t, followedUserAfter.FollowerCount, followedUserBefore.FollowerCount-1)
	assert.Equal(t, followingUserAfter.FollowingCount, followingUserBefore.FollowingCount-1)
}

func TestUserUnfollow_AlreadyNotFollowed(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	err := userService.UnfollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)
}

func TestUserUnfollow_FolloweeNotExist(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	err := userService.UnfollowOtherUser(validUser.Id, notExistUser.Id)
	assert.NoError(t, err)
}

func TestUserProfileWithRecentTweetsForFollower_Ok(t *testing.T) {
	/*
		validUser follow validUser2
		validUser2 create 11 tweets
		TODO validUser like tweet-1 (for interaction)
		validUser see validUser2 profile
	*/

	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	tweetRepository := tweet.NewRepository(ctx, pgConn, rdConn)

	userService := user.NewService(ctx, cfg, userRepository)

	tweetNumber := 11
	for i := 0; i < tweetNumber; i++ {
		tweet := models.Tweet{
			UserId:  validUser2.Id,
			Content: fmt.Sprintf("Tweet %d", i+1),
		}
		createdTweet, err := tweetRepository.CreateTweet(tweet)
		assert.NoError(t, err)
		assert.NotNil(t, createdTweet)
	}

	profile, err := userService.GetProfileByUsernameWithRecentTweetsForFollower(validUser2.Username, validUser.Id, 1)
	assert.NoError(t, err)
	assert.NotNil(t, profile)

	// first page
	assert.Len(t, profile.RecentTweets, 10)
	assert.Equal(t, 10, profile.RecentTweetsLength)

	for i, tweet := range profile.RecentTweets {
		assert.Equal(t, fmt.Sprintf("Tweet %d", tweetNumber-i), tweet.Content)
	}

	// second page
	profile, err = userService.GetProfileByUsernameWithRecentTweetsForFollower(validUser2.Username, validUser.Id, 2)
	assert.NoError(t, err)
	assert.NotNil(t, profile)

	assert.Len(t, profile.RecentTweets, 1)
	assert.Equal(t, 1, profile.RecentTweetsLength)

	for _, tweet := range profile.RecentTweets {
		assert.Equal(t, fmt.Sprintf("Tweet %d", 1), tweet.Content)
	}
}
