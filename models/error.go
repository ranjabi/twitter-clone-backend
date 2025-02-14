package models

import (
	"os"
	"strings"
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (e *AppError) Error() string {
	if strings.Contains(os.Getenv("ENV_NAME"), "test") && e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) GetCode() int {
	return e.Code
}
