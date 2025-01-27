package models

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) GetCode() int {
	return e.Code
}
