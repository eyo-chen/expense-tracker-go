package monthlytrans

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	MonthlyTrans interfaces.MonthlyTransRepo
	Transaction  interfaces.TransactionRepo
}

func New(mt interfaces.MonthlyTransRepo, t interfaces.TransactionRepo) *UC {
	return &UC{MonthlyTrans: mt, Transaction: t}
}

func (u *UC) Create(ctx context.Context, date time.Time) error {
	trans, err := u.Transaction.GetMonthlyAggregatedData(ctx, date)
	if err != nil {
		return err
	}

	return u.MonthlyTrans.Create(ctx, date, trans)
}
