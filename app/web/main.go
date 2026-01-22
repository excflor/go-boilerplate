package main

import (
	"fmt"
	"go-boilerplate/app/web/router"
	"go-boilerplate/config"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	cfg := config.NewConfig()
	e := router.SetupEcho()

	e.GET("/", heartbeat)

	if err := e.Start(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

func heartbeat(c *echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
