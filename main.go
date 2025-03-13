package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"twitter-clone-backend/config"
	"twitter-clone-backend/db"
	"twitter-clone-backend/messagebroker"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/usecases/auth"
	"twitter-clone-backend/usecases/tweet"
	"twitter-clone-backend/usecases/user"
)

func main() {
	env := os.Getenv("ENV_NAME")
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loaded ENV:", env)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err = db.SetupConnection(ctx, cfg); err != nil {
		log.Fatal(err)
	}

	rabbitMq := messagebroker.NewRabbitMq(ctx, db.RbConn, db.RbCh)
	if err = rabbitMq.DeclareFeedQeueu(); err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migrationsPath := filepath.Join(cwd, "db", "migrations")
	var seedPath string
	if strings.Contains(env, "prod") {
		seedPath = filepath.Join(cwd, "db", "seed")
		actions := []string{"migrate.reset", "migrate.up", "seed.up"}
		err = db.ApplyMigrationsAndSeed(ctx, cfg, actions, migrationsPath, seedPath, false)
		if err != nil {
			log.Fatal(err)
		}
	}

	mux := new(AppMux)
	mux.RegisterMiddleware(middleware.Logging)
	mux.RegisterMiddleware(middleware.JwtAuthorization(cfg))

	userRepository := user.NewRepository(ctx, db.PgConn, db.RdConn)
	tweetRepository := tweet.NewRepository(ctx, db.PgConn, db.RdConn)

	authService := auth.NewService(ctx, cfg, userRepository)
	userService := user.NewService(ctx, userRepository)
	tweetService := tweet.NewService(tweetRepository, userRepository, rabbitMq)

	validate := validator.New(validator.WithRequiredStructEnabled())
	authHandler := auth.NewHandler(authService, validate)
	userHandler := user.NewHandler(userService, validate)
	tweetHandler := tweet.NewHandler(tweetService)

	// use mux.Handle so the error will goes into AppHandler
	mux.Handle("POST 	/v2/register", authHandler.HandleRegisterUser)
	mux.Handle("POST 	/v2/login", authHandler.HandleLoginUser)

	mux.Handle("POST 	/v2/users/{id}/follow", userHandler.HandleFollowOtherUser)
	mux.Handle("POST 	/v2/users/{id}/unfollow", userHandler.HandleUnfollowOtherUser)
	mux.Handle("GET		/v2/users/{username}", userHandler.HandleGetProfile)
	mux.Handle("GET		/v2/users/{id}/feed", userHandler.HandleGetFeed)

	mux.Handle("POST 	/v2/tweets", tweetHandler.HandleCreateTweet)
	mux.Handle("GET 	/v2/tweets/{id}", tweetHandler.HandleGetTweet)
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
