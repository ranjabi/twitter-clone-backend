package models

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
