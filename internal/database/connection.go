package database

import (
	"GIN/internal/config"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if cfg.DBConnectionString == "" {
		slog.Error("DB connection string is empty")
		return errors.New("DB connection string is empty")
	}

	slog.Info("Connecting to DB...")

	poolConfig, err := pgxpool.ParseConfig(cfg.DBConnectionString)
	if err != nil {
		slog.Error("DB connection string is wrong", "error", err)
		return err
	}

	Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("DB connection string is wrong", "error", err)
		return err
	}

	if err := Pool.Ping(ctx); err != nil {
		slog.Error("DB connection string is wrong", "error", err)
		return err
	}

	slog.Info("Connected to database")
	return nil
}

func Close() {
	if Pool != nil {
		slog.Info("Closing database connection")
		Pool.Close()
	}
}

func GetPool() *pgxpool.Pool {
	return Pool
}

func HealthCheck(ctx context.Context) error {
	return Pool.Ping(ctx)
}
