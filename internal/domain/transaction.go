package domain

import "time"

type Transaction struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Type        string     `json:"type"`
	MainCategID int64      `json:"main_category_id"`
	SubCategID  int64      `json:"sub_category_id"`
	Price       int64      `json:"price"`
	Date        *time.Time `json:"date"`
	Note        string     `json:"note"`
}
