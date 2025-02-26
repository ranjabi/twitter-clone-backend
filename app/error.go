package app

import (
	"os"
	"strings"
)

type Error struct {
	Err     error
	Message string
	Code    int
}

func (e *Error) Error() string {
	if strings.Contains(os.Getenv("ENV_NAME"), "test") && e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) GetCode() int {
	return e.Code
}
