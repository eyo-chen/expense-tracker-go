package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionModel struct {
	DB *sql.DB
}

type Transaction struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	MainCategID int64      `json:"main_category_id"`
	SubCategID  int64      `json:"sub_category_id"`
	Price       int64      `json:"price"`
	Note        string     `json:"note"`
	Date        *time.Time `json:"date"`
}

func newTransactionModel(db *sql.DB) *TransactionModel {
	return &TransactionModel{DB: db}
}

func (t *TransactionModel) Create(ctx context.Context, transaction *domain.Transaction) error {
	trans := cvtToModelTransaction(transaction)
	query := "INSERT INTO transactions (user_id, main_category_id, sub_category_id, price, note, date) VALUES (?, ?, ?, ?, ?, ?)"

	if _, err := t.DB.ExecContext(ctx, query, trans.UserID, trans.MainCategID, trans.SubCategID, trans.Price, trans.Note, trans.Date); err != nil {
		logger.Error("t.DB.ExecContext failed", "package", "model", "err", err)
		return err
	}

	return nil

}

func (t *TransactionModel) GetAll(ctx context.Context, query *domain.GetQuery, userID int64) (*domain.TransactionResp, error) {
	// filter, err := getFilter(query, userID)
	// if err != nil {
	// 	return nil, err
	// }

	// opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: 1}})

	// cursor, err := t.DB.Collection(CollectionTransactions).Find(ctx, filter, opts)
	// if err != nil {
	// 	return nil, err
	// }

	// var transactions []*Transaction
	// if err := cursor.All(ctx, &transactions); err != nil {
	// 	return nil, err
	// }

	// return cvtToDomainTransactionResp(transactions), nil

	return nil, nil
}

func getFilter(query *domain.GetQuery, userID int64) (bson.M, error) {
	dateFilter, err := getDateFilter(query)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"user_id": userID,
		"date":    dateFilter,
	}

	return filter, nil
}

func getDateFilter(query *domain.GetQuery) (bson.M, error) {
	filter := bson.M{}

	if query.StartDate != "" {
		startDate, err := parseDateOnly(query.StartDate)
		if err != nil {
			return nil, err
		}
		filter["$gte"] = primitive.NewDateTimeFromTime(*startDate)
	}

	if query.EndDate != "" {
		endDate, err := parseDateOnly(query.EndDate)
		if err != nil {
			return nil, err
		}
		filter["$lte"] = primitive.NewDateTimeFromTime(*endDate)
	}

	return filter, nil
}

func parseDateOnly(date string) (*time.Time, error) {
	parsedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, err
	}

	return &parsedDate, nil
}
