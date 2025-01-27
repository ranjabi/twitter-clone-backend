package it

import (
	"context"

	"log"
	"path/filepath"
	"twitter-clone-backend/config"
	"twitter-clone-backend/db"
	"twitter-clone-backend/models"
	"twitter-clone-backend/usecases/tweet"
	"twitter-clone-backend/usecases/user"

	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // for pgx sql driver
	"github.com/joho/godotenv"

	"github.com/redis/go-redis/v9"
)

var (
	pgConn         *pgxpool.Pool
	rdConn         *redis.Client
	ctx            context.Context
	cfg            *config.Config
	migrationsPath string
	seedPath       string

	userRepository  user.UserRepository
	tweetRepository tweet.TweetRepository

	userService  user.Service
	tweetService tweet.Service

	validUser = models.User{
		Id:       1,
		Email:    "test@example.com",
		Username: "test",
		FullName: "Test test",
		Password: "password",
	}
	validUser2 = models.User{
		Id:       2,
		Email:    "test2@example.com",
		Username: "test2",
		FullName: "Test test 2",
		Password: "password",
	}
	notExistUser = models.User{
		Id: 100,
	}
	validTweet = models.Tweet{
		Id:      1,
		Content: "content",
		UserId:  1,
	}
	notExistTweet = models.Tweet{
		Id: 100,
	}
)

func TestMain(m *testing.M) {
	var err error
	ctx = context.Background()

	err = godotenv.Load("../../.env.dev.local")
	if err != nil {
		log.Fatal(err)
	}
	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pgConn, rdConn, err = db.Setup(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pgConn.Close()

	userRepository = user.NewRepository(ctx, pgConn, rdConn)
	tweetRepository = tweet.NewRepository(ctx, pgConn, rdConn)

	userService = user.NewService(ctx, cfg, userRepository)
	tweetService = tweet.NewService(tweetRepository, userRepository)

	// ----- Migration and seed start -----
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath = filepath.Join(cwd, "..", "..", "db", "migrations")
	seedPath = filepath.Join(cwd, "..", "..", "db", "seedtest")

	actions := []string{"migrate.reset", "migrate.up", "seed.up"}
	err = db.ApplyMigrationsAndSeed(ctx, cfg, actions, migrationsPath, seedPath, false)
	if err != nil {
		log.Fatal(err)
	}
	// ----- Migration and seed end -----

	code := m.Run()
	os.Exit(code)
}

func ResetAndSeed() error {
	actions := []string{"seed.down", "seed.up"}
	err := db.ApplyMigrationsAndSeed(ctx, cfg, actions, migrationsPath, seedPath, true)
	if err != nil {
		return err
	}
	return nil
}
