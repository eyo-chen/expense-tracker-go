package domain

// SortByType is an enumeration of sort by types
type SortByType int64

const (
	// SortByTypeUnSpecified is an enumeration of unspecified sort by type
	SortByTypeUnSpecified SortByType = iota

	// SortByTypePrice is an enumeration of sort by price type
	SortByTypePrice

	// SortByTypeDate is an enumeration of sort by date type
	SortByTypeDate

	// SortByTypeType is an enumeration of sort by transaction type type
	SortByTypeTransType
)

// IsValid checks if the sort by type is valid
func (t SortByType) IsValid() bool {
	switch t {
	case SortByTypePrice, SortByTypeDate, SortByTypeTransType:
		return true
	}
	return false
}

// CvtToSortByType converts a string to a sort by type
func CvtToSortByType(s string) SortByType {
	switch s {
	case "price":
		return SortByTypePrice
	case "date":
		return SortByTypeDate
	case "type":
		return SortByTypeTransType
	}
	return SortByTypeUnSpecified
}

// String returns the string representation of the sort by type
func (t SortByType) String() string {
	switch t {
	case SortByTypePrice:
		return "price"
	case SortByTypeDate:
		return "date"
	case SortByTypeTransType:
		return "type"
	}
	return "unspecified"
}

// GetField returns the field name of the sort by type
func (t SortByType) GetField() string {
	switch t {
	case SortByTypePrice:
		return "Price"
	case SortByTypeDate:
		return "Date"
	case SortByTypeTransType:
		return "Type"
	}
	return ""
}
