package models

type AppError struct {
	Error 	error
	Message	string
	Code	int
}