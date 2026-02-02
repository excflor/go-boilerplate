package portfolio

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, portfolio *Portfolio) error {
	if err := r.db.WithContext(ctx).Create(portfolio).Error; err != nil {
		return fmt.Errorf("failed to create portfolio: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Portfolio, error) {
	var portfolio Portfolio

	err := r.db.WithContext(ctx).
		Preload("Holdings").
		Where("id = ?", id).
		First(&portfolio).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get portfolio %s: %w", id, err)
	}

	return &portfolio, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Portfolio, error) {
	var portfolios []Portfolio

	err := r.db.WithContext(ctx).
		Preload("Holdings").
		Where("user_id = ? AND is_active = ?", userID, true).
		Find(&portfolios).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get portfolios: %w", err)
	}

	return portfolios, nil
}

func (r *repository) Update(ctx context.Context, portfolio *Portfolio) error {
	err := r.db.WithContext(ctx).
		Model(portfolio).
		Updates(portfolio).Error

	if err != nil {
		return fmt.Errorf("failed to update portfolio: %w", err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&Portfolio{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}

	return nil
}

func (r *repository) AddHolding(ctx context.Context, holding *Holding) error {
	if err := r.db.WithContext(ctx).Create(holding).Error; err != nil {
		return fmt.Errorf("failed to add holding: %w", err)
	}
	return nil
}

func (r *repository) UpdateHolding(ctx context.Context, holding *Holding) error {
	err := r.db.WithContext(ctx).
		Model(holding).
		Updates(holding).Error

	if err != nil {
		return fmt.Errorf("failed to update holding: %w", err)
	}

	return nil
}

func (r *repository) RemoveHolding(ctx context.Context, portfolioID, holdingID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND portfolio_id = ?", holdingID, portfolioID).
		Delete(&Holding{}).Error

	if err != nil {
		return fmt.Errorf("failed to remove holding: %w", err)
	}

	return nil
}

func (r *repository) GetHoldingsByPortfolioID(ctx context.Context, portfolioID uuid.UUID) ([]Holding, error) {
	var holdings []Holding

	err := r.db.WithContext(ctx).
		Where("portfolio_id = ?", portfolioID).
		Find(&holdings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get holdings: %w", err)
	}

	return holdings, nil
}
