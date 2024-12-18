package healthCheck

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"twitter-clone-backend/middleware"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"
)

func HealthCheck(db *pgxpool.Pool, ctx context.Context) middleware.AppHandler {
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
			return &models.AppError{Err: err, Message: utils.ErrMsgFailedToSerializeResponseBody, Code: http.StatusInternalServerError}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

		return nil
	}
}
