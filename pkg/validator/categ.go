package validator

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

func (v *Validator) CreateMainCateg(categ *domain.MainCateg) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(categ.IconType.IsValid(), "icon_type", "Icon type must be default or custom")
	v.Check(categ.IconData != "", "icon_data", "Icon data can't be empty")
	return v.Valid()
}

func (v *Validator) UpdateMainCateg(categ *domain.MainCateg) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(categ.IconType.IsValid(), "icon_type", "Icon type must be default or custom")
	v.Check(categ.IconData != "", "icon_data", "Icon data can't be empty")
	return v.Valid()
}

func (v *Validator) CreateSubCateg(categ *domain.SubCateg) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	return v.Valid()
}

func (v *Validator) UpdateSubCateg(categ *domain.SubCateg) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	return v.Valid()
}
