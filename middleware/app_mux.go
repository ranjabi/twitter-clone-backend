package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"twitter-clone-backend/models"

	"twitter-clone-backend/utils"
)

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
)

type AppHandler func(http.ResponseWriter, *http.Request) *models.AppError

// The ServeHTTP method called by the appHandler function and displays the returned error
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		// goes to logging
		fmt.Println(utils.ColorLog(strconv.Itoa(e.Code), RED), utils.ColorLog(http.StatusText(e.Code), RED))
		fmt.Println(utils.ColorLog("Message:", RED), utils.ColorLog(e.Message, RED))

		if e.Error != nil {
			errError := "None"
			errError = e.Error.Error()
			fmt.Println(utils.ColorLog("Error:", RED), errError)
		}

		res, err := json.Marshal(models.ErrorResponse{Message: e.Message})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Code)
		w.Write(res)
	}
}

type AppMux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

func (mux *AppMux) RegisterMiddleware(next func(next http.Handler) http.Handler) { // next func(next http.Handler) http.Handler ???
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
	log.Println(utils.ColorLog(r.Method, GREEN), utils.ColorLog(r.URL.String(), GREEN))

	requestBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body stream:", err)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	var requestBody map[string]any
	_ = json.Unmarshal(requestBodyBytes, &requestBody)

	prettyRequestBody, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}
	fmt.Println(utils.ColorLog("Request body:", BLUE))
	fmt.Printf("%s\n", prettyRequestBody)

	var current http.Handler = &mux.ServeMux // routing logic of http.ServeMux (maps URL paths to handlers)

	// the middleware wraps the current handler. ex: current = loggingMiddleware(authMiddleware(current))
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
	case http.HandlerFunc:
		wrappedHandler = h
	case http.Handler:
		wrappedHandler = h
	default:
		panic("Unsupported handler type") // todo: vs fatal?
	}

	mux.ServeMux.Handle(pattern, wrappedHandler)
}
