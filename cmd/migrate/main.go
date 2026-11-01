package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"metrics/internal/config"
	"metrics/internal/database"
	"metrics/migrations"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		log.Fatal("usage: migrate up|down (down reverts the latest migration)")
	}
	_ = godotenv.Load()
	db, err := database.ConnecPostgres(config.Load())
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := migrations.Run(ctx, sqlDB, os.Args[1]); err != nil {
		log.Fatal(err)
	}
	log.Printf("migrations %s complete", os.Args[1])
}
