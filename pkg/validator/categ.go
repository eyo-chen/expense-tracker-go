package validator

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

func (v *Validator) CreateMainCateg(categ domain.CreateMainCategInput) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(categ.IconType.IsValid(), "icon_type", "Icon type must be default or custom")
	v.Check(categ.IconID > 0, "icon_id", "Icon ID must be greater than 0")
	return v.Valid()
}

func (v *Validator) UpdateMainCateg(categ domain.UpdateMainCategInput) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.Type.IsValid(), "type", "Type must be income or expense")
	v.Check(categ.IconType.IsValid(), "icon_type", "Icon type must be default or custom")
	return v.Valid()
}

func (v *Validator) CreateSubCateg(categ domain.SubCateg) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	return v.Valid()
}

func (v *Validator) UpdateSubCateg(categ domain.SubCateg) bool {
	v.Check(len(categ.Name) > 0, "name", "Name can't be empty")
	v.Check(categ.MainCategID > 0, "main_category_id", "Main category ID must be greater than 0")
	return v.Valid()
}
