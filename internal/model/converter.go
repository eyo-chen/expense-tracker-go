package model

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToDomainTransaction(t *Transaction) *domain.Transaction {
	return &domain.Transaction{
		ID:        t.ID.Hex(),
		UserID:    t.UserID,
		Type:      cvtToDomainType(t.Type),
		MainCateg: cvtToDomainMainCateg(t.MainCateg),
		SubCateg:  cvtToDomainSubCateg(t.SubCateg),
		Price:     t.Price,
		Date:      t.Date,
		Note:      t.Note,
	}
}

func cvtToModelTransaction(t *domain.Transaction) *Transaction {
	return &Transaction{
		UserID:    t.UserID,
		Type:      cvtToModelType(t.Type),
		MainCateg: cvtToModelMainCateg(t.MainCateg),
		SubCateg:  cvtToModelSubCateg(t.SubCateg),
		Price:     t.Price,
		Date:      t.Date,
		Note:      t.Note,
	}
}

func cvtToDomainTransactionResp(transactions []*Transaction) *domain.TransactionResp {
	var result domain.TransactionResp

	for _, t := range transactions {
		if t.Type == "1" {
			result.Income += t.Price
			result.NetIncome += t.Price
		} else {
			result.Expense += t.Price
			result.NetIncome -= t.Price
		}
		result.Transactions = append(result.Transactions, cvtToDomainTransaction(t))
	}

	return &result
}

func cvtToModelMainCateg(c *domain.MainCateg) *MainCateg {
	var icon *Icon
	if c.Icon != nil {
		icon = &Icon{
			ID:  c.Icon.ID,
			URL: c.Icon.URL,
		}
	}

	return &MainCateg{
		ID:   c.ID,
		Name: c.Name,
		Type: cvtToModelType(c.Type),
		Icon: icon,
	}
}

func cvtToDomainMainCateg(c *MainCateg) *domain.MainCateg {
	var icon *domain.Icon
	if c.Icon != nil {
		icon = &domain.Icon{
			ID:  c.Icon.ID,
			URL: c.Icon.URL,
		}
	}

	return &domain.MainCateg{
		ID:   c.ID,
		Name: c.Name,
		Type: cvtToDomainType(c.Type),
		Icon: icon,
	}
}

func cvtToModelSubCateg(c *domain.SubCateg) *SubCateg {
	return &SubCateg{
		ID:          c.ID,
		Name:        c.Name,
		MainCategID: c.MainCategID,
	}
}

func cvtToDomainSubCateg(categ *SubCateg) *domain.SubCateg {
	return &domain.SubCateg{
		ID:          categ.ID,
		Name:        categ.Name,
		MainCategID: categ.MainCategID,
	}
}

func cvtToModelType(t string) string {
	if t == "income" {
		return "1"
	}
	return "2"
}

func cvtToDomainType(t string) string {
	if t == "1" {
		return "income"
	}
	return "expense"
}
