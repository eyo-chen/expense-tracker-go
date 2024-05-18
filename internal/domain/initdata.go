package domain

type InitData struct {
	Income  []InitDataMainCateg `json:"income"`
	Expense []InitDataMainCateg `json:"expense"`
}

type InitDataMainCateg struct {
	Name      string   `json:"name"`
	Icon      Icon     `json:"icon"`
	SubCategs []string `json:"sub_categories"`
}
