package main

import (
	"GIN/internal/config"
	"GIN/internal/database"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	slog.Info("Starting application...")

	cfg := config.Load()
	slog.Info("Configuration loaded")

	if err := database.Init(cfg); err != nil {
		slog.Error("Error initializing database.", "error", err)
		os.Exit(1)
	}
	defer database.Close()
}
