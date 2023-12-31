package validator

import (
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

// CreateMainCateg is a function that validates the input for creating main category.
func (v *Validator) CreateTransaction(transaction *domain.Transaction) bool {
	v.Check(transaction.Type == "income" || transaction.Type == "expense", "type", "Type must be income or expense")
	v.Check(transaction.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	v.Check(transaction.SubCategID > 0, "sub_category_id", "Sub category ID must be greater than 0")
	v.Check(transaction.Price > 0, "price", "Price must be greater than 0")
	v.Check(transaction.Date != nil, "date", "Date can't be empty")
	return v.Valid()
}

// GetTransaction is a function that validates the queries for getting transactions.
func (v *Validator) GetTransaction(query *domain.GetQuery) bool {
	v.Check(query.StartDate != "", "startDate", "Start date can't be empty")
	v.Check(query.EndDate != "", "endDate", "End date can't be empty")
	v.Check(isValidDateFormat(query.StartDate), "startDate", "Start date must be in YYYY-MM-DD format")
	v.Check(isValidDateFormat(query.EndDate), "endDate", "End date must be in YYYY-MM-DD format")

	// Check if start date is before or equal end date
	startDate, _ := time.Parse(time.DateOnly, query.StartDate)
	endDate, _ := time.Parse(time.DateOnly, query.EndDate)
	v.Check(startDate.Before(endDate) || startDate.Equal(endDate), "startDate", "Start date must be before end date")

	return v.Valid()
}

func isValidDateFormat(dateString string) bool {
	_, err := time.Parse(time.DateOnly, dateString)
	return err == nil
}
