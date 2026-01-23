package dto

import (
	"time"

	"github.com/google/uuid"
)

// Portfolio Request DTOs
type CreatePortfolioRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Currency    string  `json:"currency,omitempty" validate:"omitempty,len=3"`
}

type UpdatePortfolioRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type AddHoldingRequest struct {
	Symbol    string  `json:"symbol" validate:"required,min=1,max=10"`
	AssetType string  `json:"asset_type" validate:"required,oneof=stock crypto bond etf"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	AvgCost   float64 `json:"avg_cost" validate:"required,gt=0"`
}

// Portfolio Response DTOs
type PortfolioResponse struct {
	ID          uuid.UUID         `json:"id"`
	UserID      uuid.UUID         `json:"user_id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	TotalValue  float64           `json:"total_value"`
	Currency    string            `json:"currency"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Holdings    []HoldingResponse `json:"holdings,omitempty"`
}

type HoldingResponse struct {
	ID           uuid.UUID `json:"id"`
	PortfolioID  uuid.UUID `json:"portfolio_id"`
	Symbol       string    `json:"symbol"`
	AssetType    string    `json:"asset_type"`
	Quantity     float64   `json:"quantity"`
	AvgCost      float64   `json:"avg_cost"`
	CurrentPrice float64   `json:"current_price"`
	MarketValue  float64   `json:"market_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PortfolioSummaryResponse struct {
	PortfolioResponse
	TotalReturn    float64 `json:"total_return"`
	TotalReturnPct float64 `json:"total_return_pct"`
	DayChange      float64 `json:"day_change"`
	DayChangePct   float64 `json:"day_change_pct"`
	TotalInvested  float64 `json:"total_invested"`
	HoldingsCount  int     `json:"holdings_count"`
}

type PortfolioListResponse struct {
	Portfolios []PortfolioResponse `json:"portfolios"`
	Pagination PaginationResponse  `json:"pagination,omitempty"`
}

// Path Parameters
type PortfolioIDParam struct {
	ID string `param:"id" validate:"required,uuid"`
}

type HoldingIDParam struct {
	PortfolioID string `param:"id" validate:"required,uuid"`
	HoldingID   string `param:"holdingId" validate:"required,uuid"`
}
