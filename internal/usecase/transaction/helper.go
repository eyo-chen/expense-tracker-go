package transaction

import (
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
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
		if _, ok := dateToData[date]; ok {
			datasets = append(datasets, dateToData[date])
		} else {
			// if there is no data for the weekday, append 0
			datasets = append(datasets, 0)
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
		if _, ok := dateToData[date]; ok {
			datasets = append(datasets, dateToData[date])
		} else {
			datasets = append(datasets, 0)
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
	var prevAmount float64
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
		if _, ok := dateToData[date]; ok {
			datasets = append(datasets, dateToData[date])
			prevAmount = dateToData[date]
		} else {
			// if there is no data for the weekday, append the previous amount
			datasets = append(datasets, prevAmount)
		}
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

/*
Note that the data for line chart is already accumulated when querying from the database
So we don't need to accumulate the data again
In three months line chart, we only need to
1. Find the correct date as the label (every 3 days)
2. Find the data for the date

How to find the correct data for the date?
For every found data, update the previous amount
When it's time to append the data to the datasets, append the previous amount

For example, the dateToData is

	{
		"2024-03-01": 100,
		"2024-03-02": 200,
		"2024-03-03": 300,
		"2024-03-04": 400,
		"2024-03-10": 500,
	}

For 03/01, the data is 100
For 03/04, the data is 400
For 03/07, the data is 400
For 03/10, the data is 500
*/
func genThreeMonthsLineChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	var prevAmount float64
	var index int
	labels := make([]string, 0, timeRangeType.GetVal())
	datasets := make([]float64, 0, timeRangeType.GetVal())

	for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 0, 1) {
		date := t.Format(time.DateOnly)

		if _, ok := dateToData[date]; ok {
			prevAmount = dateToData[date] // update the previous amount when there is data
		}

		if index%3 == 0 {
			labels = append(labels, t.Format(dayFormat))
			datasets = append(datasets, prevAmount)
		}

		index++
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func genMonthlyLineChartData(dateToData domain.DateToChartData, timeRangeType domain.TimeRangeType, start, end time.Time) domain.ChartData {
	var prevAmount float64
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

		if _, ok := dateToData[date]; ok {
			datasets = append(datasets, dateToData[date])
			prevAmount = dateToData[date]
		} else {
			// if there is no data for the weekday, append the previous amount
			datasets = append(datasets, prevAmount)
		}
	}

	return domain.ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}
