package transaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "model/transaction"
)

type TransactionModel struct {
	DB *sql.DB
}

type Transaction struct {
	ID          int64
	Type        string
	UserID      int64 `efactory:"User"`
	MainCategID int64 `efactory:"MainCateg,main_categories" esql:"main_category_id"`
	SubCategID  int64 `efactory:"SubCateg,sub_categories" esql:"sub_category_id"`
	Price       float64
	Note        string
	Date        time.Time
}

func NewTransactionModel(db *sql.DB) *TransactionModel {
	return &TransactionModel{DB: db}
}

func (t *TransactionModel) Create(ctx context.Context, trans domain.CreateTransactionInput) error {
	tr := cvtToModelTransaction(trans)
	qStmt := "INSERT INTO transactions (user_id, type, main_category_id, sub_category_id, price, note, date) VALUES (?, ?, ?, ?, ?, ?, ?)"

	if _, err := t.DB.ExecContext(ctx, qStmt, tr.UserID, tr.Type, tr.MainCategID, tr.SubCategID, tr.Price, tr.Note, tr.Date); err != nil {
		logger.Error("t.DB.ExecContext failed", "package", PackageName, "err", err)
		return err
	}

	return nil
}

func (t *TransactionModel) GetAll(ctx context.Context, query domain.GetQuery, userID int64) ([]domain.Transaction, error) {
	qStmt := getAllQStmt(query)
	args := getAllArgs(query, userID)

	rows, err := t.DB.QueryContext(ctx, qStmt, args...)
	if err != nil {
		logger.Error("t.DB.QueryContext failed", "package", PackageName, "err", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var trans Transaction
		var mainCateg maincateg.MainCateg
		var subCateg subcateg.SubCateg
		var icon icon.Icon

		if err := rows.Scan(&trans.ID, &trans.UserID, &trans.Type, &trans.Price, &trans.Note, &trans.Date, &mainCateg.ID, &mainCateg.Name, &mainCateg.Type, &subCateg.ID, &subCateg.Name, &icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", PackageName, "err", err)
			return nil, err
		}

		transactions = append(transactions, cvtToDomainTransaction(trans, mainCateg, subCateg, icon))
	}
	defer rows.Close()

	return transactions, nil
}

func (t *TransactionModel) GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, userID int64) (domain.AccInfo, error) {
	qStmt := getAccInfoQStmt(query)
	args := getAccInfoArgs(query, userID)

	var accInfo domain.AccInfo
	if err := t.DB.QueryRowContext(ctx, qStmt, args...).
		Scan(&accInfo.TotalIncome, &accInfo.TotalExpense, &accInfo.TotalBalance); err != nil && err != sql.ErrNoRows {
		logger.Error("t.DB.QueryRowContext failed", "package", PackageName, "err", err)
		return domain.AccInfo{}, err
	}

	return accInfo, nil
}

func (t *TransactionModel) Delete(ctx context.Context, id int64) error {
	qStmt := "DELETE FROM transactions WHERE id = ?"

	if _, err := t.DB.ExecContext(ctx, qStmt, id); err != nil {
		logger.Error("t.DB.ExecContext failed", "package", PackageName, "err", err)
		return err
	}

	return nil
}

func (t *TransactionModel) GetByIDAndUserID(ctx context.Context, id, userID int64) (domain.Transaction, error) {
	qStmt := "SELECT id, user_id, type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE id = ? AND user_id = ?"

	var trans Transaction
	if err := t.DB.QueryRowContext(ctx, qStmt, id, userID).
		Scan(&trans.ID, &trans.UserID, &trans.Type, &trans.MainCategID, &trans.SubCategID, &trans.Price, &trans.Note, &trans.Date); err != nil {
		if err == sql.ErrNoRows {
			return domain.Transaction{}, domain.ErrTransactionDataNotFound
		}

		logger.Error("t.DB.QueryRowContext failed", "package", PackageName, "err", err)
		return domain.Transaction{}, err
	}

	return cvtToDomainTransactionWithoutCategory(trans), nil
}

func (t *TransactionModel) GetChartData(ctx context.Context, chartType domain.ChartType, dataRange domain.ChartDateRange, userID int64) (domain.ChartDataByWeekday, error) {
	qStmt := `
	  SELECT DATE_FORMAT(date, '%a'),
		       SUM(price)
		FROM transactions
		WHERE user_id = ?
		AND type = 2
		AND date BETWEEN ? AND ?
		GROUP BY date
	`

	rows, err := t.DB.QueryContext(ctx, qStmt, userID, dataRange.StartDate, dataRange.EndDate)
	if err != nil {
		logger.Error("t.DB.QueryContext failed", "package", PackageName, "err", err)
		return domain.ChartDataByWeekday{}, err
	}

	dataByTime := domain.ChartDataByWeekday{}
	for rows.Next() {
		var date string
		var price float64
		if err := rows.Scan(&date, &price); err != nil {
			logger.Error("rows.Scan failed", "package", PackageName, "err", err)
			return domain.ChartDataByWeekday{}, err
		}

		// Normally, this case will never happen
		if v, ok := dataByTime[date]; ok {
			dataByTime[date] = v + price
		} else {
			dataByTime[date] = price
		}
	}
	defer rows.Close()

	return dataByTime, nil
}
