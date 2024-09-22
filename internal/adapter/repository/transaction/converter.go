package transaction

import (
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func cvtToDomainTransaction(t Transaction, m maincateg.MainCateg, s subcateg.SubCateg) domain.Transaction {
	return domain.Transaction{
		ID:     t.ID,
		Type:   domain.CvtToTransactionType(t.Type),
		UserID: t.UserID,
		Price:  t.Price,
		Note:   t.Note,
		Date:   t.Date,
		MainCateg: domain.MainCateg{
			ID:       m.ID,
			Name:     m.Name,
			Type:     domain.CvtToTransactionType(m.Type),
			IconType: domain.CvtToIconType(m.IconType),
			IconData: m.IconData,
		},
		SubCateg: domain.SubCateg{
			ID:          s.ID,
			Name:        s.Name,
			MainCategID: m.ID, // use m.ID because in the get query, we don't reterive the subCateg.MainCategID
		},
	}
}

func cvtCreateTransInputToModelTransaction(t domain.CreateTransactionInput) Transaction {
	return Transaction{
		Type:        t.Type.ToModelValue(),
		UserID:      t.UserID,
		MainCategID: t.MainCategID,
		SubCategID:  t.SubCategID,
		Price:       t.Price,
		Note:        t.Note,
		Date:        t.Date,
	}
}

func cvtUpdateTransInputToModelTransaction(t domain.UpdateTransactionInput) Transaction {
	return Transaction{
		ID:          t.ID,
		Type:        t.Type.ToModelValue(),
		MainCategID: t.MainCategID,
		SubCategID:  t.SubCategID,
		Price:       t.Price,
		Note:        t.Note,
		Date:        t.Date,
	}
}

func cvtToDomainTransactionWithoutCategory(t Transaction) domain.Transaction {
	return domain.Transaction{
		ID:     t.ID,
		Type:   domain.CvtToTransactionType(t.Type),
		UserID: t.UserID,
		Price:  t.Price,
		Note:   t.Note,
		Date:   t.Date,
	}
}
