package handler

import (
	"github.com/OYE0303/expense-tracker-go/internal/handler/icon"
	"github.com/OYE0303/expense-tracker-go/internal/handler/initdata"
	"github.com/OYE0303/expense-tracker-go/internal/handler/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/handler/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/handler/transaction"
	"github.com/OYE0303/expense-tracker-go/internal/handler/user"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
)

type Handler struct {
	User        user.UserHandler
	MainCateg   maincateg.MainCategHandler
	SubCateg    subcateg.SubCategHandler
	Transaction transaction.TransactionHandler
	Icon        icon.IconHandler
	InitData    initdata.InitDataHlr
}

func New(u interfaces.UserUC,
	m interfaces.MainCategUC,
	s interfaces.SubCategUC,
	t interfaces.TransactionUC,
	i interfaces.IconUC,
	in interfaces.InitDataUC,
) *Handler {
	return &Handler{
		User:        *user.NewUserHandler(u),
		MainCateg:   *maincateg.NewMainCategHandler(m),
		SubCateg:    *subcateg.NewSubCategHandler(s),
		Transaction: *transaction.NewTransactionHandler(t),
		Icon:        *icon.NewIconHandler(i),
		InitData:    *initdata.New(in),
	}
}
