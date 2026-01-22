package crypto

import (
	"go-boilerplate/internal/crypto/portfolio"

	"github.com/labstack/echo/v5"
	"github.com/samber/do"
)

var Injector *do.Injector

func NewInjector() {
	Injector = do.New()

}

func NewHTTPHandlers(
	g *echo.Group,
) {
	portfolio.NewHandler(g)
}
