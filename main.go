package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func healthCheckHandler(conn *pgx.Conn, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server OK")

		var test string
		err := conn.QueryRow(ctx, "select 'OK'").Scan(&test)
		if err != nil {
			fmt.Fprintf(w, "Database NOT OK: %v\n", err)
		} else {
			fmt.Fprintf(w, "Database %v\n", test)
		}
	}
}

func main() {
	ctx := context.Background()

	databaseUrl := "postgres://postgres:123456@localhost:5432/postgres"
	conn, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}
	defer conn.Close(ctx)

	http.HandleFunc("/health-check", healthCheckHandler(conn, ctx))

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	err = http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Error starting server: ", err)
    }
}