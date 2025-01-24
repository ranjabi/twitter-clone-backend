package it

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"strings"
	"twitter-clone-backend/constants"
	"twitter-clone-backend/db"
	"twitter-clone-backend/models"
	"twitter-clone-backend/usecases/user"
	"twitter-clone-backend/utils"

	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // for pgx sql driver
	"github.com/joho/godotenv"

	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var pgConn *pgxpool.Pool
var rdConn *redis.Client
var ctx context.Context

var validUser = models.User{
	Email:    "johndoe@example.com",
	Username: "johndoe",
	FullName: "John Doe",
	Password: "password",
}

func TestMain(m *testing.M) {
	ctx = context.Background()

	err := godotenv.Load("../../.env.dev.local")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	pgConn, _, err = db.Setup(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ----- Migration and seed start -----
	log.SetPrefix("DB: ")
	log.Println("Starting migration and seed...")
	db, err := sql.Open("pgx", utils.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: make dynamic and accept from env
	migrationsPath := filepath.Join(cwd, "..", "..", "db", "migrations")
	seedPath := filepath.Join(cwd, "..", "..", "db", "seed")

	log.Println("Starting migration reset...")
	if err := goose.RunWithOptionsContext(ctx, "reset", db, migrationsPath, []string{}); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting migration up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, migrationsPath, []string{}); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting seed up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, seedPath, []string{}, goose.WithNoVersioning()); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration and seed has been applied!")
	log.SetPrefix("")
	// ----- Migration and seed end -----

	code := m.Run()

	pgConn.Close()

	os.Exit(code)
}

func TestCreateUser_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	testUser := models.User{
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	newUser, err := userService.CreateUser(testUser)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, strings.ToLower(testUser.Email), newUser.Email)
}

func TestCreateUser_EmailAlreadyExist(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	testUser := models.User{
		Email:    validUser.Email,
		Password: faker.Password(),
	}

	_, err := userService.CreateUser(testUser)
	assert.EqualError(t, err, constants.EMAIL_ALREADY_EXIST_MSG)
}

func TestLoginUser_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	user, err := userService.CheckUserCredential(validUser.Email, validUser.Password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Email, validUser.Email)
	assert.Equal(t, user.Username, validUser.Username)
	assert.Equal(t, user.FullName, validUser.FullName)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	_, err := userService.CheckUserCredential(faker.Email(), faker.Password())
	assert.EqualError(t, err, constants.USER_NOT_FOUND_MSG)
}

func TestLoginUser_WrongCredential(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, userRepository)

	_, err := userService.CheckUserCredential(validUser.Email, faker.Password())
	assert.EqualError(t, err, constants.WRONG_CREDENTIAL_MSG)
}
