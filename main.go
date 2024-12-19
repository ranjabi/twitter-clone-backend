package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"twitter-clone-backend/db"
	"twitter-clone-backend/healthCheck"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/tweet"
	"twitter-clone-backend/user"
)

func main() {
	env := os.Getenv("ENV_NAME")
	err := godotenv.Load(env)
	fmt.Println("LOADED ENV:", env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	pgConn, rdConn := db.Setup(ctx)
	defer db.ClosePostgresConnection()

	mux := new(middleware.AppMux)
	mux.RegisterMiddleware(middleware.JwtAuthorization)

	mux.Handle("/health-check", healthCheck.HealthCheck(pgConn, rdConn, ctx))

	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)
	userHandler := user.NewHandler(userService)

	tweetRepository := tweet.NewRepository(ctx, pgConn, rdConn)
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
