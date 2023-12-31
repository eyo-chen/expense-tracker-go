package model

import (
	"context"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	_, err := t.DB.Collection("transactions").InsertOne(ctx, trans)
	if err != nil {
		return err
	}
	return nil
}

func cvtToModelTransaction(t *domain.Transaction) *Transaction {
	return &Transaction{
		UserID:      t.UserID,
		Type:        cvtToModelType(t.Type),
		MainCategID: t.MainCategID,
		SubCategID:  t.SubCategID,
		Price:       t.Price,
		Date:        t.Date,
		Note:        t.Note,
	}
}
