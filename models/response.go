package models

type SuccessResponse[T any] struct {
	Message	string	`json:"message"`
	Data	T		`json:"data"`
}

type SuccessResponseMessage struct {
	Message	string	`json:"message"`
}

type SuccessResponseData struct {
	Data	any	`json:"data"`
}

type ErrorResponseMessage struct {
	Message	string	`json:"message"`
}