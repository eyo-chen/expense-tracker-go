package domain

import (
	"time"
)

// Transaction contains transaction information with main category and sub category
type Transaction struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	MainCateg MainCateg `json:"main_category"`
	SubCateg  SubCateg  `json:"sub_category"`
	Price     float64   `json:"price"`
	Date      time.Time `json:"date"`
	Note      string    `json:"note"`
}

// GetQuery contains query for getting transactions
type GetQuery struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
