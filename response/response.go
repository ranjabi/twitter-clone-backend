package response

import (
	"encoding/json"
	"net/http"
	"twitter-clone-backend/app"
	"twitter-clone-backend/errmsg"
)

type Payload struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Json(w http.ResponseWriter, statusCode int, data any) *app.Error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		data = make(map[string]interface{})
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return &app.Error{Err: err, Message: errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, Code: http.StatusInternalServerError}
	}

	return nil
}
