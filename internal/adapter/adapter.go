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
	User        user.Repo
	MainCateg   maincateg.Repo
	SubCateg    subcateg.Repo
	Icon        icon.Repo
	Transaction transaction.Repo
}

func New(mysqlDB *sql.DB) *Model {
	return &Model{
		User:        *user.New(mysqlDB),
		MainCateg:   *maincateg.New(mysqlDB),
		SubCateg:    *subcateg.New(mysqlDB),
		Icon:        *icon.New(mysqlDB),
		Transaction: *transaction.New(mysqlDB),
	}
}
