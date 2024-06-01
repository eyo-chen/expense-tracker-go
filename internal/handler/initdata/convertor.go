package initdata

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToInitData(input createInitDataInput) domain.InitData {
	return domain.InitData{
		Income:  cvtToMainCateg(input.Income),
		Expense: cvtToMainCateg(input.Expense),
	}
}

func cvtToMainCateg(input []initDataMainCateg) []domain.InitDataMainCateg {
	categ := make([]domain.InitDataMainCateg, len(input))

	for i, v := range input {
		categ[i] = domain.InitDataMainCateg{
			Name: v.Name,
			Icon: domain.Icon{
				ID:  v.Icon.ID,
				URL: v.Icon.URL,
			},
			SubCategs: v.SubCategs,
		}
	}

	return categ
}
