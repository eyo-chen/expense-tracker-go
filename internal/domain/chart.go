package domain

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
