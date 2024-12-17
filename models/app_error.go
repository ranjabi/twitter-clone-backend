package models

// 	type AppError struct {
// 	Err   error = 	type ServiceError struct {
// 						Err   error  	// contain error type from repo or nil, want to preserve repo error to display on log
// 						Message string	<---|	// business logic error
// 					}						|
// 	Message string -------------------------|
// 	Code    int
// 	}

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (e AppError) Error() string {
	// todo: pakai pointer!!!! kenapa??
	if _, ok := e.Err.(*ServiceError); ok {
		// fmt.Println("LOG: e.Err is type of ServiceError")
		return e.Err.Error()
	}

	// fmt.Println("LOG: e.Err is type of built in error")
	return e.Message + ": " + e.Err.Error()
}

type ServiceError struct {
	Err     error
	Message string
}

func (e ServiceError) Error() string {
	// Err bisa nil ketika business logic service benar <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	if e.Err != nil {
		// fmt.Println("LOG: Service error .Err NOT NIL")
		return e.Message + ": " + e.Err.Error()
	}
	// fmt.Println("LOG:  error .Err NIL")
	return e.Message + ": No underlying error"
}
