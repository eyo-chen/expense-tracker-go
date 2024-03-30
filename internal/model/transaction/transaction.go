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
	tr := cvtCreateTransInputToModelTransaction(trans)
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

func (t *TransactionModel) Update(ctx context.Context, trans domain.UpdateTransactionInput) error {
	tr := cvtUpdateTransInputToModelTransaction(trans)
	qStmt := "UPDATE transactions SET type = ?, main_category_id = ?, sub_category_id = ?, price = ?, note = ?, date = ? WHERE id = ?"

	if _, err := t.DB.ExecContext(ctx, qStmt, tr.Type, tr.MainCategID, tr.SubCategID, tr.Price, tr.Note, tr.Date, tr.ID); err != nil {
		logger.Error("t.DB.ExecContext failed", "package", PackageName, "err", err)
		return err
	}

	return nil
}

func (t *TransactionModel) Delete(ctx context.Context, id int64) error {
	qStmt := "DELETE FROM transactions WHERE id = ?"

	if _, err := t.DB.ExecContext(ctx, qStmt, id); err != nil {
		logger.Error("t.DB.ExecContext failed", "package", PackageName, "err", err)
		return err
	}

	return nil
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

func (t *TransactionModel) GetBarChartData(ctx context.Context, dataRange domain.ChartDateRange, transactionType domain.TransactionType, userID int64) (domain.ChartDataByWeekday, error) {
	qStmt := `
	  SELECT DATE_FORMAT(date, '%a'),
		       SUM(price)
		FROM transactions
		WHERE user_id = ?
		AND type = ?
		AND date BETWEEN ? AND ?
		GROUP BY date
	`

	rows, err := t.DB.QueryContext(ctx, qStmt, userID, transactionType.ToModelValue(), dataRange.StartDate, dataRange.EndDate)
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

func (t *TransactionModel) GetPieChartData(ctx context.Context, dataRange domain.ChartDateRange, transactionType domain.TransactionType, userID int64) (domain.ChartData, error) {
	qStmt := `
	  SELECT mc.name,
		       SUM(ts.price)
		FROM transactions AS ts
		INNER JOIN main_categories AS mc
		ON ts.main_category_id = mc.id
		WHERE ts.user_id = ?
		AND ts.type = ?
		AND ts.date BETWEEN ? AND ?
		GROUP BY mc.name
	`

	rows, err := t.DB.QueryContext(ctx, qStmt, userID, transactionType.ToModelValue(), dataRange.StartDate, dataRange.EndDate)
	if err != nil {
		logger.Error("t.DB.QueryContext failed", "package", PackageName, "err", err)
		return domain.ChartData{}, err
	}

	var labels []string
	var datasets []float64
	for rows.Next() {
		var name string
		var price float64
		if err := rows.Scan(&name, &price); err != nil {
			logger.Error("rows.Scan failed", "package", PackageName, "err", err)
			return domain.ChartData{}, err
		}

		labels = append(labels, name)
		datasets = append(datasets, price)
	}
	defer rows.Close()

	return domain.ChartData{Labels: labels, Datasets: datasets}, nil
}

func (t *TransactionModel) GetMonthlyData(ctx context.Context, dateRange domain.GetMonthlyDateRange, userID int64) domain.MonthDayToTransactionType {
	qStmt := `
		SELECT
		DAY(date) AS day,
		CASE
			WHEN COUNT(DISTINCT type) = 1 AND MAX(type) = 1 THEN 1
			WHEN COUNT(DISTINCT type) = 1 AND MAX(type) = 2 THEN 2
		ELSE 3
		END AS type
		FROM transactions
		WHERE user_id = ?
		AND date BETWEEN ? AND ?
		GROUP BY DAY(date)
	`

	rows, err := t.DB.QueryContext(ctx, qStmt, userID, dateRange.StartDate, dateRange.EndDate)
	if err != nil {
		logger.Error("t.DB.QueryContext failed", "package", PackageName, "err", err)
		return nil
	}
	defer rows.Close()

	data := domain.MonthDayToTransactionType{}
	for rows.Next() {
		var date string
		var t domain.TransactionType
		if err := rows.Scan(&date, &t); err != nil {
			logger.Error("rows.Scan failed", "package", PackageName, "err", err)
			return nil
		}

		data[date] = t
	}

	return data
}
