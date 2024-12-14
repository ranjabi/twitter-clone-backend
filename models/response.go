package models

type SuccessResponse[T any] struct {
	Data	T		`json:"data"`
	Message	string	`json:"message"`
}

type SuccessResponseMessage struct {
	Message	string	`json:"message"`
}