package transaction

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
)

func cvtToDomainTransaction(t *Transaction, m *maincateg.MainCateg, s *subcateg.SubCateg, i *icon.Icon) domain.Transaction {
	return domain.Transaction{
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

func cvtToDomainMainCateg(c *maincateg.MainCateg, i *icon.Icon) *domain.MainCateg {
	return &domain.MainCateg{
		ID:   c.ID,
		Name: c.Name,
		Type: domain.CvtToMainCategType(c.Type),
		Icon: domain.Icon{
			ID:  i.ID,
			URL: i.URL,
		},
	}
}

func cvtToDomainIcon(i *icon.Icon) *domain.Icon {
	if i == nil {
		return nil
	}

	return &domain.Icon{
		ID:  i.ID,
		URL: i.URL,
	}
}

func cvtToDomainSubCateg(categ *subcateg.SubCateg) *domain.SubCateg {
	return &domain.SubCateg{
		ID:          categ.ID,
		Name:        categ.Name,
		MainCategID: categ.MainCategID,
	}
}
