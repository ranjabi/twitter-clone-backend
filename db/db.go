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
	conn	*pgxpool.Pool
)

func GetDbConnection(connString string) (*pgxpool.Pool, error) {
	var err error

	once.Do(func() {
		conn, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatal("Error to create database connection:", err)
		}
		fmt.Println("Database connection created")
	})

	// todo: make it LOG: msg
	fmt.Println("Successfully obtained database connection")

	return conn, err
}

func CloseConnection() {
	if conn != nil {
		conn.Close()
		fmt.Println("Database connection closed")
	}
}
