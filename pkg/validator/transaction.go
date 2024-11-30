package validator

import (
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

// CreateMainCateg validates the input for creating main category.
func (v *Validator) CreateTransaction(t domain.CreateTransactionInput) bool {
	v.Check(t.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	v.Check(t.SubCategID > 0, "sub_category_id", "Sub category ID must be greater than 0")
	v.Check(t.Price > 0, "price", "Price must be greater than 0")
	v.Check(t.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(!t.Date.IsZero(), "date", "Date can't be empty")
	return v.Valid()
}

// GetTransaction validates the input for getting transactions.
func (v *Validator) GetTransaction(o domain.GetTransOpt) bool {
	if o.Filter.StartDate != nil && o.Filter.EndDate != nil {
		v.Check(checkStartDateBeforeEndDateTime(*o.Filter.StartDate, *o.Filter.EndDate), "startDate", "Start date must be before end date")
	}

	return v.Valid()
}

// UpdateTransaction validates the input for updating transaction.
func (v *Validator) UpdateTransaction(t domain.UpdateTransactionInput) bool {
	v.Check(t.ID > 0, "id", "ID must be greater than 0")
	v.Check(t.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	v.Check(t.SubCategID > 0, "sub_category_id", "Sub category ID must be greater than 0")
	v.Check(t.Price > 0, "price", "Price must be greater than 0")
	v.Check(t.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(!t.Date.IsZero(), "date", "Date can't be empty")
	return v.Valid()
}

// Delete validates the input for deleting transaction.
func (v *Validator) Delete(id int64) bool {
	v.Check(id > 0, "id", "ID must be greater than 0")
	return v.Valid()
}

// GetAccInfo validates the input for getting account info.
func (v *Validator) GetAccInfo(q domain.GetAccInfoQuery, timeRangeType domain.TimeRangeType) bool {
	v.Check(isValidDateFormat(q.StartDate), "startDate", "Start date must be in YYYY-MM-DD format")
	v.Check(isValidDateFormat(q.EndDate), "endDate", "End date must be in YYYY-MM-DD format")
	v.Check(checkStartDateBeforeEndDate(q.StartDate, q.EndDate), "startDate", "Start date must be before end date")
	v.Check(timeRangeType.IsValid(), "time_range", "time range is invalid")
	return v.Valid()
}

// GetChartData validates the input for getting bar chart data.
func (v *Validator) GetBarChartData(dateRange domain.ChartDateRange, transactionType domain.TransactionType, timeRangeType domain.TimeRangeType) bool {
	v.Check(checkStartDateBeforeEndDateTime(dateRange.Start, dateRange.End), "start_date", "start date must be before end date")
	v.Check(transactionType.IsValid(), "type", "transaction type must be income or expense")
	v.Check(timeRangeType.IsValid(), "time_range", "time range is invalid")
	return v.Valid()
}

// GetPieChartData validates the input for getting pie chart data.
func (v *Validator) GetPieChartData(dateRange domain.ChartDateRange, transactionType domain.TransactionType) bool {
	v.Check(checkStartDateBeforeEndDateTime(dateRange.Start, dateRange.End), "start_date", "start date must be before end date")
	v.Check(transactionType.IsValid(), "type", "transaction type must be income or expense")
	return v.Valid()
}

// GetLineChartData validates the input for getting line chart data.
func (v *Validator) GetLineChartData(dateRange domain.ChartDateRange, timeRangeType domain.TimeRangeType) bool {
	v.Check(checkStartDateBeforeEndDateTime(dateRange.Start, dateRange.End), "start_date", "start date must be before end date")
	v.Check(timeRangeType.IsValid(), "time_range", "time range is invalid")
	return v.Valid()
}

// GetMonthlyData validates the date range for getting monthly data.
func (v *Validator) GetMonthlyData(dateRange domain.GetMonthlyDateRange) bool {
	v.Check(checkStartDateBeforeEndDateTime(dateRange.StartDate, dateRange.EndDate), "start_date", "start date must be before end date")
	return v.Valid()
}

func isValidDateFormat(dateString *string) bool {
	if dateString == nil {
		return true
	}

	_, err := time.Parse(time.DateOnly, *dateString)
	return err == nil
}

func checkStartDateBeforeEndDate(startDate, endDate *string) bool {
	if startDate == nil || endDate == nil {
		return true
	}

	start, _ := time.Parse(time.DateOnly, *startDate)
	end, _ := time.Parse(time.DateOnly, *endDate)
	return start.Before(end) || start.Equal(end)
}

func checkStartDateBeforeEndDateTime(start, end time.Time) bool {
	return start.Before(end) || start.Equal(end)
}
