package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
    ErrMsgFailedToParseRequestBody = "Failed to parse request body"
	ErrMsgFailedToSerializeResponseBody = "Failed to serialize response body"
	ErrMsgMethodNotAllowed = "Method not allowed"

	RED = "\033[31m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"
	BLUE = "\033[34m"
)

type SuccessResponse[T any] struct {
	Data	T		`json:"data"`
	Message	string	`json:"message"`
}

type SuccessResponseMessage struct {
	Message	string	`json:"message"`
}

type AppMux struct {
	http.ServeMux
	middlewares		[]func(next http.Handler)	http.Handler
}

func (mux *AppMux) RegisterMiddleware(next func(next http.Handler) http.Handler) { // next func(next http.Handler) http.Handler ???
	mux.middlewares = append(mux.middlewares, next)
}

func colorLog(message string, colorCode string) string {
	return fmt.Sprintf("%s%s\033[0m", colorCode, message)
}

func (mux *AppMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(colorLog(r.Method, GREEN), colorLog(r.URL.String(), GREEN))

	requestBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body stream:", err)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	var requestBody map[string]interface{}
	_ = json.Unmarshal(requestBodyBytes, &requestBody)

	prettyRequestBody, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}
	fmt.Println(colorLog("Request body:", BLUE))
	fmt.Printf("%s\n", prettyRequestBody)

	var current http.Handler = &mux.ServeMux // routing logic of http.ServeMux (maps URL paths to handlers)

	// the middleware wraps the current handler. ex: current = loggingMiddleware(authMiddleware(current))
	for _, next := range mux.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}