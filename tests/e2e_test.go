package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"twitter-clone-backend/config"

	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // for pgx sql driver
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

const (
	newUserUsername = "username0"
	newUserEmail    = "email0@email.com"
	newUserPassword = "password0"
	newUserToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIzIiwidXNlcm5hbWUiOiJ1c2VybmFtZTk5In0.EyzXI50ozNYwu7W9YuiLS-s7xojazWLfMXXUKJeIRwQ"

	userId       = 1
	userUsername = "username"
	userEmail    = "email@email.com"
	userPassword = "password"
	userToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVtYWlsQGVtYWlsLmNvbSIsInVzZXJJZCI6MSwidXNlcm5hbWUiOiJ1c2VybmFtZSJ9.m-53CNo6HJvF5quF5ducF4nMMPyEWb8dyUALGvqrpIM"

	user2Id = 2
	user3Id = 3

	userNotExistId    = 100000
	userNotExistEmail = "not-exist@email.com"

	tweetId             = 1
	tweetUpdatedContent = "Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. In vehicula lacinia lacus. Vestibulum tincidunt dui nunc, at interdum neque posuere id."

	tweetNotExistId = 100000

	newTweetContent = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas neque turpis, posuere non tortor ac, accumsan tempus est."
	newTweetUserId  = 1
)

func TestMain(m *testing.M) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = godotenv.Load("../.env.dev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("LOADED BASE URL:", os.Getenv("TEST_BASE_URL"))

	db, err := sql.Open("pgx", cfg.PgConnString)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// migration and seed
	// TODO: make dynamic and accept from env
	migrationsPath := filepath.Join(cwd, "..", "db", "migrations")
	seedPath := filepath.Join(cwd, "..", "db", "seed")

	fmt.Println("Starting migration reset...")
	if err := goose.RunWithOptionsContext(ctx, "reset", db, migrationsPath, []string{}); err != nil {
		log.Fatal("Migration reset failed:", err)
	}

	fmt.Println("Starting migration up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, migrationsPath, []string{}); err != nil {
		log.Fatal("Migration up failed:", err)
	}

	fmt.Println("Starting seed up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, seedPath, []string{}, goose.WithNoVersioning()); err != nil {
		log.Fatal("Seed up failed:", err)
	}

	fmt.Println("Migrations has been applied!")

	m.Run()
}

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/health-check", os.Getenv("TEST_BASE_URL")), nil)
	assert.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
	"data": {
		"Database": "OK",
		"Server": "OK"
	}
}`

	assert.Equal(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestUserRegister(t *testing.T) {
	reqBody := map[string]any{
		"username": newUserUsername,
		"email":    newUserEmail,
		"password": newUserPassword,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/register", os.Getenv("TEST_BASE_URL")), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := fmt.Sprintf(`{
		"message": "Account created successfully. Please login",
		"data": {
			"username": "%s",
			"email": "%s"
		}
	}`, newUserUsername, newUserEmail)

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestUserLogin(t *testing.T) {
	reqBody := map[string]any{
		"email":    userEmail,
		"password": userPassword,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", os.Getenv("TEST_BASE_URL")), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := fmt.Sprintf(`{
    "message": "Login success",
    "data": {
        "id": %d,
        "username": "%s",
		"email": "%s",
        "token": "%s"
    }
}`, userId, userUsername, userEmail, userToken)

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestUserLoginNotExist(t *testing.T) {
	reqBody := map[string]any{
		"email":    userNotExistEmail,
		"password": userPassword,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", os.Getenv("TEST_BASE_URL")), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
    "message": "User not found. Please create an account"
}`

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.JSONEq(t, expected, string(resBodyStr))
	res.Body.Close()
}

func TestUserFollow(t *testing.T) {
	reqBody := map[string]any{
		"followingId": user3Id,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/%d/follow", os.Getenv("TEST_BASE_URL"), user3Id), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
    "message": "User has been followed"
}`

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestUserFollowAlreadyFollowed(t *testing.T) {
	reqBody := map[string]any{
		"followingId": user2Id,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/%d/follow", os.Getenv("TEST_BASE_URL"), user2Id), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
    "message": "Failed to follow"
}`

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	res.Body.Close()
}

func TestUserUnfollow(t *testing.T) {
	reqBody := map[string]any{
		"followingId": user2Id,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/%d/unfollow", os.Getenv("TEST_BASE_URL"), user2Id), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
    "message": "User has been unfollowed"
}`

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestTweetCreate(t *testing.T) {
	reqBody := map[string]any{
		"content": newTweetContent,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tweets", os.Getenv("TEST_BASE_URL")), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	if data, ok := resBodyJson["data"].(map[string]any); ok {
		delete(data, "id")
		delete(data, "createdAt")
	}

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := fmt.Sprintf(`{
    "message": "Tweet created successfully",
    "data": {
        "content": "%s",
		"userId": %d,
		"isLiked": false,
		"likeCount": 0,
		"modifiedAt": null,
		"username": ""
    }
}`, newTweetContent, userId)

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestTweetUpdate(t *testing.T) {
	reqBody := map[string]any{
		"tweetId": tweetId,
		"content": tweetUpdatedContent,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tweets", os.Getenv("TEST_BASE_URL")), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	if data, ok := resBodyJson["data"].(map[string]any); ok {
		delete(data, "modifiedAt")
		delete(data, "createdAt")
	}

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := fmt.Sprintf(`{
    "message": "Tweet updated successfully",
    "data": {
        "id": %d,
        "content": "%s",
		"userId": %d
    }
}`, tweetId, tweetUpdatedContent, userId)

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestTweetUpdateNotFound(t *testing.T) {
	reqBody := map[string]any{
		"tweetId": tweetNotExistId,
		"content": tweetUpdatedContent,
	}
	reqBodyByte, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	reqBodyStr := string(reqBodyByte)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tweets", os.Getenv("TEST_BASE_URL")), strings.NewReader(reqBodyStr))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	if data, ok := resBodyJson["data"].(map[string]any); ok {
		delete(data, "modifiedAt")
	}

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
    "message": "Tweet not found"
}`

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	res.Body.Close()
}

func TestTweetDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tweets/%d", os.Getenv("TEST_BASE_URL"), tweetId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
    "message": "Tweet deleted successfully",
	"data": null
}`

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}

func TestTweetDeleteNotFound(t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tweets/%d", os.Getenv("TEST_BASE_URL"), tweetNotExistId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var resBodyJson map[string]any
	err = json.Unmarshal(resBody, &resBodyJson)
	assert.NoError(t, err)

	resBodyStr, err := json.MarshalIndent(resBodyJson, "", "\t")
	assert.NoError(t, err)

	expected := `{
    "message": "Tweet not found"
}`

	assert.JSONEq(t, expected, string(resBodyStr))
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	res.Body.Close()
}
