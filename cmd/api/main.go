package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "metrics/docs"

	"metrics/internal/config"
	"metrics/internal/database"
	"metrics/internal/helpers"
	"metrics/internal/middleware"
	"metrics/internal/routes"
)

// @title Metrics API
// @version 1.0
// @description Metrics order API.
// @BasePath /api/v1
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	cfg := config.Load()
	if err := helpers.RegisterValidators(); err != nil {
		log.Fatalf("register validators: %v", err)
	}

	//  This opens connection pool
	db, err := database.ConnecPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// *sql.DB is not a single connection. It's a connection pool manager. under the hood gorm uses sql/database package
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	router := gin.New()
	middleware.RegisterMetrics(router)
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.ErrorHandler()) //Error handling moddleware

	if cfg.Environment == "local" {
		router.SetTrustedProxies(nil) // For local development
		pprof.Register(router)        // to check if we have goruotine leak, heap etc..
		router.GET(
			"/swagger/*any",
			ginSwagger.WrapHandler(swaggerFiles.Handler),
		)
	}

	app := &config.App{
		DB:     db,
		Config: cfg,
	}

	routes.Register(router, app)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on %s", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	//  For graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	// Close database connections
	if err := sqlDB.Close(); err != nil {
		log.Fatal("database close error:", err)
	}

	log.Println("server stopped")
}
