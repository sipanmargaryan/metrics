package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	"metrics/internal/config"
	"metrics/internal/database"
	"metrics/internal/logging"
	"metrics/migrations"
)

func main() {
	slog.SetDefault(logging.New(os.Stdout, "metrics-migrate"))
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		slog.Error("usage: migrate up|down (down reverts the latest migration)")
		os.Exit(1)
	}
	_ = godotenv.Load()
	db, err := database.ConnecPostgres(config.Load())
	if err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := migrations.Run(ctx, sqlDB, os.Args[1]); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations complete", "direction", os.Args[1])
}
