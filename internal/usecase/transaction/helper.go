package transaction

import (
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

var (
	weekDayFormat      = "Mon"
	dayFormat          = "01/02"
	yearAndMonthFormat = "2006-01"
)

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
			labels = append(labels, t.Format(dayFormat))
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
		/*
			Have to use the format "YYYY-MM"
			Because what frontend send is the random start date, and time duration
			For example, "2024-03-15", and "3 months"
			If we're using "YYYY-MM-DD" format in model and usecase, and assume the date is "01"("2024-03-01")
			Then there's no way to match the date("2024-03-15") with the data("2024-03-01")
			So we have to use "YYYY-MM" format to match the date of the month
			Because both "2024-03-15" and "2024-03-01" will be formatted to "2024-03"
		*/
		date := t.Format(yearAndMonthFormat)

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

func genLineChartData(dateToDate domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	if timeRangeType.IsDailyType() && timeRangeType != domain.TimeRangeTypeThreeMonths {
		return genDailyLineChartData(dateToDate, timeRangeType, start, end)
	}

	if timeRangeType == domain.TimeRangeTypeThreeMonths {
		return genThreeMonthsLineChartData(dateToDate, timeRangeType, start, end)
	}

	return genMonthlyLineChartData(dateToDate, timeRangeType, start, end)
}

func genDailyLineChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	labels := make([]string, 0, timeRangeType.GetVal())
	datasets := make([]float64, 0, timeRangeType.GetVal())

	accAmount := 0.0
	for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 0, 1) {
		var label string
		if timeRangeType == domain.TimeRangeTypeOneWeekDay {
			label = t.Format(weekDayFormat)
		} else {
			label = t.Format(dayFormat)
		}

		labels = append(labels, label)

		date := t.Format(time.DateOnly)
		if _, ok := dateToData[date]; ok {
			accAmount += dateToData[date]
		}

		datasets = append(datasets, accAmount)
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func genThreeMonthsLineChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
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
			labels = append(labels, t.Format(dayFormat))
			datasets = append(datasets, accAmount)
		}

		index++
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func genMonthlyLineChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	labels := make([]string, 0, timeRangeType.GetVal())
	datasets := make([]float64, 0, timeRangeType.GetVal())

	accAmount := 0.0
	for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 1, 0) {
		/*
			Have to use the format "YYYY-MM"
			Because what frontend send is the random start date, and time duration
			For example, "2024-03-15", and "3 months"
			If we're using "YYYY-MM-DD" format in model and usecase, and assume the date is "01"("2024-03-01")
			Then there's no way to match the date("2024-03-15") with the data("2024-03-01")
			So we have to use "YYYY-MM" format to match the date of the month
			Because both "2024-03-15" and "2024-03-01" will be formatted to "2024-03"
		*/
		date := t.Format(yearAndMonthFormat)

		// get the first 3 characters of the month
		shortMonth := t.Month().String()[:3]
		labels = append(labels, shortMonth)

		if _, ok := dateToData[date]; ok {
			accAmount += dateToData[date]
		}

		datasets = append(datasets, accAmount)
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}
