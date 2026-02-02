package health

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Liveness returns 200 OK as long as the application is running.
func (h *HealthHandler) Liveness(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "alive"})
}

// Readiness returns 200 OK only if the application can serve requests (e.g. DB is up).
func (h *HealthHandler) Readiness(c *echo.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unready",
			"error":  "database connection unavailable",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unready",
			"error":  "database ping failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
}
