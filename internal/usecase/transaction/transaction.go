package transaction

import (
	"context"
	"errors"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
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

func (t *TransactionUC) Create(ctx context.Context, user *domain.User, transaction *domain.Transaction) error {
	// check if the main category exists
	_, err := t.MainCateg.GetByID(transaction.MainCateg.ID, user.ID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		logger.Error("t.MainCateg.GetFullInfoByID failed", "package", "usecase", "err", err)
		return err
	}

	// check if the sub category exists
	subCateg, err := t.SubCateg.GetByID(transaction.SubCateg.ID, user.ID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		logger.Error("t.SubCateg.GetByID failed", "package", "usecase", "err", err)
		return err
	}

	// check if the sub category matches the main category
	if subCateg.MainCategID != transaction.MainCateg.ID {
		return domain.ErrDataNotFound
	}

	if err := t.Transaction.Create(ctx, transaction); err != nil {
		logger.Error("t.Transaction.Create failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}

func (t *TransactionUC) GetAll(ctx context.Context, query *domain.GetQuery, user *domain.User) ([]domain.Transaction, error) {
	transactions, err := t.Transaction.GetAll(ctx, query, user.ID)
	if err != nil {
		logger.Error("t.Transaction.GetAll failed", "package", "usecase", "err", err)
		return nil, err
	}

	return transactions, nil
}
