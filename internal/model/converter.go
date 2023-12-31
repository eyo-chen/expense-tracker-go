package model

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToDomainTransaction(t *Transaction) *domain.Transaction {
	return &domain.Transaction{
		ID:          t.ID.Hex(),
		UserID:      t.UserID,
		Type:        cvtToDomainType(t.Type),
		MainCategID: t.MainCategID,
		SubCategID:  t.SubCategID,
		Price:       t.Price,
		Date:        t.Date,
		Note:        t.Note,
	}
}

func cvtToModelTransaction(t *domain.Transaction) *Transaction {
	return &Transaction{
		UserID:      t.UserID,
		Type:        cvtToModelType(t.Type),
		MainCategID: t.MainCategID,
		SubCategID:  t.SubCategID,
		Price:       t.Price,
		Date:        t.Date,
		Note:        t.Note,
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
