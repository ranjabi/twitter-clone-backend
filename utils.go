package main

import (
	"fmt"
	"log"
	"net/http"
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
	methodColor := "\033[32m" // green
	urlColor := "\033[32m" // green
	
	log.Println(colorLog(r.Method, methodColor), colorLog(r.URL.String(), urlColor))

	var current http.Handler = &mux.ServeMux // routing logic of http.ServeMux (maps URL paths to handlers)

	// the middleware wraps the current handler. ex: current = loggingMiddleware(authMiddleware(current))
	for _, next := range mux.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}