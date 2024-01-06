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

func cvtToModelMainCateg(c *domain.MainCateg, i *domain.Icon) *MainCateg {
	return &MainCateg{
		ID:     c.ID,
		Name:   c.Name,
		Type:   cvtToModelType(c.Type),
		IconID: i.ID,
	}
}

func cvtToDomainMainCateg(c *MainCateg, i *Icon) *domain.MainCateg {
	return &domain.MainCateg{
		ID:   c.ID,
		Name: c.Name,
		Type: cvtToDomainType(c.Type),
		Icon: cvtToDomainIcon(i),
	}
}

func cvtToModelIcon(i *domain.Icon) *Icon {
	if i == nil {
		return nil
	}

	return &Icon{
		ID:  i.ID,
		URL: i.URL,
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
