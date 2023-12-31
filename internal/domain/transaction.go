package domain

import "time"

type Transaction struct {
	ID        string     `json:"id"`
	UserID    int64      `json:"user_id"`
	Type      string     `json:"type"`
	MainCateg *MainCateg `json:"main_category"`
	SubCateg  *SubCateg  `json:"sub_category"`
	Price     int64      `json:"price"`
	Date      *time.Time `json:"date"`
	Note      string     `json:"note"`
}

type GetQuery struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type TransactionResp struct {
	Transactions []*Transaction `json:"transactions"`
	Income       int64          `json:"income"`
	Expense      int64          `json:"expense"`
	NetIncome    int64          `json:"net_income"`
}
