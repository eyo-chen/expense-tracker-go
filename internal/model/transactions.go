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
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      int64              `bson:"user_id"`
	Type        string             `bson:"type"`
	MainCategID int64              `bson:"main_category_id"`
	SubCategID  int64              `bson:"sub_category_id"`
	Price       int64              `bson:"price"`
	Date        *time.Time         `bson:"date"`
	Note        string             `bson:"note,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
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
	startDate, endDate, err := parseDate(query.StartDate, query.EndDate)
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"user_id": userID,
		"date": map[string]interface{}{
			"$gte": primitive.NewDateTimeFromTime(*startDate),
			"$lte": primitive.NewDateTimeFromTime(*endDate),
		},
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

func parseDate(startDate, endDate string) (*time.Time, *time.Time, error) {
	parsedStartDate, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return nil, nil, err
	}

	parsedEndDate, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return nil, nil, err
	}

	return &parsedStartDate, &parsedEndDate, nil
}
