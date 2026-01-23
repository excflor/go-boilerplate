package response

import (
	"go-boilerplate/internal/dto"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Success sends a successful response
func Success(c *echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, dto.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a bad request error response
func BadRequest(c *echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, dto.BaseResponse{
		Success: false,
		Error:   message,
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, dto.BaseResponse{
		Success: false,
		Error:   message,
	})
}

// NotFound sends a not found error response
func NotFound(c *echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, dto.BaseResponse{
		Success: false,
		Error:   message,
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(c *echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, dto.BaseResponse{
		Success: false,
		Error:   message,
	})
}
