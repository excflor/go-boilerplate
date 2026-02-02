package crypto

import (
	"go-boilerplate/internal/crypto/portfolio"
	"go-boilerplate/internal/infra/auth"

	"github.com/labstack/echo/v5"
	"github.com/samber/do"
	"gorm.io/gorm"
)

// NewInjector creates and configures a new dependency injection container.
// Returns the injector instead of storing it globally for better testability.
func NewInjector(db *gorm.DB) *do.Injector {
	injector := do.New()

	if db != nil {
		do.Provide[*gorm.DB](injector, func(i *do.Injector) (*gorm.DB, error) {
			return db, nil
		})
	}

	newPortfolio(injector)
	return injector
}

// NewHTTPHandlers registers all HTTP handlers for the crypto domain.
// Accepts injector as parameter instead of using global state.
func NewHTTPHandlers(
	g *echo.Group,
	injector *do.Injector,
	jwtSvc auth.JWTService,
) {
	bearerMiddleware := auth.BearerAuth(jwtSvc)

	portfolio.NewHandler(
		g,
		do.MustInvoke[portfolio.Usecase](injector),
		bearerMiddleware,
	)
}

// newPortfolio registers portfolio-related dependencies in the injector.
func newPortfolio(injector *do.Injector) {
	do.Provide[portfolio.Repository](injector, func(i *do.Injector) (portfolio.Repository, error) {
		return portfolio.NewRepository(
			do.MustInvoke[*gorm.DB](i),
		), nil
	})

	do.Provide[portfolio.Usecase](injector, func(i *do.Injector) (portfolio.Usecase, error) {
		return portfolio.NewUsecase(
			do.MustInvoke[portfolio.Repository](i),
		), nil
	})
}
