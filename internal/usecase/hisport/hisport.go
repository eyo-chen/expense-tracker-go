package hisport

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	historicalPortfolioService interfaces.HistoricalPortfolioService
}

func New(historicalPortfolioService interfaces.HistoricalPortfolioService) *UC {
	return &UC{historicalPortfolioService: historicalPortfolioService}
}

func (u *UC) Create(ctx context.Context, userID int32, date time.Time) error {
	if err := u.historicalPortfolioService.Create(ctx, userID, date); err != nil {
		return err
	}

	return nil
}
