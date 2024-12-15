package tests

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("LOADED BASE URL:", os.Getenv("BASE_URL"))

	m.Run()
}

type testSuite struct {
	suite.Suite
}

func TestLoadTestSuite(t *testing.T) {
	suite.Run(t, &testSuite{})
}

func(s *testSuite) TestHealthCheck() {
	fmt.Println("BASE URL:", os.Getenv("BASE_URL"))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/health-check", os.Getenv("BASE_URL")), nil)
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	s.NoError(err)

	s.Equal(`{"data":{"Database":"OK","Server":"OK"}}`, string(resBody))
	res.Body.Close()
}