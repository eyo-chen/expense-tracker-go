package maincateg

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
)

type MainCategFactory struct {
	mainCateg *testutil.Factory[MainCateg]
	user      *testutil.Factory[user.User]
	icon      *testutil.Factory[icon.Icon]
}

func NewMainCategFactory(db *sql.DB) *MainCategFactory {
	return &MainCategFactory{
		mainCateg: testutil.NewFactory[MainCateg](db, MainCateg{}, BluePrintMainCategory, InserterMainCategory),
		user:      testutil.NewFactory[user.User](db, user.User{}, BluePrintUser, InserterUser),
		icon:      testutil.NewFactory[icon.Icon](db, icon.Icon{}, BluePrintIcon, InsertIcon),
	}
}

func (mf *MainCategFactory) PrepareUsers(i int, overwrites ...*user.User) *MainCategFactory {
	mf.user.BuildList(i).Overwrites(overwrites)
	return mf
}

func (mf *MainCategFactory) PrepareIcons(i int, overwrites ...*icon.Icon) *MainCategFactory {
	mf.icon.BuildList(i).Overwrites(overwrites)
	return mf
}

func (mf *MainCategFactory) InsertUserAndIcon() ([]user.User, []icon.Icon, error) {
	users, err := mf.user.InsertList()
	if err != nil {
		return nil, nil, err
	}

	icons, err := mf.icon.InsertList()
	if err != nil {
		return nil, nil, err
	}

	return users, icons, nil
}

func (mf *MainCategFactory) PrepareMainCateies(i int, overwrites ...*MainCateg) *MainCategFactory {
	mf.mainCateg.BuildList(i).Overwrites(overwrites)
	return mf
}

func (mf *MainCategFactory) InsertMainCateies() ([]MainCateg, error) {
	return mf.mainCateg.InsertList()
}

func (mf *MainCategFactory) PrepareMainCateg(overwrite *MainCateg) *MainCategFactory {
	mf.mainCateg.Build().Overwrite(overwrite)
	return mf
}

func (mf *MainCategFactory) InsertMainCateg() (MainCateg, error) {
	return mf.mainCateg.Insert()
}

func (mf *MainCategFactory) Reset() {
	mf.mainCateg.Reset()
	mf.user.Reset()
	mf.icon.Reset()
}
