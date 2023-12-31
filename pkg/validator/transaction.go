package validator

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func (v *Validator) CreateTransaction(transaction *domain.Transaction) bool {
	v.Check(transaction.Type == "income" || transaction.Type == "expense", "type", "Type must be income or expense")
	v.Check(transaction.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	v.Check(transaction.SubCategID > 0, "sub_category_id", "Sub category ID must be greater than 0")
	v.Check(transaction.Price > 0, "price", "Price must be greater than 0")
	v.Check(transaction.Date != nil, "date", "Date can't be empty")
	return v.Valid()
}
