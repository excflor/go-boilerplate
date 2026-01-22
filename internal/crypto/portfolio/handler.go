package portfolio

import "github.com/labstack/echo/v5"

type Handler struct {
}

func NewHandler(g *echo.Group) {
	handler := &Handler{}

	g.GET("/", handler.GetAll)
}

func (h *Handler) GetAll(c *echo.Context) error {
	return nil
}
