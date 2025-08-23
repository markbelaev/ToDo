package main

import (
	"GIN/internal/config"
	"GIN/internal/database"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.Info("Starting...")

	cfg := config.Load()
	slog.Info("Config loaded")

	if err := database.Init(cfg); err != nil {
		slog.Error("database init failed", "err", err)
		os.Exit(1)
	}
	defer database.Close()
}
