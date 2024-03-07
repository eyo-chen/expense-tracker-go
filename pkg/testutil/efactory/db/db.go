package db

// Database is responsible for inserting data into the database
type Database interface {
	// insert inserts a single data into the database
	Insert(InserParams) (interface{}, error)

	// insertList inserts a list of data into the database
	InsertList(InserListParams) ([]interface{}, error)

	// SetIDField sets the ID field of the given value
	// arg must be a pointer to a struct
	SetIDField(interface{}, int) error
}

type InserParams struct {
	StorageName string
	Value       interface{}
}

type InserListParams struct {
	StorageName string
	Values      []interface{}
}
