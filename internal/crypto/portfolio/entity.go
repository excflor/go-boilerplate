package portfolio

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Asset struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"index" json:"user_id"`
	Symbol    string          `gorm:"size:10;index" json:"symbol"`
	Amount    decimal.Decimal `gorm:"type:numeric(32,18)" json:"amount"`
	CostBasis decimal.Decimal `gorm:"type:numeric(32,18)" json:"cost_basis"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type AssetUsecase interface {
	GetUserPortfolio(ctx context.Context, userID uint) ([]Asset, error)
}
