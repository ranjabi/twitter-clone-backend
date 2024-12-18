package models

type User struct {
	Id             int     `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	FollowerCount  int     `json:"followerCount"`
	FollowingCount int     `json:"followingCount"`
	RecentTweets   []Tweet `json:"recentTweets"`
	Password       string  `json:"password"`
	Token          string  `json:"token"`
}
