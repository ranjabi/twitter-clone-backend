package models

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	Id                 int         `json:"id" db:"user_id"`
	Username           string      `json:"username" db:"user_username"`
	FullName           string      `json:"fullName" db:"user_full_name"`
	Email              string      `json:"email" db:"user_email"`
	Password           string      `json:"-" db:"user_password"`
	ProfileImage       pgtype.Text `json:"profileImage" db:"user_profile_image"`
	FollowerCount      int         `json:"followerCount" db:"user_follower_count"`
	FollowingCount     int         `json:"followingCount" db:"user_following_count"`
	RecentTweetsLength int         `json:"tweetsLength" db:"-"`
	RecentTweets       []Tweet     `json:"tweets" db:"-"`
	Token              string      `json:"token" db:"-"`
	IsFollowed         bool        `json:"isFollowed" db:"-"`
	NextPageId         *int        `json:"nextPageId" db:"-"`
}
