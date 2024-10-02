package domain

// MainCateg contains main category information with icon	info
type MainCateg struct {
	ID       int64           `json:"id"`
	Name     string          `json:"name"`
	Type     TransactionType `json:"type"`
	IconType IconType        `json:"icon_type"`
	IconData string          `json:"icon_data"`
}

// CreateMainCategInput is the input for creating a main category
type CreateMainCategInput struct {
	Name     string
	Type     TransactionType
	IconType IconType
	IconID   int64
}

// UpdateMainCategInput is the input for updating a main category
type UpdateMainCategInput struct {
	ID       int64
	Name     string
	Type     TransactionType
	IconType IconType
	IconID   int64
}
