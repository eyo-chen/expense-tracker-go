package usecase

import (
	"context"
	"errors"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type transactionUC struct {
	Transaction TransactionModel
	MainCateg   MainCategModel
	SubCateg    SubCategModel
}

func newTransactionUC(t TransactionModel, m MainCategModel, s SubCategModel) *transactionUC {
	return &transactionUC{
		Transaction: t,
		MainCateg:   m,
		SubCateg:    s,
	}
}

func (t *transactionUC) Create(ctx context.Context, user *domain.User, transaction *domain.Transaction) error {
	// check if the main category exists
	mainCateg, err := t.MainCateg.GetFullInfoByID(transaction.MainCateg.ID, user.ID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		logger.Error("t.MainCateg.GetFullInfoByID failed", "package", "usecase", "err", err)
		return err
	}

	// check if the main category type matches the transaction type
	if mainCateg.Type != transaction.Type {
		return domain.ErrDataNotFound
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

	transaction.MainCateg = mainCateg
	transaction.SubCateg = subCateg
	if err := t.Transaction.Create(ctx, transaction); err != nil {
		logger.Error("t.Transaction.Create failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}

func (t *transactionUC) GetAll(ctx context.Context, query *domain.GetQuery, user *domain.User) (*domain.TransactionResp, error) {
	transactions, err := t.Transaction.GetAll(ctx, query, user.ID)
	if err != nil {
		logger.Error("t.Transaction.GetAll failed", "package", "usecase", "err", err)
		return nil, err
	}

	return transactions, nil
}
