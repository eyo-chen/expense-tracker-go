package monthlytrans

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/gofacto"
	"github.com/eyo-chen/gofacto/db/mysqlf"
)

type factory struct {
	user         *gofacto.Factory[user.User]
	monthlyTrans *gofacto.Factory[MonthlyTrans]
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		user: gofacto.New(user.User{}).WithDB(mysqlf.NewConfig(db)),
		monthlyTrans: gofacto.New(MonthlyTrans{}).
			WithDB(mysqlf.NewConfig(db)).
			WithStorageName("monthly_transactions"),
	}
}

func (f *factory) InsertUsers(ctx context.Context, userI int) ([]user.User, error) {
	users, err := f.user.BuildList(ctx, userI).Insert()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (f *factory) InsertManyMonthlyTransWithOneUser(ctx context.Context, i int, ows []MonthlyTrans) (user.User, []MonthlyTrans, error) {
	user := user.User{}
	mt, err := f.monthlyTrans.BuildList(ctx, i).
		Overwrites(ows...).
		WithOne(&user).
		Insert()
	if err != nil {
		return user, nil, err
	}

	return user, mt, nil
}

func (f *factory) Reset() {
	f.user.Reset()
}
