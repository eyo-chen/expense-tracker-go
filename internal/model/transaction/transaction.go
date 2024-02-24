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
		logger.Error("t.DB.ExecContext failed", "package", "model", "err", err)
		return err
	}

	return nil

}

func (t *TransactionModel) GetAll(ctx context.Context, query *domain.GetQuery, userID int64) ([]domain.Transaction, error) {
	qStmt := getQStmt(query, userID)
	args := getArgs(query, userID)

	rows, err := t.DB.QueryContext(ctx, qStmt, args...)
	if err != nil {
		logger.Error("t.DB.QueryContext failed", "package", "model", "err", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var trans Transaction
		var mainCateg maincateg.MainCateg
		var subCateg subcateg.SubCateg
		var icon icon.Icon

		if err := rows.Scan(&trans.ID, &trans.UserID, &trans.Price, &trans.Note, &trans.Date, &mainCateg.ID, &mainCateg.Name, &mainCateg.Type, &subCateg.ID, &subCateg.Name, &icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", "model", "err", err)
			return nil, err
		}

		transactions = append(transactions, cvtToDomainTransaction(trans, mainCateg, subCateg, icon))
	}

	return transactions, nil
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

	if query == nil {
		return qStmt
	}

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

	if query == nil {
		return args
	}

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
