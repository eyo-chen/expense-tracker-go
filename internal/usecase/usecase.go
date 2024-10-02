package usecase

import (
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/initdata"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/user"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/usericon"
)

type Usecase struct {
	User        *user.UC
	MainCateg   *maincateg.UC
	SubCateg    *subcateg.UC
	Transaction *transaction.UC
	Icon        *icon.UC
	UserIcon    *usericon.UC
	InitData    *initdata.UC
}

func New(u interfaces.UserRepo,
	m interfaces.MainCategRepo,
	s interfaces.SubCategRepo,
	i interfaces.IconRepo,
	t interfaces.TransactionRepo,
	r interfaces.RedisService,
	ui interfaces.UserIconRepo,
	s3 interfaces.S3Service,
) *Usecase {
	return &Usecase{
		User:        user.New(u, r),
		MainCateg:   maincateg.New(m, i, ui, r, s3),
		SubCateg:    subcateg.New(s, m),
		Transaction: transaction.New(t, m, s, r, s3),
		Icon:        icon.New(i, ui, r, s3),
		UserIcon:    usericon.New(s3, ui),
		InitData:    initdata.New(i, m, s, u),
	}
}
