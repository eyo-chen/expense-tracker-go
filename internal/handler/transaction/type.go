package transaction

import "time"

type createTransactionReq struct {
	Type        string    `json:"type"`
	MainCategID int64     `json:"main_category_id"`
	SubCategID  int64     `json:"sub_category_id"`
	Price       float64   `json:"price"`
	Date        time.Time `json:"date"`
	Note        string    `json:"note"`
}

type updateTransactionReq struct {
	Type        string    `json:"type"`
	MainCategID int64     `json:"main_category_id"`
	SubCategID  int64     `json:"sub_category_id"`
	Price       float64   `json:"price"`
	Date        time.Time `json:"date"`
	Note        string    `json:"note"`
}

type getTransactionResp struct {
	Transactions []transaction `json:"transactions"`
}

type mainCateg struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IconType string `json:"icon_type"`
	IconData string `json:"icon_data"`
}

type subCateg struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type transaction struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	MainCateg mainCateg `json:"main_category"`
	SubCateg  subCateg  `json:"sub_category"`
	Price     float64   `json:"price"`
	Note      string    `json:"note"`
	Date      time.Time `json:"date"`
}
