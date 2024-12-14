package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"twitter-clone-backend/utils"

	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/health-check" {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Unauthorized access", 400)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
		fmt.Println("token string: ", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.JWT_SIGNATURE_KEY), nil
		})
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims) // type assertion
		if !ok || !token.Valid {
			http.Error(w, err.Error(), 400)
			return
		}

		fmt.Printf("claims: \n%v\n", claims)

		ctx := context.WithValue(context.Background(), "userInfo", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}