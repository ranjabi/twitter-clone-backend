package types

import "time"

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type TweetWithUserResponse struct {
	Id        int          `json:"id"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
	LikeCount int          `json:"likeCount"`
	User      UserResponse `json:"user"`
}
