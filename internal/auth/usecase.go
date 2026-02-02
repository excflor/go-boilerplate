package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	infraAuth "go-boilerplate/internal/infra/auth"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
)

type Usecase interface {
	Login(ctx context.Context, username, password string) (string, string, error)
	RefreshToken(ctx context.Context, refreshTokenStr string) (string, error)
	Logout(ctx context.Context, refreshTokenStr string) error
}

type usecase struct {
	repo   Repository
	jwtSvc infraAuth.JWTService
}

func NewUsecase(repo Repository, jwtSvc infraAuth.JWTService) Usecase {
	return &usecase{
		repo:   repo,
		jwtSvc: jwtSvc,
	}
}

func (u *usecase) Login(ctx context.Context, username, password string) (string, string, error) {
	// In a real app, validate credentials against the DB
	// For this boilerplate, we'll use a mock success
	if username != "admin" || password != "admin" {
		return "", "", ErrInvalidCredentials
	}

	userID := "7ea078fa-aac0-4364-8f5f-ba69b136b8f7"
	accessToken, refreshTokenStr, err := u.jwtSvc.GeneratePair(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token pair: %w", err)
	}

	refreshToken := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := u.repo.Create(ctx, refreshToken); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshTokenStr, nil
}

func (u *usecase) RefreshToken(ctx context.Context, refreshTokenStr string) (string, error) {
	token, err := u.repo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return "", ErrInvalidToken
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = u.repo.DeleteByToken(ctx, refreshTokenStr)
		return "", ErrTokenExpired
	}

	accessToken, err := u.jwtSvc.GenerateToken(token.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func (u *usecase) Logout(ctx context.Context, refreshTokenStr string) error {
	return u.repo.DeleteByToken(ctx, refreshTokenStr)
}
