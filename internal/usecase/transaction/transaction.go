package transaction

import (
	"context"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/codeutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "usecase/transaction"
)

type TransactionUC struct {
	Transaction interfaces.TransactionRepo
	MainCateg   interfaces.MainCategRepo
	SubCateg    interfaces.SubCategRepo
}

func NewTransactionUC(t interfaces.TransactionRepo, m interfaces.MainCategRepo, s interfaces.SubCategRepo) *TransactionUC {
	return &TransactionUC{
		Transaction: t,
		MainCateg:   m,
		SubCateg:    s,
	}
}

func (t *TransactionUC) Create(ctx context.Context, trans domain.CreateTransactionInput) error {
	// check if the main category exists
	mainCateg, err := t.MainCateg.GetByID(trans.MainCategID, trans.UserID)
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
	if err != nil {
		return err
	}

	// check if the sub category matches the main category
	if subCateg.MainCategID != trans.MainCategID {
		logger.Error("Create Transaction failed", "package", PackageName, "err", domain.ErrMainCategNotConsistent)
		return domain.ErrMainCategNotConsistent
	}

	return t.Transaction.Create(ctx, trans)
}

func (t *TransactionUC) GetAll(ctx context.Context, opt domain.GetTransOpt, user domain.User) ([]domain.Transaction, domain.Cursor, error) {
	trans, decodedNextKeys, err := t.Transaction.GetAll(ctx, opt, user.ID)
	if err != nil {
		return nil, domain.Cursor{}, err
	}

	var cursor domain.Cursor
	if opt.Cursor.Size != 0 && len(trans) == opt.Cursor.Size {
		cursor.Size = opt.Cursor.Size

		// if it's the first page, we need to initialize the nextKey
		// note that the order of decodedNextKeys does matter
		// the query will be like:
		// AND col_1 < or > val_1
		// OR (col_1 = val_1 AND col_2 < or > val_2)
		// the base field(ID)(col_2) should be the last field
		// so we need to make sure the ID is the last field in the decodedNextKeys
		if opt.Cursor.NextKey == "" {
			decodedNextKeys = domain.DecodedNextKeys{}

			if opt.Sort != nil && opt.Sort.By.IsValid() {
				decodedNextKeys = append(decodedNextKeys, domain.DecodedNextKeyInfo{
					Field: opt.Sort.By.GetField(),
				})
			}

			decodedNextKeys = append(decodedNextKeys, domain.DecodedNextKeyInfo{
				Field: "ID",
			})
		}

		// encode the nextKey to string
		encodedNextKey, err := codeutil.EncodeNextKeys(decodedNextKeys, trans[len(trans)-1])
		if err != nil {
			return nil, domain.Cursor{}, err
		}

		cursor.NextKey = encodedNextKey
	} else {
		cursor.NextKey = ""
		cursor.Size = 0
	}

	return trans, cursor, nil
}

func (t *TransactionUC) Update(ctx context.Context, trans domain.UpdateTransactionInput, user domain.User) error {
	// check if the main category exists
	mainCateg, err := t.MainCateg.GetByID(trans.MainCategID, user.ID)
	if err != nil {
		return err
	}

	// check if the type in main category matches the transaction type
	if trans.Type != mainCateg.Type {
		logger.Error("Update Transaction failed", "package", PackageName, "err", domain.ErrTypeNotConsistent)
		return domain.ErrTypeNotConsistent
	}

	// check if the sub category exists
	subCateg, err := t.SubCateg.GetByID(trans.SubCategID, user.ID)
	if err != nil {
		return err
	}

	// check if the sub category matches the main category
	if trans.MainCategID != subCateg.MainCategID {
		logger.Error("Update Transaction failed", "package", PackageName, "err", domain.ErrMainCategNotConsistent)
		return domain.ErrMainCategNotConsistent
	}

	// check permission
	if _, err := t.Transaction.GetByIDAndUserID(ctx, trans.ID, user.ID); err != nil {
		return err
	}

	return t.Transaction.Update(ctx, trans)
}

func (t *TransactionUC) Delete(ctx context.Context, id int64, user domain.User) error {
	// check permission
	if _, err := t.Transaction.GetByIDAndUserID(ctx, id, user.ID); err != nil {
		return err
	}

	return t.Transaction.Delete(ctx, id)
}

func (t *TransactionUC) GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, user domain.User) (domain.AccInfo, error) {
	return t.Transaction.GetAccInfo(ctx, query, user.ID)
}

func (t *TransactionUC) GetBarChartData(ctx context.Context, chartDateRange domain.ChartDateRange, timeRangeType domain.TimeRangeType, transactionType domain.TransactionType, mainCategIDs []int64, user domain.User) (domain.ChartData, error) {
	var dateToData domain.DateToChartData
	var err error
	if timeRangeType.IsDailyType() {
		dateToData, err = t.Transaction.GetDailyBarChartData(ctx, chartDateRange, transactionType, mainCategIDs, user.ID)
		if err != nil {
			return domain.ChartData{}, err
		}
	} else {
		dateToData, err = t.Transaction.GetMonthlyBarChartData(ctx, chartDateRange, transactionType, mainCategIDs, user.ID)
		if err != nil {
			return domain.ChartData{}, err
		}
	}

	return genChartData(dateToData, timeRangeType, chartDateRange.Start, chartDateRange.End), nil
}

func (t *TransactionUC) GetPieChartData(ctx context.Context, chartDateRange domain.ChartDateRange, transactionType domain.TransactionType, user domain.User) (domain.ChartData, error) {
	return t.Transaction.GetPieChartData(ctx, chartDateRange, transactionType, user.ID)
}

func (t *TransactionUC) GetLineChartData(ctx context.Context, chartDateRange domain.ChartDateRange, timeRangeType domain.TimeRangeType, user domain.User) (domain.ChartData, error) {
	var dateToData domain.DateToChartData
	var err error
	if timeRangeType.IsDailyType() {
		dateToData, err = t.Transaction.GetDailyLineChartData(ctx, chartDateRange, user.ID)
		if err != nil {
			return domain.ChartData{}, err
		}
	} else {
		dateToData, err = t.Transaction.GetMonthlyLineChartData(ctx, chartDateRange, user.ID)
		if err != nil {
			return domain.ChartData{}, err
		}
	}

	return genLineChartData(dateToData, timeRangeType, chartDateRange.Start, chartDateRange.End), nil
}

func (t *TransactionUC) GetMonthlyData(ctx context.Context, dateRange domain.GetMonthlyDateRange, user domain.User) ([]domain.TransactionType, error) {
	data := make([]domain.TransactionType, 0, dateRange.EndDate.Day())

	monthlyData, err := t.Transaction.GetMonthlyData(ctx, dateRange, user.ID)
	if err != nil {
		return data, err
	}

	// loop from 1 to the last day of the month(30 or 31)
	// Note that it's important to start at index 1, not 0
	// because monthlyData contains the data from day 1, there's no data for day 0
	// inside the loop, we use `append` to help us to insert the data to the correct index
	for t := dateRange.StartDate; t.Before(dateRange.EndDate) || t.Equal(dateRange.EndDate); t = t.AddDate(0, 0, 1) {
		day := t.Day()

		if transactionType, ok := monthlyData[day]; ok {
			data = append(data, transactionType)
		} else {
			data = append(data, domain.TransactionTypeUnSpecified)
		}
	}

	return data, nil
}
