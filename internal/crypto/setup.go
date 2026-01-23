package crypto

import (
	"go-boilerplate/internal/crypto/portfolio"

	"github.com/labstack/echo/v5"
	"github.com/samber/do"
	"gorm.io/gorm"
)

var Injector *do.Injector

func NewInjector(
	db *gorm.DB,
) {
	Injector = do.New()

	if db != nil {
		do.Provide[*gorm.DB](Injector, func(i *do.Injector) (*gorm.DB, error) {
			return db, nil
		})
	}

	newPortfolio()
}

func NewHTTPHandlers(
	g *echo.Group,
	bearerMiddleware echo.MiddlewareFunc,
) {
	portfolio.NewHandler(
		g,
		do.MustInvoke[portfolio.Usecase](Injector),
		bearerMiddleware,
	)
}

func newPortfolio() {
	do.Provide[portfolio.Repository](Injector, func(i *do.Injector) (portfolio.Repository, error) {
		return portfolio.NewRepository(
			do.MustInvoke[*gorm.DB](i),
		), nil
	})

	do.Provide[portfolio.Usecase](Injector, func(i *do.Injector) (portfolio.Usecase, error) {
		return portfolio.NewUsecase(
			do.MustInvoke[portfolio.Repository](i),
		), nil
	})
}
