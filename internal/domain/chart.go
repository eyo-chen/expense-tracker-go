package domain

// ChartData contains chart data
type ChartData struct {
	Labels   []string
	Datasets []float64
}

// ChartDateRange contains start date and end date for chart data
type ChartDateRange struct {
	StartDate string
	EndDate   string
}
