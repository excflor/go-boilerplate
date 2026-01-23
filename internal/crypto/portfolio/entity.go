package portfolio

import (
	"context"
	"go-boilerplate/internal/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Portfolio struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	TotalValue  float64        `json:"total_value" gorm:"type:decimal(15,2);default:0"`
	Currency    string         `json:"currency" gorm:"type:varchar(3);default:'USD'"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Holdings []Holding `json:"holdings,omitempty" gorm:"foreignKey:PortfolioID"`
}
type Holding struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PortfolioID  uuid.UUID      `json:"portfolio_id" gorm:"type:uuid;not null;index"`
	Symbol       string         `json:"symbol" gorm:"type:varchar(10);not null"`
	AssetType    string         `json:"asset_type" gorm:"type:varchar(20);not null"`
	Quantity     float64        `json:"quantity" gorm:"type:decimal(15,8);not null"`
	AvgCost      float64        `json:"avg_cost" gorm:"type:decimal(15,2);not null"`
	CurrentPrice float64        `json:"current_price" gorm:"type:decimal(15,2);default:0"`
	MarketValue  float64        `json:"market_value" gorm:"type:decimal(15,2);default:0"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type PortfolioSummary struct {
	Portfolio
	TotalReturn    float64 `json:"total_return"`
	TotalReturnPct float64 `json:"total_return_pct"`
	DayChange      float64 `json:"day_change"`
	DayChangePct   float64 `json:"day_change_pct"`
	TotalInvested  float64 `json:"total_invested"`
	HoldingsCount  int     `json:"holdings_count"`
}

func (Portfolio) TableName() string {
	return "portfolios"
}

func (Holding) TableName() string {
	return "holdings"
}

type Usecase interface {
	CreatePortfolio(ctx context.Context, userID uuid.UUID, req dto.CreatePortfolioRequest) (*Portfolio, error)
	GetPortfolio(ctx context.Context, userID, portfolioID uuid.UUID) (*Portfolio, error)
	GetUserPortfolios(ctx context.Context, userID uuid.UUID) ([]Portfolio, error)
	UpdatePortfolio(ctx context.Context, userID, portfolioID uuid.UUID, req dto.UpdatePortfolioRequest) (*Portfolio, error)
	DeletePortfolio(ctx context.Context, userID, portfolioID uuid.UUID) error
	AddHolding(ctx context.Context, userID, portfolioID uuid.UUID, req dto.AddHoldingRequest) (*Holding, error)
	RemoveHolding(ctx context.Context, userID, portfolioID, holdingID uuid.UUID) error
	GetPortfolioSummary(ctx context.Context, userID, portfolioID uuid.UUID) (*PortfolioSummary, error)
}

type Repository interface {
	Create(ctx context.Context, portfolio *Portfolio) error
	GetByID(ctx context.Context, id uuid.UUID) (*Portfolio, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Portfolio, error)
	Update(ctx context.Context, portfolio *Portfolio) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddHolding(ctx context.Context, holding *Holding) error
	UpdateHolding(ctx context.Context, holding *Holding) error
	RemoveHolding(ctx context.Context, portfolioID, holdingID uuid.UUID) error
	GetHoldingsByPortfolioID(ctx context.Context, portfolioID uuid.UUID) ([]Holding, error)
}
