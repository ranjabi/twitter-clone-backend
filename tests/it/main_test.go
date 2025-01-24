package it

import (
	"context"
	"database/sql"

	"log"
	"path/filepath"
	"twitter-clone-backend/config"
	"twitter-clone-backend/db"
	"twitter-clone-backend/models"

	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // for pgx sql driver
	"github.com/joho/godotenv"

	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

var pgConn *pgxpool.Pool
var rdConn *redis.Client
var ctx context.Context

var validUser = models.User{
	Id:       11,
	Email:    "test@example.com",
	Username: "test",
	FullName: "Test test",
	Password: "password",
}
var validUser2 = models.User{
	Id:       12,
	Email:    "test2@example.com",
	Username: "test2",
	FullName: "Test test 2",
	Password: "password",
}
var notExistUser = models.User{
	Id: 100,
}
var cfg *config.Config

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

	pgConn, _, err = db.Setup(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// ----- Migration and seed start -----
	log.SetPrefix("DB: ")
	log.Println("Starting migration and seed...")
	db, err := sql.Open("pgx", cfg.PgConnString)
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
