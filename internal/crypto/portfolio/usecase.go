package portfolio

import "context"

type Usecase struct {
}

func NewUsecase() AssetUsecase {
	return &Usecase{}
}

// GetUserPortfolio implements [AssetUsecase].
func (u *Usecase) GetUserPortfolio(ctx context.Context, userID uint) ([]Asset, error) {
	panic("unimplemented")
}
