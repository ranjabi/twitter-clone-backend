package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"twitter-clone-backend/db"
	"twitter-clone-backend/handler"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/tweet"
	"twitter-clone-backend/user"
	"twitter-clone-backend/utils"
)

func main() {
	ctx := context.Background()

	conn, err := db.GetDbConnection(utils.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error getting database connection:", err)
	}
	defer db.CloseConnection()

	mux := new(middleware.AppMux)
	mux.RegisterMiddleware(middleware.JwtAuthorization)

	mux.Handle("/health-check", handler.HealthCheck(conn, ctx))
	mux.Handle("/users/unfollow", handler.Unfollow(conn, ctx))

	userRepository := user.NewRepository(conn, ctx)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	tweetRepository := tweet.NewRepository(conn, ctx)
	tweetService := tweet.NewService(tweetRepository)
	tweetHandler := tweet.NewHandler(tweetService)

	// VERSION 2
	// if use mux.Handle then will goes into AppHandler
	mux.Handle("POST /v2/register", userHandler.HandleUserRegister)
	mux.Handle("POST /v2/login", userHandler.HandleUserLogin)

	mux.Handle("POST /v2/user/follow", userHandler.HandleFollowOtherUser)

	mux.Handle("POST /v2/tweet", tweetHandler.HandleTweetCreate)
	mux.Handle("PUT /v2/tweet", tweetHandler.HandleUpdateTweet)
	mux.Handle("DELETE /v2/tweet", tweetHandler.HandleDeleteTweet)

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = mux

	fmt.Printf("Server started at http://localhost%s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
