package handler

import (
	"github.com/eyo-chen/expense-tracker-go/internal/handler/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/initdata"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/user"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type Handler struct {
	User        user.Hlr
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
