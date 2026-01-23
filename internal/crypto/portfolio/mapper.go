package portfolio

import "go-boilerplate/internal/dto"

func ToPortfolioResponse(p *Portfolio) dto.PortfolioResponse {
	holdings := make([]dto.HoldingResponse, len(p.Holdings))
	for i, holding := range p.Holdings {
		holdings[i] = dto.HoldingResponse{
			ID:           holding.ID,
			PortfolioID:  holding.PortfolioID,
			Symbol:       holding.Symbol,
			AssetType:    holding.AssetType,
			Quantity:     holding.Quantity,
			AvgCost:      holding.AvgCost,
			CurrentPrice: holding.CurrentPrice,
			MarketValue:  holding.MarketValue,
			CreatedAt:    holding.CreatedAt,
			UpdatedAt:    holding.UpdatedAt,
		}
	}

	return dto.PortfolioResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Name:        p.Name,
		Description: p.Description,
		TotalValue:  p.TotalValue,
		Currency:    p.Currency,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Holdings:    holdings,
	}
}

func ToPortfolioListResponse(p []Portfolio) dto.PortfolioListResponse {
	portfolios := make([]dto.PortfolioResponse, len(p))
	for i, port := range p {
		portfolios[i] = ToPortfolioResponse(&port)
	}

	return dto.PortfolioListResponse{
		Portfolios: portfolios,
	}
}
