package usecase

import (
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/initdata"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/user"
)

type Usecase struct {
	User        *user.UC
	MainCateg   *maincateg.UC
	SubCateg    *subcateg.UC
	Transaction *transaction.UC
	Icon        *icon.UC
	InitData    *initdata.UC
}

func New(u interfaces.UserRepo,
	m interfaces.MainCategRepo,
	s interfaces.SubCategRepo,
	i interfaces.IconRepo,
	t interfaces.TransactionRepo,
	r interfaces.RedisService,
	ui interfaces.UserIconRepo,
) *Usecase {
	return &Usecase{
		User:        user.New(u, r),
		MainCateg:   maincateg.New(m, i),
		SubCateg:    subcateg.New(s, m),
		Transaction: transaction.New(t, m, s),
		Icon:        icon.New(i, ui, r, nil),
		InitData:    initdata.New(i, m, s, u),
	}
}
