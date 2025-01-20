package models

type User struct {
	Id                 int     `json:"id"`
	Username           string  `json:"username"`
	FullName           string  `json:"fullName"`
	Email              string  `json:"email"`
	ProfileImage       string  `json:"profileImage"`
	FollowerCount      int     `json:"followerCount"`
	FollowingCount     int     `json:"followingCount"`
	RecentTweetsLength int     `json:"tweetsLength"`
	RecentTweets       []Tweet `json:"tweets"`
	Password           string  `json:"password"`
	Token              string  `json:"token"`
	IsFollowed         bool    `json:"isFollowed"`
	NextPageId         *int    `json:"nextPageId"`
}
