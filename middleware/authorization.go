package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/utils"

	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthorization(next http.Handler) http.Handler {
	// diwrap pakai http.HandlerFunc supaya fungsi di bawah bisa jadi http handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedPaths := []string{
			"/v2/health-check",
			"/v2/login",
			"/v2/register",
		}

		for _, path := range allowedPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			res, err := json.Marshal(models.ErrorResponse{Message: "Unauthorized access"})
			if err != nil {
				http.Error(w, errmsg.FAILED_TO_SERIALIZE_RESPONSE_BODY, http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(utils.JWT_SIGNATURE_KEY), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims) // type assertion
		if !ok || !token.Valid {
			// TODO: where this goes?
			http.Error(w, "Jwt claims failed", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(context.Background(), utils.UserInfoKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
