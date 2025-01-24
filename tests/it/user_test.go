package it

import (
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/usecases/user"

	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // for pgx sql driver
	"github.com/stretchr/testify/assert"
)

func TestFollow_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

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

func TestFollow_AlreadyFollowed(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

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

func TestFollow_FolloweeNotExist(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	err := userService.FollowOtherUser(validUser.Id, notExistUser.Id)
	assert.EqualError(t, err, errmsg.USER_NOT_FOUND)
}

func TestUnfollow_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

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

func TestUnfollow_AlreadyNotFollowed(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	err := userService.UnfollowOtherUser(validUser.Id, validUser2.Id)
	assert.NoError(t, err)
}

func TestUnfollow_FolloweeNotExist(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	err := userService.UnfollowOtherUser(validUser.Id, notExistUser.Id)
	assert.NoError(t, err)
}
