package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"twitter-clone-backend/db"
	"twitter-clone-backend/healthCheck"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/tweet"
	"twitter-clone-backend/user"
	"twitter-clone-backend/utils"
)

func main() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	conn, err := db.GetDbConnection(utils.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error getting database connection:", err)
	}
	defer db.CloseConnection()

	mux := new(middleware.AppMux)
	mux.RegisterMiddleware(middleware.JwtAuthorization)

	mux.Handle("/health-check", healthCheck.HealthCheck(conn, ctx))

	userRepository := user.NewRepository(conn, ctx)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	tweetRepository := tweet.NewRepository(conn, ctx)
	tweetService := tweet.NewService(tweetRepository)
	tweetHandler := tweet.NewHandler(tweetService)

	// if use mux.Handle then will goes into AppHandler
	mux.Handle("POST 	/v2/register", userHandler.HandleRegisterUser)
	mux.Handle("POST 	/v2/login", userHandler.HandleLoginUser)

	mux.Handle("POST 	/v2/user/follow", userHandler.HandleFollowOtherUser)
	mux.Handle("POST 	/v2/user/unfollow", userHandler.HandleUnfollowOtherUser)
	mux.Handle("GET		/v2/users/{id}", userHandler.HandleGetUserProfile)

	mux.Handle("POST 	/v2/tweet", tweetHandler.HandleCreateTweet)
	mux.Handle("PUT 	/v2/tweet", tweetHandler.HandleUpdateTweet)
	mux.Handle("DELETE 	/v2/tweet", tweetHandler.HandleDeleteTweet)
	mux.Handle("POST 	/v2/tweet/{id}/like", tweetHandler.HandleLikeTweet)
	mux.Handle("POST 	/v2/tweet/{id}/unlike", tweetHandler.HandleUnlikeTweet)

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = mux

	fmt.Printf("Server started at http://localhost%s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
