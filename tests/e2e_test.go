package tests

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("LOADED BASE URL:", os.Getenv("BASE_URL"))

	m.Run()
}

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/health-check", os.Getenv("BASE_URL")), nil)
	assert.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{"data":{"Database":"OK","Server":"OK"}}`, string(resBody))
	res.Body.Close()
}