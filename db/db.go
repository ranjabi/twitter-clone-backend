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

	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	migrationsPath := filepath.Join(cwd, "db", "migrations")
	var seedPath string
	env := os.Getenv("SEED")
	if strings.Contains(env, "test") {
		// called from root project, so cwd will be root
		seedPath = filepath.Join(cwd, "db", "seedtest")
	} else {
		seedPath = filepath.Join(cwd, "db", "seed")
	}
	actions := []string{"migrate.reset", "migrate.up", "seed.up"}
	err = ApplyMigrationsAndSeed(ctx, cfg, actions, migrationsPath, seedPath, false)
	if err != nil {
		return nil, nil, err
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

func ApplyMigrationsAndSeed(ctx context.Context, cfg *config.Config, actions []string, migrationsPath string, seedPath string, isSilent bool) error {
	if !isSilent {
		log.Println("Applying migrations and seed...")
	}

	if isSilent {
		goose.SetLogger(goose.NopLogger())
	}

	db, err := sql.Open("pgx", cfg.PgConnString)
	if err != nil {
		return fmt.Errorf("fail to open database connection: %w", err)
	}

	for _, action := range actions {
		if !isSilent {
			log.Printf("Starting %s...\n", action)
		}
		var options []goose.OptionsFunc

		parts := strings.Split(action, ".")
		path, cmd := parts[0], parts[1]
		if path == "migrate" {
			path = migrationsPath
		} else if path == "seed" {
			path = seedPath
			options = append(options, goose.WithNoVersioning())
		}

		if err := goose.RunWithOptionsContext(ctx, cmd, db, path, []string{}, options...); err != nil {
			return fmt.Errorf("db operation failed: %w", err)
		}
	}

	if !isSilent {
		log.Println("Migration or seed has been applied!")
	}
	return nil
}
