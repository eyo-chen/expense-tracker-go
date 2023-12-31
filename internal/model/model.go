package model

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

type Model struct {
	User        UserModel
	MainCateg   MainCategModel
	SubCateg    SubCategModel
	Icon        IconModel
	Transaction TransactionModel
}

func New(mysqlDB *sql.DB, mongoDB *mongo.Database) *Model {
	return &Model{
		User:        *newUserModel(mysqlDB),
		MainCateg:   *newMainCategModel(mysqlDB),
		SubCateg:    *newSubCategModel(mysqlDB),
		Icon:        *newIconModel(mysqlDB),
		Transaction: *newTransactionModel(mongoDB),
	}
}
