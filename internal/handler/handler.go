package handler

import (
	"github.com/OYE0303/expense-tracker-go/internal/handler/interfaces"
	"github.com/OYE0303/expense-tracker-go/internal/handler/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/handler/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/handler/transaction"
	"github.com/OYE0303/expense-tracker-go/internal/handler/user"
)

type Handler struct {
	User        user.UserHandler
	MainCateg   maincateg.MainCategHandler
	SubCateg    subcateg.SubCategHandler
	Transaction transaction.TransactionHandler
}

func New(u interfaces.UserUC, m interfaces.MainCategUC, s interfaces.SubCategUC, t interfaces.TransactionUC) *Handler {
	return &Handler{
		User:        *user.NewUserHandler(u),
		MainCateg:   *maincateg.NewMainCategHandler(m),
		SubCateg:    *subcateg.NewSubCategHandler(s),
		Transaction: *transaction.NewTransactionHandler(t),
	}
}
