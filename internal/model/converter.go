package model

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToDomainTransaction(t *Transaction, m *MainCateg, s *SubCateg, i *Icon) *domain.Transaction {
	return &domain.Transaction{
		ID:        t.ID,
		UserID:    t.UserID,
		MainCateg: cvtToDomainMainCateg(m, i),
		SubCateg:  cvtToDomainSubCateg(s),
		Price:     t.Price,
		Note:      t.Note,
		Date:      t.Date,
	}
}

func cvtToModelTransaction(t *domain.Transaction) *Transaction {
	return &Transaction{
		UserID:      t.UserID,
		MainCategID: t.MainCateg.ID,
		SubCategID:  t.SubCateg.ID,
		Price:       t.Price,
		Note:        t.Note,
		Date:        t.Date,
	}
}

func cvtToDomainMainCateg(c *MainCateg, i *Icon) *domain.MainCateg {
	return &domain.MainCateg{
		ID:   c.ID,
		Name: c.Name,
		Type: domain.CvtToMainCategType(c.Type),
		Icon: cvtToDomainIcon(i),
	}
}

func cvtToDomainIcon(i *Icon) *domain.Icon {
	if i == nil {
		return nil
	}

	return &domain.Icon{
		ID:  i.ID,
		URL: i.URL,
	}
}

func cvtToDomainSubCateg(categ *SubCateg) *domain.SubCateg {
	return &domain.SubCateg{
		ID:          categ.ID,
		Name:        categ.Name,
		MainCategID: categ.MainCategID,
	}
}
