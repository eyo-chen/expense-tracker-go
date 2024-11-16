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

// UpdateTransactionInput represents input for updating transaction
type UpdateTransactionInput struct {
	ID          int64           `json:"id"`
	Type        TransactionType `json:"type"`
	MainCategID int64           `json:"main_category_id"`
	SubCategID  int64           `json:"sub_category_id"`
	Price       float64         `json:"price"`
	Date        time.Time       `json:"date"`
	Note        string          `json:"note"`
}

// AccInfo contains accumulated information
type AccInfo struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	TotalBalance float64 `json:"total_balance"`
}

// Filter contains filter for getting transactions
type Filter struct {
	StartDate    *time.Time
	EndDate      *time.Time
	MinPrice     *float64
	MaxPrice     *float64
	MainCategIDs []int64
	SubCategIDs  []int64
}

// Sort contains sort by and sort direction
type Sort struct {
	By  SortByType  `json:"sort_by"`
	Dir SortDirType `json:"sort_direction"`
}

// Search contains keyword for searching transactions
type Search struct {
	Keyword *string `json:"keyword"`
}

// Cursor contains next key for pagination
type Cursor struct {
	NextKey string `json:"next_key"`
	Size    int    `json:"size"`
}

// DecodedNextKeys is a slice of DecodedNextKeyInfo
type DecodedNextKeys []DecodedNextKeyInfo

// DecodedNextKeyInfo contains field and value of next key
type DecodedNextKeyInfo struct {
	Field string
	Value string
}

// GetTransOpt contains options for getting transactions
type GetTransOpt struct {
	Filter Filter `json:"filter"`
	Sort   *Sort  `json:"sort"`
	Search Search `json:"search"`
	Cursor Cursor `json:"cursor"`
}

// GetAccInfoQuery contains query for getting accumulated information
type GetAccInfoQuery struct {
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

// GetMonthlyDateRange contains date range for monthly data
type GetMonthlyDateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// MonthDayToTransactionType contains mapping from month day to transaction type
type MonthDayToTransactionType map[int]TransactionType

// MonthlyAggregatedData contains aggregated data for a month
type MonthlyAggregatedData struct {
	UserID       int64
	TotalIncome  float64
	TotalExpense float64
}
