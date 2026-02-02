package main

import (
	"context"
	"fmt"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/crypto"
	"go-boilerplate/internal/database"
	"go-boilerplate/internal/infra/auth"
	"go-boilerplate/internal/infra/health"
	"go-boilerplate/internal/router"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: e,
	}

	go func() {
		slog.Info("starting server", "port", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
	}

	if err := database.Close(db); err != nil {
		slog.Error("failed to close database connection", "error", err)
	} else {
		slog.Info("database connection closed")
	}

	slog.Info("server stopped")
}
