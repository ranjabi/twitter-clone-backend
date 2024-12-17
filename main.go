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
	mux.Handle("/register", handler.Register(conn, ctx))
	mux.Handle("/login", handler.Login(conn, ctx))
	mux.Handle("/tweet", handler.Tweet(conn, ctx))
	mux.Handle("/users/follow", handler.Follow(conn, ctx))
	mux.Handle("/users/unfollow", handler.Unfollow(conn, ctx))

	userRepository := user.NewRepository(conn, ctx)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	mux.HandleFunc("/v2/register", userHandler.HandleCreateUser)

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = mux

	fmt.Println("Server started at http://localhost:8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
