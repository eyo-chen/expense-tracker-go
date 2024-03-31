package domain

import "time"

// ChartDataByWeekday contains chart data mapped by weekday
// e.g. Mon -> 12.0
type ChartDataByWeekday map[string]float64

// DateToChartData contains mapping from date to chart data
// e.g. 2021-01-01 -> 12.0
type DateToChartData map[string]float64

// ChartData contains chart data
type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []float64 `json:"datasets"`
}

// ChartDateRange contains start date and end date for chart data
type ChartDateRange struct {
	Start time.Time
	End   time.Time
}
