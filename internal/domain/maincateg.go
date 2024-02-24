package domain

// MainCateg contains main category information with icon	info
type MainCateg struct {
	ID   int64           `json:"id"`
	Name string          `json:"name"`
	Type TransactionType `json:"type"`
	Icon Icon            `json:"icon"`
}
