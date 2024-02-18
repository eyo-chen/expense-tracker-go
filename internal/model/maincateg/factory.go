package maincateg

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
)

type MainCategFactory struct {
	mainCateg *testutil.Factory[MainCateg]
	user      *testutil.Factory[user.User]
	icon      *testutil.Factory[icon.Icon]
}

func setIncomeType(maincateg *MainCateg) {
	maincateg.Type = domain.Income.ModelValue()
}

func setExpenseType(maincateg *MainCateg) {
	maincateg.Type = domain.Expense.ModelValue()
}

func NewMainCategFactory(db *sql.DB) *MainCategFactory {
	return &MainCategFactory{
		mainCateg: testutil.NewFactory[MainCateg](db, MainCateg{}, BluePrint, Inserter).SetTrait("income", setIncomeType).SetTrait("expense", setExpenseType),
		user:      testutil.NewFactory[user.User](db, user.User{}, BluePrintUser, InserterUser),
		icon:      testutil.NewFactory[icon.Icon](db, icon.Icon{}, BluePrintIcon, InsertIcon),
	}
}

func (mf *MainCategFactory) PrepareUsers(i int, overwrites ...user.User) *MainCategFactory {
	mf.user.BuildList(i).Overwrites(overwrites)
	return mf
}

func (mf *MainCategFactory) PrepareIcons(i int, overwrites ...icon.Icon) *MainCategFactory {
	mf.icon.BuildList(i).Overwrites(overwrites)
	return mf
}

func (mf *MainCategFactory) InsertIcons() ([]icon.Icon, error) {
	return mf.icon.InsertList()
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

func (mf *MainCategFactory) PrepareMainCateies(i int, overwrites ...MainCateg) *MainCategFactory {
	mf.mainCateg.BuildList(i).Overwrites(overwrites)
	return mf
}

func (mf *MainCategFactory) InsertMainCateies() ([]MainCateg, error) {
	return mf.mainCateg.InsertList()
}

func (mf *MainCategFactory) PrepareMainCateg(overwrite MainCateg) *MainCategFactory {
	mf.mainCateg.Build().Overwrite(overwrite)
	return mf
}

func (mf *MainCategFactory) InsertMainCateg() (MainCateg, error) {
	return mf.mainCateg.Insert()
}

func (mf *MainCategFactory) InsertMainCategWithAss(ow MainCateg) (MainCateg, user.User, icon.Icon, error) {
	user := &user.User{}
	icon := &icon.Icon{}

	maincateg, _, err := mf.mainCateg.Build().Overwrite(ow).WithOne(user).WithOne(icon).InsertWithAss()

	return maincateg, *user, *icon, err
}

func (mf *MainCategFactory) InsertMainCategListWithAss(i int, userIdx int, iconIdx int, traitName ...string) ([]MainCateg, []user.User, []icon.Icon, error) {
	iconPtrList := make([]interface{}, 0, iconIdx)
	for k := 0; k < iconIdx; k++ {
		iconPtrList = append(iconPtrList, &icon.Icon{})
	}

	userPtrList := make([]interface{}, 0, userIdx)
	for k := 0; k < userIdx; k++ {
		userPtrList = append(userPtrList, &user.User{})
	}

	maincategList, _, err := mf.mainCateg.BuildList(i).WithTraits(traitName).WithMany(userIdx, userPtrList...).WithMany(iconIdx, iconPtrList...).InsertListWithAss()
	if err != nil {
		return nil, nil, nil, err
	}

	icons := make([]icon.Icon, 0, i)
	for _, v := range iconPtrList {
		icons = append(icons, *v.(*icon.Icon))
	}

	users := make([]user.User, 0, i)
	for _, v := range userPtrList {
		users = append(users, *v.(*user.User))
	}

	return maincategList, users, icons, nil
}

func (mf *MainCategFactory) Reset() {
	mf.mainCateg.Reset()
	mf.user.Reset()
	mf.icon.Reset()
}
