package domain

// SubCateg contains sub category information
type SubCateg struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	MainCategID int64  `json:"main_category_id"`
}
