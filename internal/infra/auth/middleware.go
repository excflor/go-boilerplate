package auth

import (
	"go-boilerplate/pkg/response"
	"strings"

	"github.com/labstack/echo/v5"
)

func BearerAuth(jwtSvc JWTService) echo.MiddlewareFunc {
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

			claims, err := jwtSvc.ValidateToken(token)
			if err != nil {
				return response.Unauthorized(c, "invalid or expired token")
			}

			c.Set("user_id", claims.UserID)

			return next(c)
		}
	}
}
