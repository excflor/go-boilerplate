package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTService_TokenLifecycle(t *testing.T) {
	// Arrange
	secret := "test-secret"
	expiryHours := 1
	svc := NewJWTService(secret, expiryHours)
	userID := "user-123"

	// Act - Generate
	token, err := svc.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Act - Validate
	claims, err := svc.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestJWTService_InvalidToken(t *testing.T) {
	// Arrange
	svc := NewJWTService("secret", 1)

	// Act
	claims, err := svc.ValidateToken("invalid.token.here")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTService_ExpiredToken(t *testing.T) {
	// Arrange
	secret := "secret"
	svc := &jwtService{
		secret:       []byte(secret),
		accessExpiry: -1 * time.Hour, // Expired
	}

	token, _ := svc.GenerateToken("user-123")

	// Act
	claims, err := svc.ValidateToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}
