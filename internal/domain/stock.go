package domain

import "time"

type CreateStock struct {
	UserID     int32
	Symbol     string
	Price      float64
	Quantity   int32
	ActionType string
	StockType  string
	Date       time.Time
	CreatedAt  time.Time
}

type Portfolio struct {
	UserID              int32
	TotalPortfolioValue float64
	TotalGain           float64
	ROI                 float64
}

type StockInfo struct {
	Symbol     string  `json:"symbol"`
	Quantity   int32   `json:"quantity"`
	Price      float64 `json:"price"`
	AvgCost    float64 `json:"avg_cost"`
	Percentage float64 `json:"percentage"`
}

type AllStockInfo struct {
	Stocks []StockInfo `json:"stocks"`
	ETF    []StockInfo `json:"etf"`
	Cash   []StockInfo `json:"cash"`
}
