package portfolio

import (
	"context"
	"testing"

	"go-boilerplate/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a manual mock of the Repository interface.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, p *Portfolio) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*Portfolio, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Portfolio), args.Error(1)
}

func (m *MockRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Portfolio, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Portfolio), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, p *Portfolio) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) AddHolding(ctx context.Context, h *Holding) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

func (m *MockRepository) UpdateHolding(ctx context.Context, h *Holding) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

func (m *MockRepository) RemoveHolding(ctx context.Context, pID, hID uuid.UUID) error {
	args := m.Called(ctx, pID, hID)
	return args.Error(0)
}

func (m *MockRepository) GetHoldingsByPortfolioID(ctx context.Context, pID uuid.UUID) ([]Holding, error) {
	args := m.Called(ctx, pID)
	return args.Get(0).([]Holding), args.Error(1)
}

func TestCreatePortfolio(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	u := NewUsecase(mockRepo)

	userID := uuid.New()
	name := "My Crypto Portfolio"
	desc := "Testing description"
	req := dto.CreatePortfolioRequest{
		Name:        name,
		Description: &desc,
		Currency:    "USD",
	}

	// Set expectation: repository Create method should be called once
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*portfolio.Portfolio")).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(1).(*Portfolio)
			assert.Equal(t, userID, p.UserID)
			assert.Equal(t, name, p.Name)
			assert.Equal(t, "USD", p.Currency)
		})

	// Act
	result, err := u.CreatePortfolio(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, userID, result.UserID)
	mockRepo.AssertExpectations(t)
}

func TestGetPortfolio_Unauthorized(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	u := NewUsecase(mockRepo)

	userID := uuid.New()
	otherUserID := uuid.New()
	portfolioID := uuid.New()

	existingPortfolio := &Portfolio{
		ID:     portfolioID,
		UserID: otherUserID, // Different owner
		Name:   "Someone Else's Portfolio",
	}

	mockRepo.On("GetByID", mock.Anything, portfolioID).Return(existingPortfolio, nil)

	// Act
	result, err := u.GetPortfolio(context.Background(), userID, portfolioID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, assert.IsType(t, ErrUnauthorized, err) || err == ErrUnauthorized)
	mockRepo.AssertExpectations(t)
}
