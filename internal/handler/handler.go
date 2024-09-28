package handler

import (
	"github.com/eyo-chen/expense-tracker-go/internal/handler/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/initdata"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/interfaces"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/user"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/usericon"
)

type Handler struct {
	User        *user.Hlr
	MainCateg   *maincateg.Hlr
	SubCateg    *subcateg.Hlr
	Transaction *transaction.Hlr
	Icon        *icon.Hlr
	UserIcon    *usericon.Hlr
	InitData    *initdata.Hlr
}

func New(u interfaces.UserUC,
	m interfaces.MainCategUC,
	s interfaces.SubCategUC,
	t interfaces.TransactionUC,
	i interfaces.IconUC,
	ui interfaces.UserIconUC,
	in interfaces.InitDataUC,
) *Handler {
	return &Handler{
		User:        user.New(u),
		MainCateg:   maincateg.New(m),
		SubCateg:    subcateg.New(s),
		Transaction: transaction.New(t),
		Icon:        icon.New(i),
		UserIcon:    usericon.New(ui),
		InitData:    initdata.New(in),
	}
}
