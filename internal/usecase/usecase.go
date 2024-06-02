package usecase

import (
	"github.com/OYE0303/expense-tracker-go/internal/model/interfaces"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/icon"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/initdata"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/transaction"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/user"
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
) *Usecase {
	return &Usecase{
		User:        *user.NewUserUC(u),
		MainCateg:   *maincateg.NewMainCategUC(m, i),
		SubCateg:    *subcateg.NewSubCategUC(s, m),
		Transaction: *transaction.NewTransactionUC(t, m, s),
		Icon:        *icon.NewIconUC(i),
		InitData:    *initdata.NewInitDataUC(i, m, s, u),
	}
}
