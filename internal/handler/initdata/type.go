package initdata

type createInitDataInput struct {
	Income  []initDataMainCateg `json:"income"`
	Expense []initDataMainCateg `json:"expense"`
}

type initDataMainCateg struct {
	Name      string   `json:"name"`
	Icon      icon     `json:"icon"`
	SubCategs []string `json:"sub_categories"`
}

type icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}
