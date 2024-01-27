package model

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/transaction"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
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
