package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ranjabi/twitter-clone-backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	pgOnce sync.Once
	PgConn *pgxpool.Pool
	rdOnce sync.Once
	RdConn *redis.Client
	RbConn *amqp.Connection
	RbCh   *amqp.Channel
	err    error
)

func SetupConnection(ctx context.Context, cfg *config.Config) error {
	log.SetPrefix("DB: ")
	defer log.SetPrefix("")

	PgConn, err = GetPostgresConnection(ctx, cfg.PgConnString)
	if err != nil {
		return err
	}

	RdConn, err = GetRedisConnection()
	if err != nil {
		return err
	}

	RbConn, err = GetRabbitMqConnection()
	if err != nil {
		return err
	}
	RbCh, err = RbConn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func GetPostgresConnection(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	var initErr error

	pgOnce.Do(func() {
		var err error
		PgConn, err = pgxpool.New(ctx, connString)
		if err != nil {
			initErr = fmt.Errorf("error creating postgres database connection: %w", err)
			return
		}

		var testResult int
		err = PgConn.QueryRow(ctx, "SELECT 1").Scan(&testResult)
		if err != nil {
			initErr = fmt.Errorf("postgres failed to run test query: %w", err)
			return
		}

		log.Println("Postgres database connection successfully established")
	})

	if initErr != nil {
		return nil, initErr
	}

	return PgConn, nil
}

func GetRedisConnection() (*redis.Client, error) {
	rdOnce.Do(func() {
		RdConn = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			DB:   0, // Use default DB
		})
	})

	_, err := RdConn.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis database connection successfully obtained:")

	return RdConn, nil
}

func GetRabbitMqConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	return conn, nil
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
