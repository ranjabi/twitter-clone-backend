package types

import "time"

type User struct {
	Id       int
	Username string
	Email    string
}

type Tweet struct {
	Id        int
	Content   string
	CreatedAt time.Time
	LikeCount int
}

type TweetWithUser struct {
	Tweet
	User User
}
