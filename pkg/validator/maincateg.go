package validator

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func (v *Validator) AddMainCateg(categ *domain.MainCateg) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.IconID > 0, "icon_id", "Icon ID must be greater than 0")
	v.Check(categ.Type == "income" || categ.Type == "expense", "type", "Type must be income or expense")
	return v.Valid()
}
