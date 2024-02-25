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
	ID          int64     `json:"id"`
	Type        string    `json:"type"`
	UserID      int64     `json:"user_id" factory:"User,users"`
	MainCategID int64     `json:"main_category_id" factory:"MainCateg,main_categories"`
	SubCategID  int64     `json:"sub_category_id" factory:"SubCateg,sub_categories"`
	Price       float64   `json:"price"`
	Note        string    `json:"note"`
	Date        time.Time `json:"date"`
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
	qStmt := getQStmt(query, userID)
	args := getArgs(query, userID)

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

	return transactions, nil
}

func (t *TransactionModel) GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, userID int64) (domain.AccInfo, error) {
	qStmt := getAccInfoQStmt(query, userID)
	args := getAccInfoArgs(query, userID)

	var accInfo domain.AccInfo
	if err := t.DB.QueryRowContext(ctx, qStmt, args...).
		Scan(&accInfo.TotalIncome, &accInfo.TotalExpense, &accInfo.TotalBalance); err != nil && err != sql.ErrNoRows {
		logger.Error("t.DB.QueryRowContext failed", "package", PackageName, "err", err)
		return domain.AccInfo{}, err
	}

	return accInfo, nil
}
