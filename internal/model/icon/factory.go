package icon

import (
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory/db/esql"
)

type factory struct {
	i *efactory.Factory[Icon]
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		i: efactory.New(Icon{}).SetConfig(efactory.Config[Icon]{
			DB: &esql.Config{DB: db},
		}),
	}
}

func (f *factory) InsertMany(i int, ow ...Icon) ([]Icon, error) {
	return f.i.BuildList(i).Overwrites(ow...).Insert()
}

func (f *factory) Reset() {
	f.i.Reset()
}
