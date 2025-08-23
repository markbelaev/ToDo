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
		os.Exit(1)
	}
	slog.Info("config file loaded")

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	return &Config{
		DBConnectionString: dbConnectionString,
	}
}
