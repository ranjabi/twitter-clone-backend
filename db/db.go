package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once	sync.Once
	pool	*pgxpool.Pool
)

func GetDbPool(connString string) (*pgxpool.Pool, error) {
	var err error

	once.Do(func() {
		pool, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatal("Error to create database pool:", err)
		}
		fmt.Println("Database pool created")
	})
	
	// todo: make it LOG: msg
	fmt.Println("Successfully obtained database pool")

	return pool, err
}

func ClosePool() {
	if pool != nil {
		pool.Close()
		fmt.Println("Database pool closed")
	}
}