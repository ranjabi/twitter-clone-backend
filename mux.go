package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"
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
		wrappedHandler = AppHandler(h)
	case http.Handler:
		wrappedHandler = h
	default:
		log.Fatal("Unsupported handler type")
	}

	mux.ServeMux.Handle(pattern, wrappedHandler)
}

type AppHandler func(http.ResponseWriter, *http.Request) *models.AppError

// TODO: confirm this by looking at error trace when err is nul at ServiceError <--- The ServeHTTP method called by the appHandler function and displays the returned error
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		if e.Code == 0 {
			e.Code = http.StatusInternalServerError
		}

		fmt.Println(utils.ColorLog(strconv.Itoa(e.Code), utils.RED), utils.ColorLog(http.StatusText(e.Code), utils.RED))
		fmt.Println(utils.ColorLog(e.Error(), utils.RED))

		res, err := json.Marshal(models.ErrorResponse{Message: e.Message})
		if err != nil {
			http.Error(w, errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Code)
		w.Write(res)
	}
}
