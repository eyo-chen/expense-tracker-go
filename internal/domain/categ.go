package domain

// MainCateg contains main category information
type MainCateg struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Icon *Icon  `json:"icon"`
}

// SubCateg contains sub category information
type SubCateg struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	MainCategID int64  `json:"main_category_id"`
}
