package models

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type SuccessResponseMessage struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
