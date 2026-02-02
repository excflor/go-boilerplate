package main

import (
	"context"
	"fmt"
	"go-boilerplate/internal/auth"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/crypto"
	"go-boilerplate/internal/crypto/portfolio"
	"go-boilerplate/internal/database"
	infraAuth "go-boilerplate/internal/infra/auth"
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

	// Auto-migrate domain entities
	if err := db.AutoMigrate(&auth.RefreshToken{}, &portfolio.Portfolio{}, &portfolio.Holding{}); err != nil {
		slog.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	jwtSvc := infraAuth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiryHours)

	e := router.NewRouter(cfg)

	// Health check endpoints
	healthHandler := health.NewHealthHandler(db)
	e.GET("/health/live", healthHandler.Liveness)
	e.GET("/health/ready", healthHandler.Readiness)

	// Auth domain setup
	authInjector := auth.NewInjector(db, jwtSvc)
	auth.RegisterHandlers(e, authInjector)

	// Crypto domain setup
	cryptoGroup := e.Group("/crypto-api")
	cryptoInjector := crypto.NewInjector(db)
	crypto.NewHTTPHandlers(cryptoGroup, cryptoInjector, jwtSvc)

	// Configure http.Server explicitly for better control and graceful shutdown support in Echo v5
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: e,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("starting server", "port", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	// Create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
	}

	// Close database connection
	if err := database.Close(db); err != nil {
		slog.Error("failed to close database connection", "error", err)
	} else {
		slog.Info("database connection closed")
	}

	slog.Info("server stopped")
}
