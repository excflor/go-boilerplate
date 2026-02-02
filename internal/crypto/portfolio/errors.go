package portfolio

import "errors"

// Sentinel errors for portfolio domain.
// These errors can be checked using errors.Is() for type-safe error handling.
var (
	// ErrNotFound is returned when a portfolio or holding is not found.
	ErrNotFound = errors.New("portfolio not found")

	// ErrUnauthorized is returned when a user tries to access a portfolio they don't own.
	ErrUnauthorized = errors.New("unauthorized access to portfolio")

	// ErrHoldingNotFound is returned when a holding is not found.
	ErrHoldingNotFound = errors.New("holding not found")

	// ErrInvalidInput is returned when input validation fails.
	ErrInvalidInput = errors.New("invalid input")
)
