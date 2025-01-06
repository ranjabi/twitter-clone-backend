package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"twitter-clone-backend/db"
	"twitter-clone-backend/healthcheck"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/tweet"
	"twitter-clone-backend/user"
)

func main() {
	env := os.Getenv("ENV_NAME")
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if env != "" {
		log.Println("LOADED ENV:", env)
	} else {
		log.Println("LOADED ENV: .env")
	}

	ctx := context.Background()
	pgConn, rdConn := db.Setup(ctx)
	defer db.ClosePostgresConnection()

	mux := new(AppMux)
	mux.RegisterMiddleware(middleware.Logging)
	mux.RegisterMiddleware(middleware.JwtAuthorization)

	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	tweetRepository := tweet.NewRepository(ctx, pgConn, rdConn)

	userService := user.NewService(ctx, userRepository)
	tweetService := tweet.NewService(tweetRepository, userRepository)

	userHandler := user.NewHandler(userService)
	tweetHandler := tweet.NewHandler(tweetService)

	// use mux.Handle so the error will goes into AppHandler
	mux.Handle("GET		/v2/health-check", healthcheck.HealthCheck(pgConn, rdConn, ctx))
	mux.Handle("POST 	/v2/register", userHandler.HandleRegisterUser)
	mux.Handle("POST 	/v2/login", userHandler.HandleLoginUser)

	mux.Handle("POST 	/v2/users/{id}/follow", userHandler.HandleFollowOtherUser)
	mux.Handle("POST 	/v2/users/{id}/unfollow", userHandler.HandleUnfollowOtherUser)
	mux.Handle("GET		/v2/users/{username}", userHandler.HandleGetProfile)
	mux.Handle("GET		/v2/users/{id}/feed", userHandler.HandleGetFeed)

	mux.Handle("POST 	/v2/tweets", tweetHandler.HandleCreateTweet)
	mux.Handle("PUT 	/v2/tweets", tweetHandler.HandleUpdateTweet)
	mux.Handle("DELETE 	/v2/tweets/{id}", tweetHandler.HandleDeleteTweet)
	mux.Handle("POST 	/v2/tweets/{id}/like", tweetHandler.HandleLikeTweet)
	mux.Handle("POST 	/v2/tweets/{id}/unlike", tweetHandler.HandleUnlikeTweet)

	server := new(http.Server)
	server.Addr = ":8080"
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		Debug:            true,
	})
	server.Handler = c.Handler(mux)

	fmt.Printf("Server started at http://localhost%s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
