package transaction

import (
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

var (
	weekDayFormat = "Mon"
	dayFormat     = "01/02"
	monthFormat   = "Jun"
)

func cvtDateToTime(startDate, endDate string) (time.Time, time.Time, error) {
	if startDate == "" || endDate == "" {
		return time.Time{}, time.Time{}, nil
	}

	start, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func genChartData(dateToDate domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	if timeRangeType.IsDailyType() && timeRangeType != domain.TimeRangeTypeThreeMonths {
		return genDailyChartData(dateToDate, timeRangeType, start, end)
	}

	if timeRangeType == domain.TimeRangeTypeThreeMonths {
		return genThreeMonthsChartData(dateToDate, timeRangeType, start, end)
	}

	return genMonthlyChartData(dateToDate, timeRangeType, start, end)
}

func genDailyChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	labels := make([]string, 0, timeRangeType.GetVal())
	datasets := make([]float64, 0, timeRangeType.GetVal())

	for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 0, 1) {
		var label string
		if timeRangeType == domain.TimeRangeTypeOneWeekDay {
			label = t.Format(weekDayFormat)
		} else {
			label = t.Format(dayFormat)
		}

		labels = append(labels, label)

		date := t.Format(time.DateOnly)
		// if there is no data for the weekday, append 0
		if _, ok := dateToData[date]; !ok {
			datasets = append(datasets, 0)
		} else {
			datasets = append(datasets, dateToData[date])
		}
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func genThreeMonthsChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	var accAmount float64
	var index int
	labels := make([]string, 0, timeRangeType.GetVal())
	datasets := make([]float64, 0, timeRangeType.GetVal())

	for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 0, 1) {
		date := t.Format(time.DateOnly)

		if _, ok := dateToData[date]; ok {
			accAmount += dateToData[date]
		}

		if index%3 == 0 {
			labels = append(labels, t.Format(monthFormat))
			datasets = append(datasets, accAmount)
			accAmount = 0
		}

		index++
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func genMonthlyChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	labels := make([]string, 0, timeRangeType.GetVal())
	datasets := make([]float64, 0, timeRangeType.GetVal())

	for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 1, 0) {
		date := t.Format(time.DateOnly)

		// get the first 3 characters of the month
		shortMonth := t.Month().String()[:3]
		labels = append(labels, shortMonth)

		// if there is no data for the weekday, append 0
		if _, ok := dateToData[date]; !ok {
			datasets = append(datasets, 0)
		} else {
			datasets = append(datasets, dateToData[date])
		}
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}
