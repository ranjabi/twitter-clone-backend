package config

import (
	"fmt"
	"os"
)

type Config struct {
	PgConnString string
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
		HashAlg:   "HS256",
		JwtSecret: "secret",
		SaltRound: 10,
	}

	return &config, nil
}
