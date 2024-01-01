package model

import (
	"context"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionTransactions = "transactions"
)

type TransactionModel struct {
	DB *mongo.Database
}

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    int64              `bson:"user_id"`
	Type      string             `bson:"type"`
	MainCateg *MainCateg         `bson:"main_category"`
	SubCateg  *SubCateg          `bson:"sub_category"`
	Price     int64              `bson:"price"`
	Date      *time.Time         `bson:"date"`
	Note      string             `bson:"note,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func newTransactionModel(db *mongo.Database) *TransactionModel {
	return &TransactionModel{DB: db}
}

func (t *TransactionModel) Create(ctx context.Context, transaction *domain.Transaction) error {
	curTime := time.Now()
	trans := cvtToModelTransaction(transaction)
	trans.CreatedAt = curTime
	trans.UpdatedAt = curTime

	_, err := t.DB.Collection(CollectionTransactions).InsertOne(ctx, trans)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionModel) GetAll(ctx context.Context, query *domain.GetQuery, userID int64) (*domain.TransactionResp, error) {
	filter, err := getFilter(query, userID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: 1}})

	cursor, err := t.DB.Collection(CollectionTransactions).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var transactions []*Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return cvtToDomainTransactionResp(transactions), nil
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
