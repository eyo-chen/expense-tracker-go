package usericon

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/gofacto"
	"github.com/eyo-chen/gofacto/db/mysqlf"
)

type factory struct {
	ui *gofacto.Factory[userIcon]
	u  *gofacto.Factory[user.User]
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		ui: gofacto.New(userIcon{}).WithDB(mysqlf.NewConfig(db)),
		u:  gofacto.New(user.User{}).WithDB(mysqlf.NewConfig(db)),
	}
}

func (f *factory) InsertManyWithOneUser(ctx context.Context, i int, ow ...userIcon) ([]userIcon, user.User, error) {
	u := user.User{}
	userIcons, err := f.ui.BuildList(ctx, i).WithOne(&u).Overwrites(ow...).Insert()
	if err != nil {
		return nil, user.User{}, err
	}

	return userIcons, u, nil
}

func (f *factory) InsertUser(ctx context.Context) (user.User, error) {
	return f.u.Build(ctx).Insert()
}

func (f *factory) Reset() {
	f.ui.Reset()
}
