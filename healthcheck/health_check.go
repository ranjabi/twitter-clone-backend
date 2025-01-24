package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/handlers"
	"twitter-clone-backend/models"
)

func HealthCheck(db *pgxpool.Pool, rdConn *redis.Client, ctx context.Context) handlers.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) *models.AppError {
		serverStatus := "OK"
		dbStatus := "OK"

		var test string
		err := db.QueryRow(ctx, "select 'OK'").Scan(&test)
		if err != nil {
			dbStatus = "NOT OK"
		}

		res, err := json.Marshal(map[string]any{
			"data": map[string]string{
				"Server":   serverStatus,
				"Database": dbStatus,
			},
		})
		if err != nil {
			return &models.AppError{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

		return nil
	}
}
