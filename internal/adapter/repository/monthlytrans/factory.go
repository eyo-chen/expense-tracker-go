package monthlytrans

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/gofacto"
	"github.com/eyo-chen/gofacto/db/mysqlf"
)

type factory struct {
	user *gofacto.Factory[user.User]
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		user: gofacto.New(user.User{}).WithDB(mysqlf.NewConfig(db)),
	}
}

func (f *factory) InsertUsers(ctx context.Context, userI int) ([]user.User, error) {
	users, err := f.user.BuildList(ctx, userI).Insert()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (f *factory) Reset() {
	f.user.Reset()
}
