package subcateg

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory/db/esql"
)

type Factory struct {
	subCateg  *efactory.Factory[SubCateg]
	maincateg *efactory.Factory[maincateg.MainCateg]
	user      *efactory.Factory[user.User]
}

func NewFactory(db *sql.DB) *Factory {
	subcategConfig := efactory.Config[SubCateg]{
		DB:          &esql.Config{DB: db},
		StorageName: "sub_categories",
	}

	maincategConfig := efactory.Config[maincateg.MainCateg]{
		DB:          &esql.Config{DB: db},
		StorageName: "main_categories",
	}

	userConfig := efactory.Config[user.User]{
		DB: &esql.Config{DB: db},
	}

	return &Factory{
		subCateg:  efactory.New(SubCateg{}).SetConfig(subcategConfig),
		maincateg: efactory.New(maincateg.MainCateg{}).SetConfig(maincategConfig),
		user:      efactory.New(user.User{}).SetConfig(userConfig),
	}
}

func (f *Factory) InsertUserAndMaincateg() (user.User, maincateg.MainCateg, error) {
	u := user.User{}
	ow := maincateg.MainCateg{Type: domain.TransactionTypeExpense.ToModelValue()}

	m, _, err := f.maincateg.Build().WithOne(&u).WithOne(&icon.Icon{}).Overwrite(ow).InsertWithAss()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, err
	}

	return u, m, nil
}

func (f *Factory) InsertSubcategs(n int, ows ...SubCateg) ([]SubCateg, user.User, maincateg.MainCateg, error) {
	u := user.User{}
	ow := maincateg.MainCateg{Type: domain.TransactionTypeExpense.ToModelValue()}
	m, _, err := f.maincateg.Build().WithOne(&u).WithOne(&icon.Icon{}).Overwrite(ow).InsertWithAss()
	if err != nil {
		return nil, user.User{}, maincateg.MainCateg{}, err
	}

	subcategOWs := make([]SubCateg, n)
	for i := 0; i < n; i++ {
		if ows != nil {
			subcategOWs[i] = ows[i]
		}
		subcategOWs[i].MainCategID = m.ID
		subcategOWs[i].UserID = u.ID
	}

	ss, err := f.subCateg.BuildList(n).Overwrites(subcategOWs...).Insert()
	if err != nil {
		return nil, user.User{}, maincateg.MainCateg{}, err
	}

	return ss, u, m, nil
}

func (f *Factory) Reset() {
	f.subCateg.Reset()
	f.user.Reset()
	f.maincateg.Reset()
}
