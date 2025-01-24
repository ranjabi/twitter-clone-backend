package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"
	"sync"
	"twitter-clone-backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

var (
	pgOnce sync.Once
	pgConn *pgxpool.Pool
	rdOnce sync.Once
	rdConn *redis.Client
)

func Setup(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, *redis.Client, error) {
	log.SetPrefix("DB: ")
	defer log.SetPrefix("")

	pgConn, err := GetPostgresConnection(ctx, cfg.PgConnString)
	if err != nil {
		return nil, nil, err
	}

	rdConn, err := GetRedisConnection()
	if err != nil {
		return nil, nil, err
	}

	if env := os.Getenv("ENV_NAME"); strings.Contains(env, "prod") {
		applyMigrationsAndSeed(ctx, cfg)
	}

	return pgConn, rdConn, nil
}

func GetPostgresConnection(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	var initErr error

	pgOnce.Do(func() {
		var err error
		pgConn, err = pgxpool.New(ctx, connString)
		if err != nil {
			initErr = fmt.Errorf("error creating postgres database connection: %w", err)
			return
		}

		var testResult int
		err = pgConn.QueryRow(ctx, "SELECT 1").Scan(&testResult)
		if err != nil {
			initErr = fmt.Errorf("postgres failed to run test query: %w", err)
			return
		}

		log.Println("Postgres database connection successfully established")
	})

	if initErr != nil {
		return nil, initErr
	}

	return pgConn, nil
}

func GetRedisConnection() (*redis.Client, error) {
	rdOnce.Do(func() {
		rdConn = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			DB:   0, // Use default DB
		})
	})

	_, err := rdConn.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis database connection successfully obtained:")

	return rdConn, nil
}

func applyMigrationsAndSeed(ctx context.Context, cfg *config.Config) {
	log.Println("Applying migrations and seed...")

	db, err := sql.Open("pgx", cfg.PgConnString)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath := filepath.Join(cwd, "db", "migrations")
	seedPath := filepath.Join(cwd, "db", "seed")

	log.Println("Starting migration up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, migrationsPath, []string{}); err != nil {
		log.Fatal("Migration up failed:", err)
	}

	log.Println("Starting seed up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, seedPath, []string{}, goose.WithNoVersioning()); err != nil {
		log.Fatal("Seed up failed:", err)
	}

	log.Println("Migrations has been applied!")
}

func ClosePostgresConnection() {
	if pgConn != nil {
		pgConn.Close()
		log.Println("Postgres database connection closed")
	}
}
