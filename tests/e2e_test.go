package tests

import (
	"context"
	"database/sql"
	"log"
	"math"
	"os"
	"path/filepath"
	"testing"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib" // for sql driver
	"github.com/pressly/goose/v3"
)

func TestAbs(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("pgx", "postgres://postgres:123456@localhost:5432/postgres")
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}
	
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: make dynamic and accept from env
	migrationsPath := filepath.Join(cwd, "..", "db", "seed")
	if err := goose.RunWithOptionsContext(context.Background(), "status", db, migrationsPath, []string{}, goose.WithNoVersioning()); err != nil {
		log.Fatal("Goose failed: ", err)
	}

    got := math.Abs(-1)
    if got != 1 {
        t.Errorf("Abs(-1) = %v; want 1", got)
    }
}

// func TestAbs2(t *testing.T) {
//     got := math.Abs(-1)
//     if got != 2 {
//         t.Errorf("Abs(-1) = %v; want 1", got)
//     }
// }