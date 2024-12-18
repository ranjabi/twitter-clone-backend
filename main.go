package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	_ "github.com/jackc/pgx/v5/stdlib"

	"twitter-clone-backend/db"
	"twitter-clone-backend/healthCheck"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/tweet"
	"twitter-clone-backend/user"
	"twitter-clone-backend/utils"
)

const ENV = ".env.development"

func applyMigrationsAndSeed(ctx context.Context) {
	fmt.Println("Applying migrations and seed...")

	db, err := sql.Open("pgx", utils.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	
	migrationsPath := filepath.Join(cwd, "db", "migrations")
	seedPath := filepath.Join(cwd, "db", "seed")

	fmt.Println("Starting migration reset...")
	if err := goose.RunWithOptionsContext(ctx, "reset", db, migrationsPath, []string{}); err != nil {
		log.Fatal("Migration reset failed:", err)
	}

	fmt.Println("Starting migration up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, migrationsPath, []string{}); err != nil {
		log.Fatal("Migration up failed:", err)
	}

	fmt.Println("Starting seed up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, seedPath, []string{}, goose.WithNoVersioning()); err != nil {
		log.Fatal("Seed up failed:", err)
	}

	fmt.Println("Migrations has been applied!")
}

func main() {
	err := godotenv.Load(ENV)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	conn, err := db.GetDbConnection(utils.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error getting database connection:", err)
	}
	defer db.CloseConnection()

	if ENV != ".env.development" {
		applyMigrationsAndSeed(ctx)
	}

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
