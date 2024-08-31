package usecase

import (
	"github.com/eyo-chen/expense-tracker-go/internal/model/interfaces"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/initdata"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/user"
)

type Usecase struct {
	User        user.UserUC
	MainCateg   maincateg.MainCategUC
	SubCateg    subcateg.SubCategUC
	Transaction transaction.TransactionUC
	Icon        icon.IconUC
	InitData    initdata.InitDataUC
}

func New(u interfaces.UserModel,
	m interfaces.MainCategModel,
	s interfaces.SubCategModel,
	i interfaces.IconModel,
	t interfaces.TransactionModel,
	r interfaces.RedisService,
) *Usecase {
	return &Usecase{
		User:        *user.NewUserUC(u),
		MainCateg:   *maincateg.NewMainCategUC(m, i),
		SubCateg:    *subcateg.NewSubCategUC(s, m),
		Transaction: *transaction.NewTransactionUC(t, m, s),
		Icon:        *icon.NewIconUC(i, r),
		InitData:    *initdata.NewInitDataUC(i, m, s, u),
	}
}
