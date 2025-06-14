package stock

import (
	"context"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	stockService interfaces.StockService
}

func New(stockService interfaces.StockService) *UC {
	return &UC{stockService: stockService}
}

func (u *UC) Create(ctx context.Context, stock domain.CreateStock) (string, error) {
	return u.stockService.Create(ctx, stock)
}

func (u *UC) GetPortfolioInfo(ctx context.Context, userID int32) (domain.Portfolio, error) {
	return u.stockService.GetPortfolioInfo(ctx, userID)
}

func (u *UC) GetStockInfo(ctx context.Context, userID int32) (domain.AllStockInfo, error) {
	return u.stockService.GetStockInfo(ctx, userID)
}
