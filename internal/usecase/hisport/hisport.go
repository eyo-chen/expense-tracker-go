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

func (u *UC) GetPortfolioValue(ctx context.Context, userID int32, dateOption string) ([]string, []float64, error) {
	dates, values, err := u.historicalPortfolioService.GetPortfolioValue(ctx, userID, dateOption)
	if err != nil {
		return nil, nil, err
	}

	return dates, values, nil
}

func (u *UC) GetGain(ctx context.Context, userID int32, dateOption string) ([]string, []float64, error) {
	dates, values, err := u.historicalPortfolioService.GetGain(ctx, userID, dateOption)
	if err != nil {
		return nil, nil, err
	}

	return dates, values, nil
}
