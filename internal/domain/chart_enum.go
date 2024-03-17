package domain

// ChartType is an enumeration of chart types
type ChartType int32

const (
	// UnSpecified is an enumeration of unspecified chart type
	UnSpecifiedChart ChartType = iota
	// ChartTypeBar is an enumeration of bar chart type
	ChartTypeBar
	// ChartTypePie is an enumeration of pie chart type
	ChartTypePie
	// ChartTypeLine is an enumeration of line chart type
	ChartTypeLine
)

// IsValid checks if the ChartType is valid
func (c ChartType) IsValid() bool {
	switch c {
	case ChartTypeBar, ChartTypePie, ChartTypeLine:
		return true
	}
	return false
}

// ToString returns the string representation of ChartType
func CvtToChartType(s string) ChartType {
	switch s {
	case "bar":
		return ChartTypeBar
	case "pie":
		return ChartTypePie
	case "line":
		return ChartTypeLine
	}
	return UnSpecifiedChart
}
