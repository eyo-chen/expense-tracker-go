package domain

// ChartData contains chart data
type ChartData struct {
	Labels   []string
	Datasets []float64
}

// DateRange contains start date and end date
type ChartDateRange struct {
	StartDate string
	EndDate   string
}
