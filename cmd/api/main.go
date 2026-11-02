package main

import (
	"context"
	"errors"
	"log/slog"
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
	"metrics/internal/logging"
	"metrics/internal/middleware"
	"metrics/internal/routes"
)

// @title Metrics API
// @version 1.0
// @description Metrics order API.
// @BasePath /api/v1
func main() {
	slog.SetDefault(logging.New(os.Stdout, "metrics"))

	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found")
	}

	//  OLTP opentelemetry log collector
	shutdownLogs, err := logging.SetupOTLP("metrics")
	if err != nil {
		fatal("initialize log export", err)
	}
	defer shutdownLogs()

	cfg := config.Load()
	if err := helpers.RegisterValidators(); err != nil {
		fatal("register validators", err)
	}

	//  This opens connection pool
	db, err := database.ConnecPostgres(cfg)
	if err != nil {
		fatal("application failed", err)
	}

	// *sql.DB is not a single connection. It's a connection pool manager. under the hood gorm uses sql/database package
	sqlDB, err := db.DB()
	if err != nil {
		fatal("application failed", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.RequestLogger(slog.Default()))
	middleware.RegisterMetrics(router)
	router.Use(middleware.RecoverPanics())
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
		slog.Info("server starting", "address", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			fatal("application failed", err)
		}
	}()

	//  For graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("server shutting down")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fatal("server forced to shutdown", err)
	}

	// Close database connections
	if err := sqlDB.Close(); err != nil {
		fatal("database close failed", err)
	}

	slog.Info("server stopped")
}

func fatal(message string, err error) {
	slog.Error(message, "error", err)
	os.Exit(1)
}
