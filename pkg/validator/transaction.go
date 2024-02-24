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
	v.Check(checkStartDateBeforeEndDate(q), "startDate", "Start date must be before end date")

	return v.Valid()
}

func isValidDateFormat(dateString string) bool {
	if dateString == "" {
		return true
	}

	_, err := time.Parse(time.DateOnly, dateString)
	return err == nil
}

func checkStartDateBeforeEndDate(query domain.GetQuery) bool {
	if query.StartDate == "" || query.EndDate == "" {
		return true
	}

	startDate, _ := time.Parse(time.DateOnly, query.StartDate)
	endDate, _ := time.Parse(time.DateOnly, query.EndDate)
	return startDate.Before(endDate) || startDate.Equal(endDate)
}
