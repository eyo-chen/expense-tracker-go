package domain

// MainCateg contains main category information
type MainCateg struct {
	ID     int64
	Name   string
	Type   string
	IconID int64
}

// SubCateg contains sub category information
type SubCateg struct {
	ID          int64
	Name        string
	MainCategID int64
}
