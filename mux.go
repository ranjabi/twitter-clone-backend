package main

import (
	"log"
	"net/http"
	"twitter-clone-backend/handlers"
	"twitter-clone-backend/models"
)

type AppMux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

func (mux *AppMux) RegisterMiddleware(next func(next http.Handler) http.Handler) {
	mux.middlewares = append(mux.middlewares, next)
}

/*
Not in order with struct

	struct {
		Username	string	`json:"username"`
		Email		string	`json:"email"`
		Password	string	`json:"password"`
	}

Request body:

	{
	  "email": "Heaven_Hegmann50@hotmail.com",
	  "password": "example",
	  "username": "Garrick"
	}
*/
func (mux *AppMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var current http.Handler = &mux.ServeMux

	// the middlewares wrap the current handler. ex: current = loggingMiddleware(authMiddleware(current))
	for _, next := range mux.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}

func (mux *AppMux) Handle(pattern string, handler any) {
	var wrappedHandler http.Handler

	switch h := handler.(type) {
	case func(http.ResponseWriter, *http.Request) *models.AppError:
		wrappedHandler = handlers.AppHandler(h)
	case http.Handler:
		wrappedHandler = h
	default:
		log.Fatal("Unsupported handler type")
	}

	mux.ServeMux.Handle(pattern, wrappedHandler)
}
