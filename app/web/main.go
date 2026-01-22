package main

import (
	"fmt"
	"go-boilerplate/app/web/router"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	e := router.SetupEcho()

	e.GET("/", heartbeat)

	if err := e.Start(fmt.Sprintf(":%s", "8080")); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

func heartbeat(c *echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
