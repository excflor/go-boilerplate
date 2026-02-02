package auth

import (
	"errors"
	"go-boilerplate/pkg/response"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(e *echo.Echo, usecase Usecase) {
	h := &Handler{usecase: usecase}

	authGroup := e.Group("/auth")
	authGroup.POST("/login", h.Login)
	authGroup.POST("/refresh", h.RefreshToken)
	authGroup.POST("/logout", h.Logout)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	accessToken, refreshToken, err := h.usecase.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return response.Unauthorized(c, err.Error())
		}
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "login successful", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *Handler) RefreshToken(c *echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	accessToken, err := h.usecase.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrTokenExpired) {
			return response.Unauthorized(c, err.Error())
		}
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "token refreshed", map[string]string{
		"access_token": accessToken,
	})
}

func (h *Handler) Logout(c *echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.usecase.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return response.InternalServerError(c, "failed to logout")
	}

	return response.Success(c, "logout successful", nil)
}
