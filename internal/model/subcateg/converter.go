package subcateg

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToDomainSubCateg(categ *SubCateg) *domain.SubCateg {
	return &domain.SubCateg{
		ID:          categ.ID,
		Name:        categ.Name,
		MainCategID: categ.MainCategID,
	}
}

func cvtToSubCateg(categ *domain.SubCateg, userID int64) *SubCateg {
	return &SubCateg{
		ID:          categ.ID,
		Name:        categ.Name,
		UserID:      userID,
		MainCategID: categ.MainCategID,
	}
}
