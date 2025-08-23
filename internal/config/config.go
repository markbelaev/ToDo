package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionString string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Error loading .env file")
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	return &Config{
		DBConnectionString: dbConnectionString,
	}
}
