package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"twitter-clone-backend/middleware"
	"twitter-clone-backend/handler"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	databaseUrl := "postgres://postgres:123456@localhost:5432/postgres"
	db, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}
	defer db.Close()

	mux := new(middleware.AppMux)
	mux.RegisterMiddleware(middleware.JwtAuthorization)

	mux.Handle("/health-check", handler.HealthCheck(db, ctx))
	mux.Handle("/register", handler.Register(db, ctx))
	mux.Handle("/login", handler.Login(db, ctx))
	mux.Handle("/tweet", handler.Tweet(db, ctx))
	mux.Handle("/users/follow", handler.Follow(db, ctx))

	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = mux

	fmt.Println("Server started at http://localhost:8080")
	err = server.ListenAndServe()
    if err != nil {
        log.Fatal("Error starting server: ", err)
    }
}