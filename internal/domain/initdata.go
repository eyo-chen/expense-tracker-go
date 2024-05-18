package domain

// InitData represents the initial data of the application.
type InitData struct {
	Income  []InitDataMainCateg `json:"income"`
	Expense []InitDataMainCateg `json:"expense"`
}

// InitDataMainCateg represents the main category of the initial data.
type InitDataMainCateg struct {
	Name      string   `json:"name"`
	Icon      Icon     `json:"icon"`
	SubCategs []string `json:"sub_categories"`
}
