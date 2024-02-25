package transaction

import (
	"context"
	"errors"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "usecase/transaction"
)

type TransactionUC struct {
	Transaction interfaces.TransactionModel
	MainCateg   interfaces.MainCategModel
	SubCateg    interfaces.SubCategModel
}

func NewTransactionUC(t interfaces.TransactionModel, m interfaces.MainCategModel, s interfaces.SubCategModel) *TransactionUC {
	return &TransactionUC{
		Transaction: t,
		MainCateg:   m,
		SubCateg:    s,
	}
}

func (t *TransactionUC) Create(ctx context.Context, trans domain.CreateTransactionInput) error {
	// check if the main category exists
	mainCateg, err := t.MainCateg.GetByID(trans.MainCategID, trans.UserID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		return err
	}

	// check if the type in main category matches the transaction type
	if trans.Type != mainCateg.Type {
		logger.Error("Create Transaction failed", "package", PackageName, "err", domain.ErrTypeNotConsistent)
		return domain.ErrTypeNotConsistent
	}

	// check if the sub category exists
	subCateg, err := t.SubCateg.GetByID(trans.SubCategID, trans.UserID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		return err
	}

	// check if the sub category matches the main category
	if subCateg.MainCategID != trans.MainCategID {
		logger.Error("Create Transaction failed", "package", PackageName, "err", domain.ErrMainCategNotConsistent)
		return domain.ErrMainCategNotConsistent
	}

	if err := t.Transaction.Create(ctx, trans); err != nil {
		return err
	}

	return nil
}

func (t *TransactionUC) GetAll(ctx context.Context, query domain.GetQuery, user domain.User) ([]domain.Transaction, error) {
	transactions, err := t.Transaction.GetAll(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *TransactionUC) GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, user domain.User) (domain.AccInfo, error) {
	accInfo, err := t.Transaction.GetAccInfo(ctx, query, user.ID)
	if err != nil {
		return domain.AccInfo{}, err
	}

	return accInfo, nil
}
