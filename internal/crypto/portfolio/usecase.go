package portfolio

import (
	"context"
	"fmt"
	"go-boilerplate/internal/dto"
	"time"

	"github.com/google/uuid"
)

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreatePortfolio(ctx context.Context, userID uuid.UUID, req dto.CreatePortfolioRequest) (*Portfolio, error) {
	portfolio := &Portfolio{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Currency:    "USD", // default
		IsActive:    true,
		TotalValue:  0,
	}

	if req.Currency != "" {
		portfolio.Currency = req.Currency
	}

	if err := u.repo.Create(ctx, portfolio); err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %w", err)
	}

	return portfolio, nil
}

func (u *usecase) GetPortfolio(ctx context.Context, userID, portfolioID uuid.UUID) (*Portfolio, error) {
	portfolio, err := u.repo.GetByID(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	if portfolio.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to portfolio")
	}

	return portfolio, nil
}

func (u *usecase) GetUserPortfolios(ctx context.Context, userID uuid.UUID) ([]Portfolio, error) {
	portfolios, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolios: %w", err)
	}

	return portfolios, nil
}

func (u *usecase) UpdatePortfolio(ctx context.Context, userID, portfolioID uuid.UUID, req dto.UpdatePortfolioRequest) (*Portfolio, error) {
	portfolio, err := u.GetPortfolio(ctx, userID, portfolioID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		portfolio.Name = *req.Name
	}
	if req.Description != nil {
		portfolio.Description = req.Description
	}
	if req.IsActive != nil {
		portfolio.IsActive = *req.IsActive
	}

	portfolio.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, portfolio); err != nil {
		return nil, fmt.Errorf("failed to update portfolio: %w", err)
	}

	return portfolio, nil
}

func (u *usecase) DeletePortfolio(ctx context.Context, userID, portfolioID uuid.UUID) error {
	// Verify ownership
	_, err := u.GetPortfolio(ctx, userID, portfolioID)
	if err != nil {
		return err
	}

	if err := u.repo.Delete(ctx, portfolioID); err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}

	return nil
}

func (u *usecase) AddHolding(ctx context.Context, userID, portfolioID uuid.UUID, req dto.AddHoldingRequest) (*Holding, error) {
	// Verify portfolio ownership
	_, err := u.GetPortfolio(ctx, userID, portfolioID)
	if err != nil {
		return nil, err
	}

	holding := &Holding{
		PortfolioID:  portfolioID,
		Symbol:       req.Symbol,
		AssetType:    req.AssetType,
		Quantity:     req.Quantity,
		AvgCost:      req.AvgCost,
		CurrentPrice: req.AvgCost, // Initially set to avg cost
		MarketValue:  req.Quantity * req.AvgCost,
	}

	if err := u.repo.AddHolding(ctx, holding); err != nil {
		return nil, fmt.Errorf("failed to add holding: %w", err)
	}

	// Update portfolio total value
	if err := u.recalculatePortfolioValue(ctx, portfolioID); err != nil {
		// Log error but don't fail the operation
		// In production, consider using background job for this
		fmt.Printf("Warning: failed to recalculate portfolio value: %v\n", err)
	}

	return holding, nil
}

func (u *usecase) RemoveHolding(ctx context.Context, userID, portfolioID, holdingID uuid.UUID) error {
	// Verify portfolio ownership
	_, err := u.GetPortfolio(ctx, userID, portfolioID)
	if err != nil {
		return err
	}

	if err := u.repo.RemoveHolding(ctx, portfolioID, holdingID); err != nil {
		return fmt.Errorf("failed to remove holding: %w", err)
	}

	// Update portfolio total value
	if err := u.recalculatePortfolioValue(ctx, portfolioID); err != nil {
		fmt.Printf("Warning: failed to recalculate portfolio value: %v\n", err)
	}

	return nil
}

func (u *usecase) GetPortfolioSummary(ctx context.Context, userID, portfolioID uuid.UUID) (*PortfolioSummary, error) {
	portfolio, err := u.GetPortfolio(ctx, userID, portfolioID)
	if err != nil {
		return nil, err
	}

	// Calculate metrics
	var totalInvested, totalReturn float64

	for _, holding := range portfolio.Holdings {
		totalInvested += holding.AvgCost * holding.Quantity
		totalReturn += (holding.CurrentPrice - holding.AvgCost) * holding.Quantity
	}

	var totalReturnPct float64
	if totalInvested > 0 {
		totalReturnPct = (totalReturn / totalInvested) * 100
	}

	summary := &PortfolioSummary{
		Portfolio:      *portfolio,
		TotalReturn:    totalReturn,
		TotalReturnPct: totalReturnPct,
		TotalInvested:  totalInvested,
		HoldingsCount:  len(portfolio.Holdings),
		DayChange:      0, // Would need historical data
		DayChangePct:   0, // Would need historical data
	}

	return summary, nil
}

func (u *usecase) recalculatePortfolioValue(ctx context.Context, portfolioID uuid.UUID) error {
	holdings, err := u.repo.GetHoldingsByPortfolioID(ctx, portfolioID)
	if err != nil {
		return err
	}

	var totalValue float64
	for _, holding := range holdings {
		totalValue += holding.MarketValue
	}

	portfolio := &Portfolio{
		ID:         portfolioID,
		TotalValue: totalValue,
		UpdatedAt:  time.Now(),
	}

	return u.repo.Update(ctx, portfolio)
}
