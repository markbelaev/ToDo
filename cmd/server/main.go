package main

import (
	"GIN/internal/config"
	"GIN/internal/database"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	slog.Info("Starting server...")
	if err := r.Run(":8080"); err != nil {
		slog.Error("server run failed", "err", err)
		os.Exit(1)
	}
}
