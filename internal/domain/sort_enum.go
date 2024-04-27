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

// SortDirType is an enumeration of sort direction types
type SortDirType int8

const (
	// SortDirTypeUnSpecified is an enumeration of unspecified sort direction type
	SortDirTypeUnSpecified SortDirType = iota

	// SortDirTypeAsc is an enumeration of ascending sort direction type
	SortDirTypeAsc

	// SortDirTypeDesc is an enumeration of descending sort direction type
	SortDirTypeDesc
)

// IsValid checks if the sort direction type is valid
func (t SortDirType) IsValid() bool {
	switch t {
	case SortDirTypeAsc, SortDirTypeDesc:
		return true
	}
	return false
}

// CvtToSortDirType converts a string to a sort direction type
func CvtToSortDirType(s string) SortDirType {
	switch s {
	case "asc":
		return SortDirTypeAsc
	case "desc":
		return SortDirTypeDesc
	}
	return SortDirTypeUnSpecified
}

// String returns the string representation of the sort direction type
func (t SortDirType) String() string {
	switch t {
	case SortDirTypeAsc:
		return "asc"
	case SortDirTypeDesc:
		return "desc"
	}
	return "unspecified"
}

// GetOperand returns the operand of the sort direction type
func (t SortDirType) GetOperand() string {
	switch t {
	case SortDirTypeAsc:
		return ">"
	case SortDirTypeDesc:
		return "<"
	}

	// In MySQL, the default sort direction is ascending
	return ">"
}
