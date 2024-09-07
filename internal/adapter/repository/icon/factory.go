package icon

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/gofacto"
	"github.com/eyo-chen/gofacto/db/mysqlf"
)

type factory struct {
	i *gofacto.Factory[Icon]
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		i: gofacto.New(Icon{}).WithDB(mysqlf.NewConfig(db)),
	}
}

func (f *factory) InsertMany(ctx context.Context, i int, ow ...Icon) ([]Icon, error) {
	return f.i.BuildList(ctx, i).Overwrites(ow...).Insert()
}

func (f *factory) Reset() {
	f.i.Reset()
}
