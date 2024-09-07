package model

import (
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
)

type Model struct {
	User        user.UserModel
	MainCateg   maincateg.MainCategModel
	SubCateg    subcateg.SubCategModel
	Icon        icon.IconModel
	Transaction transaction.TransactionModel
}

func New(mysqlDB *sql.DB) *Model {
	return &Model{
		User:        *user.NewUserModel(mysqlDB),
		MainCateg:   *maincateg.NewMainCategModel(mysqlDB),
		SubCateg:    *subcateg.NewSubCategModel(mysqlDB),
		Icon:        *icon.NewIconModel(mysqlDB),
		Transaction: *transaction.NewTransactionModel(mysqlDB),
	}
}
