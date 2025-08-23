package main

import (
	"GIN/internal/config"
	"GIN/internal/database"
	"GIN/internal/routes"
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

	router := gin.Default()

	routes.SetupRoutes(router)

	slog.Info("Starting server...")
	if err := router.Run(":8080"); err != nil {
		slog.Error("server run failed", "err", err)
		os.Exit(1)
	}
}
