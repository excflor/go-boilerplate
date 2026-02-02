package auth

import (
	infraAuth "go-boilerplate/internal/infra/auth"

	"github.com/labstack/echo/v5"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func NewInjector(db *gorm.DB, jwtSvc infraAuth.JWTService) *do.Injector {
	injector := do.New()

	do.Provide(injector, func(i *do.Injector) (Repository, error) {
		return NewRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (Usecase, error) {
		repo := do.MustInvoke[Repository](i)
		return NewUsecase(repo, jwtSvc), nil
	})

	return injector
}

func RegisterHandlers(e *echo.Echo, injector *do.Injector) {
	usecase := do.MustInvoke[Usecase](injector)
	NewHandler(e, usecase)
}
