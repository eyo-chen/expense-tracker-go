package validator

import (
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

// CreateMainCateg is a function that validates the input for creating main category.
func (v *Validator) CreateTransaction(t domain.CreateTransactionInput) bool {
	v.Check(t.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	v.Check(t.SubCategID > 0, "sub_category_id", "Sub category ID must be greater than 0")
	v.Check(t.Price > 0, "price", "Price must be greater than 0")
	v.Check(t.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(!t.Date.IsZero(), "date", "Date can't be empty")
	return v.Valid()
}

// GetTransaction is a function that validates the queries for getting transactions.
func (v *Validator) GetTransaction(q domain.GetQuery) bool {
	v.Check(isValidDateFormat(q.StartDate), "startDate", "Start date must be in YYYY-MM-DD format")
	v.Check(isValidDateFormat(q.EndDate), "endDate", "End date must be in YYYY-MM-DD format")
	v.Check(checkStartDateBeforeEndDate(q.StartDate, q.EndDate), "startDate", "Start date must be before end date")

	return v.Valid()
}

func (v *Validator) UpdateTransaction(t domain.UpdateTransactionInput) bool {
	v.Check(t.ID > 0, "id", "ID must be greater than 0")
	v.Check(t.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	v.Check(t.SubCategID > 0, "sub_category_id", "Sub category ID must be greater than 0")
	v.Check(t.Price > 0, "price", "Price must be greater than 0")
	v.Check(t.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(!t.Date.IsZero(), "date", "Date can't be empty")
	return v.Valid()
}

func (v *Validator) Delete(id int64) bool {
	v.Check(id > 0, "id", "ID must be greater than 0")
	return v.Valid()
}

func (v *Validator) GetAccInfo(q domain.GetAccInfoQuery) bool {
	v.Check(isValidDateFormat(q.StartDate), "startDate", "Start date must be in YYYY-MM-DD format")
	v.Check(isValidDateFormat(q.EndDate), "endDate", "End date must be in YYYY-MM-DD format")
	v.Check(checkStartDateBeforeEndDate(q.StartDate, q.EndDate), "startDate", "Start date must be before end date")

	return v.Valid()
}

func (v *Validator) GetChartData(dateRange domain.ChartDateRange, transactionType domain.TransactionType) bool {
	v.Check(isValidDateFormat(&dateRange.StartDate), "start_date", "Start date must be in YYYY-MM-DD format")
	v.Check(isValidDateFormat(&dateRange.EndDate), "end_date", "End date must be in YYYY-MM-DD format")
	v.Check(checkStartDateBeforeEndDate(&dateRange.StartDate, &dateRange.EndDate), "start_date", "Start date must be before end date")
	v.Check(transactionType.IsValid(), "type", "Transaction type must be income or expense")
	return v.Valid()
}

func (v *Validator) GetPieChartData(dateRange domain.ChartDateRange, transactionType domain.TransactionType) bool {
	v.Check(isValidDateFormat(&dateRange.StartDate), "start_date", "Start date must be in YYYY-MM-DD format")
	v.Check(isValidDateFormat(&dateRange.EndDate), "end_date", "End date must be in YYYY-MM-DD format")
	v.Check(checkStartDateBeforeEndDate(&dateRange.StartDate, &dateRange.EndDate), "start_date", "Start date must be before end date")
	v.Check(transactionType.IsValid(), "type", "Transaction type must be income or expense")
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
