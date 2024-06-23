package maincateg

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory/db/esql"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory/utils"
)

type factory struct {
	MainCateg *efactory.Factory[MainCateg]
	User      *efactory.Factory[user.User]
	Icon      *efactory.Factory[icon.Icon]
}

func setIncomeType(m *MainCateg) {
	m.Type = domain.TransactionTypeIncome.ToModelValue()
}

func setExpenseType(m *MainCateg) {
	m.Type = domain.TransactionTypeExpense.ToModelValue()
}

func newFactory(db *sql.DB) *factory {
	categConfig := efactory.Config[MainCateg]{
		DB:          &esql.Config{DB: db},
		StorageName: "main_categories",
		BluePrint:   BluePrint,
	}

	userConfig := efactory.Config[user.User]{
		DB: &esql.Config{DB: db},
	}

	iconConfig := efactory.Config[icon.Icon]{
		DB: &esql.Config{DB: db},
	}

	return &factory{
		MainCateg: efactory.New(MainCateg{}).SetConfig(categConfig).
			SetTrait("income", setIncomeType).
			SetTrait("expense", setExpenseType),
		User: efactory.New(user.User{}).SetConfig(userConfig),
		Icon: efactory.New(icon.Icon{}).SetConfig(iconConfig),
	}
}

// InsertUsersAndIcons inserts many users and icons
func (mf *factory) InsertUsersAndIcons(userI int, iconI int) ([]user.User, []icon.Icon, error) {
	users, err := mf.User.BuildList(userI).Insert()
	if err != nil {
		return nil, nil, err
	}

	icons, err := mf.Icon.BuildList(iconI).Insert()
	if err != nil {
		return nil, nil, err
	}

	return users, icons, nil
}

// InsertMainCateg inserts a main category
func (mf *factory) InsertMainCategWithAss(ow MainCateg) (MainCateg, user.User, icon.Icon, error) {
	user := &user.User{}
	icon := &icon.Icon{}

	maincateg, _, err := mf.MainCateg.Build().Overwrite(ow).WithOne(user).WithOne(icon).InsertWithAss()

	return maincateg, *user, *icon, err
}

// InsertMainCategList inserts many main categories with associations and traits
func (mf *factory) InsertMainCategListWithAss(i int, userIdx int, iconIdx int, traitName ...string) ([]MainCateg, []user.User, []icon.Icon, error) {
	iconPtrList := utils.CvtToAnysWithOW[icon.Icon](iconIdx, nil)
	userPtrList := utils.CvtToAnysWithOW[user.User](userIdx, nil)

	maincategList, _, err := mf.MainCateg.BuildList(i).
		WithTraits(traitName...).
		WithMany(userPtrList...).
		WithMany(iconPtrList...).
		InsertWithAss()
	if err != nil {
		return nil, nil, nil, err
	}

	icons := utils.CvtToT[icon.Icon](iconPtrList)
	users := utils.CvtToT[user.User](userPtrList)
	return maincategList, users, icons, nil
}

// Reset resets the factory
func (mf *factory) Reset() {
	mf.MainCateg.Reset()
	mf.User.Reset()
	mf.Icon.Reset()
}
