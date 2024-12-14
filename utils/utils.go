package utils

import (
	"fmt"
)

const (
    ErrMsgFailedToParseRequestBody = "Failed to parse request body"
	ErrMsgFailedToSerializeResponseBody = "Failed to serialize response body"
	ErrMsgMethodNotAllowed = "Method not allowed"

	RED = "\033[31m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"
	BLUE = "\033[34m"

	JWT_SIGNATURE_KEY = "secret"
)

func ColorLog(message string, colorCode string) string {
	return fmt.Sprintf("%s%s\033[0m", colorCode, message)
}