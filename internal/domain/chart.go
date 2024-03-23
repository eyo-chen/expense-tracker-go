package domain

// ChartDataByWeekday contains chart data mapped by weekday
// e.g. Mon -> 12.0
type ChartDataByWeekday map[string]float64

// ChartData contains chart data
type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []float64 `json:"datasets"`
}

// ChartDateRange contains start date and end date for chart data
type ChartDateRange struct {
	StartDate string
	EndDate   string
}
