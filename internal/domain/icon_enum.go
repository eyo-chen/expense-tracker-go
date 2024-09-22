package domain

// IconType is an enumeration of icon types
type IconType int64

const (
	// IconTypeUnspecified is an enumeration of unspecified icon type
	IconTypeUnspecified IconType = iota
	// IconTypeDefault is an enumeration of default icon type
	IconTypeDefault
	// IconTypeCustom is an enumeration of custom icon type
	IconTypeCustom
)

// ToString returns the string representation of IconType
func (t IconType) ToString() string {
	switch t {
	case IconTypeDefault:
		return "default"
	case IconTypeCustom:
		return "custom"
	default:
		return "unknown type"
	}
}

// ToModelValue returns the string enum of mysql
func (t IconType) ToModelValue() string {
	switch t {
	case IconTypeDefault:
		return "1"
	case IconTypeCustom:
		return "2"
	}
	return "0"
}

// CvtToIconType converts string to IconType
func CvtToIconType(s string) IconType {
	switch s {
	case "1":
		return IconTypeDefault
	case "2":
		return IconTypeCustom
	case "default":
		return IconTypeDefault
	case "custom":
		return IconTypeCustom
	}
	return IconTypeUnspecified
}
