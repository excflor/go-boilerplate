package main

import (
	"fmt"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/crypto"
	"go-boilerplate/internal/databse"
	"go-boilerplate/internal/infra/auth"
	"go-boilerplate/internal/router"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := databse.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	bearerMiddleware := auth.BearerAuth()

	e := router.NewRouter(cfg)

	cryptoGroup := e.Group("/crypto-api")
	crypto.NewInjector(db)
	crypto.NewHTTPHandlers(cryptoGroup, bearerMiddleware)

	e.GET("/", heartbeat)

	if err := e.Start(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

func heartbeat(c *echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
