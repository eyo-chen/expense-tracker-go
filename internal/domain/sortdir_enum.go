package domain

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

// GetOperandFromSort returns the operand of the sort direction type
func GetOperandFromSort(s *Sort) string {
	if s == nil {
		return ">"
	}

	switch s.Dir {
	case SortDirTypeAsc:
		return ">"
	case SortDirTypeDesc:
		return "<"
	}

	// In MySQL, the default sort direction is ascending
	return ">"
}
