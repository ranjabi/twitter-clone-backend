package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	pgOnce sync.Once
	pgConn *pgxpool.Pool
	rdOnce sync.Once
	rdConn *redis.Client
)

func GetRedisConnection() *redis.Client {
	rdOnce.Do(func() {
		rdConn = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: "", // No password set
			DB:       0,  // Use default DB
		})
		fmt.Println("Redis database connection created")
	})

	fmt.Println("Redis database connection successfully obtained") // todo: this will never fail

	return rdConn
}

func GetPostgresConnection(connString string) (*pgxpool.Pool, error) {
	var err error

	pgOnce.Do(func() {
		pgConn, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatal("Error to create postgres database connection:", err)
		}
		fmt.Println("Postgres database connection created")
	})

	// todo: make it LOG: msg
	fmt.Println("Postgres database connection successfully obtained")

	return pgConn, err
}

func ClosePostgresConnection() {
	if pgConn != nil {
		pgConn.Close()
		fmt.Println("Postgres database connection closed")
	}
}
