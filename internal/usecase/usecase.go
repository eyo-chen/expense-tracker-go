package usecase

import (
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/interfaces"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/initdata"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/user"
)

type Usecase struct {
	User        *user.UC
	MainCateg   *maincateg.MainCategUC
	SubCateg    *subcateg.SubCategUC
	Transaction *transaction.TransactionUC
	Icon        *icon.IconUC
	InitData    *initdata.InitDataUC
}

func New(u interfaces.UserRepo,
	m interfaces.MainCategRepo,
	s interfaces.SubCategRepo,
	i interfaces.IconRepo,
	t interfaces.TransactionRepo,
	r interfaces.RedisService,
) *Usecase {
	return &Usecase{
		User:        user.New(u, r),
		MainCateg:   maincateg.NewMainCategUC(m, i),
		SubCateg:    subcateg.NewSubCategUC(s, m),
		Transaction: transaction.NewTransactionUC(t, m, s),
		Icon:        icon.NewIconUC(i, r),
		InitData:    initdata.NewInitDataUC(i, m, s, u),
	}
}
