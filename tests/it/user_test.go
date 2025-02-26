package it

import (
	"fmt"
	"net/http"
	"twitter-clone-backend/app"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/types"
)

func (s *TestSuite) TestUserFollow_Ok() {
	// Before
	followerUserBefore, err := s.userRepository.FindById(s.validUser.Id)
	s.NoError(err)
	s.NotNil(followerUserBefore)
	followeeUserBefore, err := s.userRepository.FindById(s.validUser2.Id)
	s.NoError(err)
	s.NotNil(followeeUserBefore)

	// Test
	err = s.userService.FollowOtherUser(s.validUser.Id, s.validUser2.Id)
	s.NoError(err)

	// After
	followerUserAfter, err := s.userRepository.FindById(s.validUser.Id)
	s.NoError(err)
	s.NotNil(followerUserAfter)
	followeeUserAfter, err := s.userRepository.FindById(s.validUser2.Id)
	s.NoError(err)
	s.NotNil(followeeUserAfter)

	s.Equal(followerUserBefore.FollowingCount+1, followerUserAfter.FollowingCount)
	s.Equal(followeeUserBefore.FollowerCount+1, followeeUserAfter.FollowerCount)
}

func (s *TestSuite) TestUserFollow_AlreadyFollowed() {
	// Before
	followerUserBefore, err := s.userRepository.FindById(s.validUser.Id)
	s.NoError(err)
	s.NotNil(followerUserBefore)
	followeeUserBefore, err := s.userRepository.FindById(s.validUser2.Id)
	s.NoError(err)
	s.NotNil(followeeUserBefore)

	// Test
	err = s.userService.FollowOtherUser(s.validUser.Id, s.validUser2.Id)
	s.NoError(err)

	// After
	err = s.userService.FollowOtherUser(s.validUser.Id, s.validUser2.Id)
	s.NoError(err)
}

func (s *TestSuite) TestUserFollow_FolloweeNotExist() {
	// Test
	err := s.userService.FollowOtherUser(s.validUser.Id, s.notExistUser.Id)

	// After
	s.EqualError(err, errmsg.USER_NOT_FOUND)
	s.IsType(&app.Error{}, err)
	s.Equal(http.StatusNotFound, err.(*app.Error).GetCode())
}

func (s *TestSuite) TestUserUnfollow_Ok() {
	// Before
	err := s.userRepository.FollowOtherUser(s.validUser.Id, s.validUser2.Id)
	s.NoError(err)
	followerUserBefore, err := s.userRepository.FindById(s.validUser.Id)
	s.NoError(err)
	s.NotNil(followerUserBefore)
	followeeUserBefore, err := s.userRepository.FindById(s.validUser2.Id)
	s.NoError(err)
	s.NotNil(followeeUserBefore)

	// Test
	err = s.userService.UnfollowOtherUser(s.validUser.Id, s.validUser2.Id)
	s.NoError(err)

	// After
	followerUserAfter, err := s.userRepository.FindById(s.validUser.Id)
	s.NoError(err)
	s.NotNil(followerUserAfter)
	followeeUserAfter, err := s.userRepository.FindById(s.validUser2.Id)
	s.NoError(err)
	s.NotNil(followeeUserAfter)

	s.Equal(followeeUserBefore.FollowerCount-1, followeeUserAfter.FollowerCount)
	s.Equal(followerUserBefore.FollowingCount-1, followerUserAfter.FollowingCount)
}

func (s *TestSuite) TestUserUnfollow_AlreadyNotFollowed() {
	err := s.userService.UnfollowOtherUser(s.validUser.Id, s.validUser2.Id)
	s.NoError(err)
}

func (s *TestSuite) TestUserUnfollow_FolloweeNotExist() {
	err := s.userService.UnfollowOtherUser(s.validUser.Id, s.notExistUser.Id)
	s.NoError(err)
}

func (s *TestSuite) TestUserProfileWithRecentTweetsForFollower_Ok() {
	// validUser follow validUser2
	// validUser2 create 11 tweets
	// validUser see validUser2 profile

	// Before
	tweetNumber := 11
	for i := 0; i < tweetNumber; i++ {
		tweet := types.Tweet{
			User: types.User{
				Id: s.validUser2.Id,
			},
			Content: fmt.Sprintf("Tweet %d content", i+1),
		}
		createdTweet, err := s.tweetRepository.Create(tweet)
		s.NoError(err)
		s.NotNil(createdTweet)
	}

	// Test
	// first page
	profile, err := s.userService.GetProfileByUsernameWithRecentTweetsForFollower(s.validUser2.Username, s.validUser.Id, 1)
	s.NoError(err)
	s.NotNil(profile)

	s.Len(profile.RecentTweets, 10)
	s.Equal(10, profile.RecentTweetsLength)

	for i, tweet := range profile.RecentTweets {
		s.Equal(fmt.Sprintf("Tweet %d content", tweetNumber-i), tweet.Content)
	}

	// second page
	profile, err = s.userService.GetProfileByUsernameWithRecentTweetsForFollower(s.validUser2.Username, s.validUser.Id, 2)
	s.NoError(err)
	s.NotNil(profile)

	s.Len(profile.RecentTweets, 1)
	s.Equal(1, profile.RecentTweetsLength)

	for _, tweet := range profile.RecentTweets {
		s.Equal(fmt.Sprintf("Tweet %d content", 1), tweet.Content)
	}
}
