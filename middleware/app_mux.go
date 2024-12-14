package middleware

import (
	"fmt"
	"net/http"
	"twitter-clone-backend/models"
	"bytes"
	"encoding/json"
	"io"
	"log"

	"twitter-clone-backend/utils"
)

const (
	RED = "\033[31m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"
	BLUE = "\033[34m"
)

type AppHandler func(http.ResponseWriter, *http.Request) *models.AppError

// The ServeHTTP method calls the appHandler function and displays the returned error
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		fmt.Println(utils.ColorLog("error:", RED), e.Error.Error()) // goes to logging
		http.Error(w, e.Message, e.Code) // returned as response
	}
}

type AppMux struct {
	http.ServeMux
	middlewares		[]func(next http.Handler)	http.Handler
}

// func (mux *AppMux) registerMiddleware(next func(next http.Handler) http.Handler) { // next func(next http.Handler) http.Handler ???
// 	mux.middlewares = append(mux.middlewares, next)
// }

func (mux *AppMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.ColorLog(r.Method, GREEN), utils.ColorLog(r.URL.String(), GREEN))

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
	fmt.Println(utils.ColorLog("Request body:", BLUE))
	fmt.Printf("%s\n", prettyRequestBody)

	var current http.Handler = &mux.ServeMux // routing logic of http.ServeMux (maps URL paths to handlers)

	// the middleware wraps the current handler. ex: current = loggingMiddleware(authMiddleware(current))
	for _, next := range mux.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}

func (mux *AppMux) Handle(pattern string, handler interface{}) {
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