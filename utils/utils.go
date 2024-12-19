package utils

import (
	"fmt"
	"os"
	"twitter-clone-backend/models"
)

type contextKey string

const (
	ErrMsgFailedToParsePathValue      = "Failed to parse path value"
	ErrMsgFailedToParseRequestBody      = "Failed to parse request body"
	ErrMsgFailedToSerializeResponseBody = "Failed to serialize response body"
	ErrMsgMethodNotAllowed              = "Method not allowed"

	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"

	JWT_SIGNATURE_KEY = "secret"

	UserInfoKey contextKey = "userInfo"
)

func ColorLog(message string, colorCode string) string {
	return fmt.Sprintf("%s%s\033[0m", colorCode, message)
}

func GetDbConnectionUrl(username string, password string, host string, port string, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
}

func GetDbConnectionUrlFromEnv() string {
	// postgres://[user]:[password]@[host]:[port]/[dbname]
	connString := GetDbConnectionUrl(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Println("LOG connString:", connString)

	return connString
}

func PrintV(name string, input any) {
	fmt.Printf("LOG %s:\n%v\n", name, input)
}

func PrintVStruct(name string, input any) {
	fmt.Printf("LOG %s:\n%#v\n", name, input)
}

func HandleErr(err error) *models.AppError {
	if e, ok := err.(*models.AppError); ok {
		return e
	} else {
		fmt.Println("LOG: not AppError instance")
		return &models.AppError{Err: err, Message: err.Error()}
	}
}