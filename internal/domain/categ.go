package domain

// MainCateg contains main category information
type MainCateg struct {
	ID   int64
	Name string
	Type string
	Icon *Icon
}

// SubCateg contains sub category information
type SubCateg struct {
	ID          int64
	Name        string
	MainCategID int64
}
