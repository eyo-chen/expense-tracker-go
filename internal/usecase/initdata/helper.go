package initdata

import (
	"fmt"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

// genAllMainCategs generates all types of main categories.
func genAllMainCategs(data domain.InitData) []domain.MainCateg {
	expenseMainCategs := genMainCategs(data.Expense, domain.TransactionTypeExpense)
	incomeMainCategs := genMainCategs(data.Income, domain.TransactionTypeIncome)
	mainCategs := append(expenseMainCategs, incomeMainCategs...)

	return mainCategs
}

// genMainCategs generates main categories based on the transaction type.
func genMainCategs(categs []domain.InitDataMainCateg, t domain.TransactionType) []domain.MainCateg {
	mainCategs := make([]domain.MainCateg, len(categs))
	for i, c := range categs {
		mainCategs[i] = domain.MainCateg{
			Name: c.Name,
			Icon: c.Icon,
			Type: t,
		}
	}

	return mainCategs
}

// genAllSubCategs generates all types of sub categories.
func genAllSubCategs(data domain.InitData, categs []domain.MainCateg) []domain.SubCateg {
	keyToMainCategIDMap := genKeyToMainCategIDMap(categs)
	expenseSubCategs := genSubCategs(data.Expense, domain.TransactionTypeExpense, keyToMainCategIDMap)
	incomeSubCategs := genSubCategs(data.Income, domain.TransactionTypeIncome, keyToMainCategIDMap)
	subCategs := append(expenseSubCategs, incomeSubCategs...)

	return subCategs
}

// genKeyToMainCategIDMap generates a map from the key to the main category ID.
func genKeyToMainCategIDMap(categs []domain.MainCateg) map[string]int64 {
	keyToMainCategID := make(map[string]int64)
	for _, c := range categs {
		key := fmt.Sprintf("%s:%s", c.Name, c.Type.ToString())
		keyToMainCategID[key] = c.ID
	}

	return keyToMainCategID
}

// genSubCategs generates sub categories based on the transaction type.
func genSubCategs(categs []domain.InitDataMainCateg, t domain.TransactionType, keyToMainCategIDMap map[string]int64) []domain.SubCateg {
	subCategs := make([]domain.SubCateg, 0, len(categs))

	for _, c := range categs {
		key := fmt.Sprintf("%s:%s", c.Name, t.ToString())
		mainCategID, ok := keyToMainCategIDMap[key]
		if !ok {
			continue
		}

		for _, subCategName := range c.SubCategs {
			subCategs = append(subCategs, domain.SubCateg{
				Name:        subCategName,
				MainCategID: mainCategID,
			})
		}
	}

	return subCategs
}
