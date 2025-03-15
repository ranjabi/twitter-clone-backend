package config

import (
	"fmt"
	"os"
)

type Config struct {
	PgConnString string
	RbConnString string
	HashAlg      string
	JwtSecret    string
	SaltRound    int
}

func Load() (*Config, error) {
	config := Config{
		PgConnString: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"),
		),
		RbConnString: "amqp://guest:guest@localhost:5672/",
		HashAlg:      "HS256",
		JwtSecret:    "secret",
		SaltRound:    10,
	}

	return &config, nil
}
