package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"twitter-clone-backend/utils"

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

func GetPostgresConnection(connString string) (*pgxpool.Pool, error) {
	var err error

	pgOnce.Do(func() {
		pgConn, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatal("Error to create postgres database connection:", err)
		}

		var testResult int
		err = pgConn.QueryRow(context.Background(), "SELECT 1").Scan(&testResult)
		if err != nil {
			log.Fatal("Postgres failed to connect:", err)
		}

		fmt.Println("Postgres database connection successfully obtained")
	})

	return pgConn, err
}

func ClosePostgresConnection() {
	if pgConn != nil {
		pgConn.Close()
		fmt.Println("Postgres database connection closed")
	}
}

func GetRedisConnection() *redis.Client {
	rdOnce.Do(func() {
		rdConn = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: "", // No password set
			DB:       0,  // Use default DB
		})
	})

	_, err := rdConn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis failed to connect:", err)
	}

	fmt.Println("Redis database connection successfully obtained")

	return rdConn
}

func applyMigrationsAndSeed(ctx context.Context) {
	fmt.Println("Applying migrations and seed...")

	db, err := sql.Open("pgx", utils.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath := filepath.Join(cwd, "db", "migrations")
	seedPath := filepath.Join(cwd, "db", "seed")

	fmt.Println("Starting migration up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, migrationsPath, []string{}); err != nil {
		log.Fatal("Migration up failed:", err)
	}

	fmt.Println("Starting seed up...")
	if err := goose.RunWithOptionsContext(ctx, "up", db, seedPath, []string{}, goose.WithNoVersioning()); err != nil {
		log.Fatal("Seed up failed:", err)
	}

	fmt.Println("Migrations has been applied!")
}

func Setup(ctx context.Context) (*pgxpool.Pool, *redis.Client) {
	pgConn, err := GetPostgresConnection(utils.GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error getting database connection:", err)
	}

	rdConn := GetRedisConnection()

	applyMigrationsAndSeed(ctx)

	return pgConn, rdConn
}
