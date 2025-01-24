package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"twitter-clone-backend/models"
)

type contextKey string

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"

	UserInfoKey contextKey = "userInfo"
)

func ColorLog(message string, colorCode string) string {
	return fmt.Sprintf("%s%s\033[0m", colorCode, message)
}

func HandleErr(err error) *models.AppError {
	// TODO kelemahan: cuman bisa nampilin track sampai handler
	if e, ok := err.(*models.AppError); ok {
		return e
	} else {
		HandleErrLog(err.Error())
		return &models.AppError{Err: err, Message: err.Error()}
	}
}

func CacheLog(content any) {
	// TODO Get the caller information (skip 1 frame to get the caller of CacheLog)
	_, fullFilePath, line, ok := runtime.Caller(1)
	if !ok {
		fullFilePath = "unknown_file"
		line = 0
	}
	file := filepath.Base(fullFilePath)

	appLog := log.New(os.Stdout, "cache: ", 0)
	appLog.Printf("%s:%d: %v\n", file, line, content)
}

func HandleErrLog(content any) {
	_, fullFilePath, line, ok := runtime.Caller(1)
	if !ok {
		fullFilePath = "unknown_file"
		line = 0
	}
	file := filepath.Base(fullFilePath)

	appLog := log.New(os.Stdout, "HandleErr: ", 0)
	appLog.Printf("%s:%d: %v\n", file, line, content)
}
