package domain

import "time"

type CreateStock struct {
	UserID     int32
	Symbol     string
	Price      float64
	Quantity   int32
	ActionType string
	StockType  string
	CreatedAt  time.Time
}

type Portfolio struct {
	UserID              int32
	TotalPortfolioValue float64
	TotalGain           float64
	ROI                 float64
}

type StockInfo struct {
	Symbol     string
	Quantity   int32
	Price      float64
	AvgCost    float64
	Percentage float64
}

type AllStockInfo struct {
	Stocks []StockInfo
	ETF    []StockInfo
	Cash   []StockInfo
}
