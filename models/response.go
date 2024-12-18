package models

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type SuccessResponseMessage struct {
	Message string `json:"message"`
}

// todo: add Data: nil with NewErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
}

// todo: why not method of ErrorResponse?
// func NewErrorResponse(message string) *ErrorResponse {
// 	return &ErrorResponse{Message: message, Data: nil}
// }
