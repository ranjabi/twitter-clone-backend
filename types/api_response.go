package types

import "time"

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// TODO test
type SuccessResponseMessage struct {
	Message string `json:"message"`
}

// TODO: add Data: nil with NewErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Token    string `json:"token"`
}

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
