package domain

import (
	"time"
)

// Transaction contains transaction information with main category and sub category
type Transaction struct {
	ID        int64           `json:"id"`
	Type      TransactionType `json:"type"`
	UserID    int64           `json:"user_id"`
	MainCateg MainCateg       `json:"main_category"`
	SubCateg  SubCateg        `json:"sub_category"`
	Price     float64         `json:"price"`
	Date      time.Time       `json:"date"`
	Note      string          `json:"note"`
}

// CreateTransactionInput represents input for creating transaction
type CreateTransactionInput struct {
	UserID      int64           `json:"user_id"`
	Type        TransactionType `json:"type"`
	MainCategID int64           `json:"main_category_id"`
	SubCategID  int64           `json:"sub_category_id"`
	Price       float64         `json:"price"`
	Date        time.Time       `json:"date"`
	Note        string          `json:"note"`
}

// GetQuery contains query for getting transactions
type GetQuery struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
