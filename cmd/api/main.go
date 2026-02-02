package main

import (
	"fmt"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/crypto"
	"go-boilerplate/internal/database"
	"go-boilerplate/internal/infra/auth"
	"go-boilerplate/internal/infra/health"
	"go-boilerplate/internal/router"
	"log/slog"
	"os"
)

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	bearerMiddleware := auth.BearerAuth()

	e := router.NewRouter(cfg)

	// Health check endpoints
	healthHandler := health.NewHealthHandler(db)
	e.GET("/health/live", healthHandler.Liveness)
	e.GET("/health/ready", healthHandler.Readiness)

	cryptoGroup := e.Group("/crypto-api")
	injector := crypto.NewInjector(db)
	crypto.NewHTTPHandlers(cryptoGroup, injector, bearerMiddleware)

	slog.Info("starting server", "port", cfg.AppPort)
	if err := e.Start(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
