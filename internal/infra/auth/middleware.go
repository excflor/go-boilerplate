package auth

import (
	"go-boilerplate/pkg/response"
	"strings"

	"github.com/labstack/echo/v5"
)

func BearerAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, "missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Unauthorized(c, "invalid token format")
			}

			token := parts[1]

			// hardcoded for example, should have token validator
			if token != "secret-token-123" {
				return response.Unauthorized(c, "invalid or expired token")
			}

			// hardcoded for example, user id should be from validated token
			c.Set("user_id", "7ea078fa-aac0-4364-8f5f-ba69b136b8f7")

			return next(c)
		}
	}
}
