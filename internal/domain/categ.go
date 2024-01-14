package domain

// MainCateg contains main category information with icon	info
type MainCateg struct {
	ID   int64         `json:"id"`
	Name string        `json:"name"`
	Type MainCategType `json:"type"`
	Icon *Icon         `json:"icon"`
}

// SubCateg contains sub category information
type SubCateg struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	MainCategID int64  `json:"main_category_id"`
}

// MainCategType is an enumeration of main category types
type MainCategType uint64

const (
	Income MainCategType = iota
	Expense
)

// String returns the string representation of MainCategType
func (m MainCategType) String() string {
	switch m {
	case Income:
		return "income"
	case Expense:
		return "expense"
	}
	return "unknown type"
}

// ModelValue returns the string enum of mysql
func (m MainCategType) ModelValue() string {
	switch m {
	case Income:
		return "1"
	case Expense:
		return "2"
	}
	return "0"
}

// IsValid checks if the MainCategType is valid
func (m MainCategType) IsValid() bool {
	switch m {
	case Income, Expense:
		return true
	}
	return false
}

// CvtToMainCategType converts string to MainCategType
func CvtToMainCategType(s string) MainCategType {
	switch s {
	case "income":
		return Income
	case "expense":
		return Expense
	case "1":
		return Income
	case "2":
		return Expense
	}
	return 0
}
