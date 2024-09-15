package initdata

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

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
			Icon: domain.DefaultIcon{
				ID:  v.Icon.ID,
				URL: v.Icon.URL,
			},
			SubCategs: v.SubCategs,
		}
	}

	return categ
}
