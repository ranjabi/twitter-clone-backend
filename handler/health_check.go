package handler

import (
	"context"
	"fmt"
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
	
	"twitter-clone-backend/models"
	"twitter-clone-backend/middleware"
)

func HealthCheck(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		fmt.Fprintln(w, "Server OK")

		var test string
		err := db.QueryRow(ctx, "select 'OK'").Scan(&test)
		if err != nil {
			fmt.Fprintf(w, "Database NOT OK: %v\n", err)
		} else {
			fmt.Fprintf(w, "Database %v\n", test)
		}

		return nil
	}
}