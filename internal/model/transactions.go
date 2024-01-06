package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type TransactionModel struct {
	DB *sql.DB
}

type Transaction struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	MainCategID int64      `json:"main_category_id"`
	SubCategID  int64      `json:"sub_category_id"`
	Price       float64    `json:"price"`
	Note        string     `json:"note"`
	Date        *time.Time `json:"date"`
}

func newTransactionModel(db *sql.DB) *TransactionModel {
	return &TransactionModel{DB: db}
}

func (t *TransactionModel) Create(ctx context.Context, transaction *domain.Transaction) error {
	trans := cvtToModelTransaction(transaction)
	qStmt := "INSERT INTO transactions (user_id, main_category_id, sub_category_id, price, note, date) VALUES (?, ?, ?, ?, ?, ?)"

	if _, err := t.DB.ExecContext(ctx, qStmt, trans.UserID, trans.MainCategID, trans.SubCategID, trans.Price, trans.Note, trans.Date); err != nil {
		logger.Error("t.DB.ExecContext failed", "package", "model", "err", err)
		return err
	}

	return nil

}

func (t *TransactionModel) GetAll(ctx context.Context, query *domain.GetQuery, userID int64) (*domain.TransactionResp, error) {
	qStmt := getQStmt(query, userID)
	args := getArgs(query, userID)

	rows, err := t.DB.QueryContext(ctx, qStmt, args...)
	if err != nil {
		logger.Error("t.DB.QueryContext failed", "package", "model", "err", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	var income, expense float64
	for rows.Next() {
		var trans Transaction
		var mainCateg MainCateg
		var subCateg SubCateg
		var icon Icon

		if err := rows.Scan(&trans.ID, &trans.UserID, &trans.Price, &trans.Note, &trans.Date, &mainCateg.ID, &mainCateg.Name, &mainCateg.Type, &subCateg.ID, &subCateg.Name, &icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", "model", "err", err)
			return nil, err
		}

		if mainCateg.Type == "1" {
			income += trans.Price
		} else {
			expense += trans.Price
		}

		transactions = append(transactions, cvtToDomainTransaction(&trans, &mainCateg, &subCateg, &icon))
	}

	var result domain.TransactionResp
	result.DataList = transactions
	result.Income = income
	result.Expense = expense
	result.NetIncome = income - expense

	return &result, nil
}

func getQStmt(query *domain.GetQuery, userID int64) string {
	qStmt := `SELECT t.id, t.user_id, t.price, t.note, t.date, mc.id, mc.name, mc.type, sc.id, sc.name, i.id, i.url
						FROM transactions AS t
						INNER JOIN main_categories AS mc 
						ON t.main_category_id = mc.id
						INNER JOIN sub_categories AS sc 
						ON t.sub_category_id = sc.id
						INNER JOIN icons AS i
						ON mc.icon_id = i.id
						WHERE t.user_id = ?`

	if query.StartDate != "" && query.EndDate != "" {
		qStmt += " AND date BETWEEN ? AND ?"
	}

	if query.StartDate != "" {
		qStmt += " AND date >= ?"
	}

	if query.EndDate != "" {
		qStmt += " AND date <= ?"
	}

	return qStmt
}

func getArgs(query *domain.GetQuery, userID int64) []interface{} {
	var args []interface{}

	args = append(args, userID)

	if query.StartDate != "" && query.EndDate != "" {
		args = append(args, query.StartDate, query.EndDate)
	}

	if query.StartDate != "" {
		args = append(args, query.StartDate)
	}

	if query.EndDate != "" {
		args = append(args, query.EndDate)
	}

	return args
}
