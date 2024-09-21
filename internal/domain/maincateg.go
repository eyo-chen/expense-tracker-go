package domain

// MainCateg contains main category information with icon	info
type MainCateg struct {
	ID       int64           `json:"id"`
	Name     string          `json:"name"`
	Type     TransactionType `json:"type"`
	Icon     DefaultIcon     `json:"icon"`
	IconType IconType        `json:"icon_type"`
	IconData string          `json:"icon_data"`
}
