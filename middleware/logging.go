package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"twitter-clone-backend/utils"
)

func Logging(next http.Handler) http.Handler {
	// diwrap pakai http.HandlerFunc supaya fungsi di bawah bisa jadi http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(utils.ColorLog(r.Method, utils.GREEN), utils.ColorLog(r.URL.String(), utils.GREEN))

		requestBodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading request body stream:", err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		var requestBody map[string]any
		_ = json.Unmarshal(requestBodyBytes, &requestBody)

		prettyRequestBody, err := json.MarshalIndent(requestBody, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling:", err)
			return
		}
		fmt.Println(utils.ColorLog("Request body:", utils.BLUE))
		fmt.Printf("%s\n", prettyRequestBody)

		next.ServeHTTP(w, r)
	})
}
