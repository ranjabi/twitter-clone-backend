package models

type SuccessResponse[T any] struct {
	Message	string	`json:"message"`
	Data	T		`json:"data"`
}

type SuccessResponseMessage struct {
	Message	string	`json:"message"`
}

type ErrorResponse struct {
	Message	string	`json:"message"`
}