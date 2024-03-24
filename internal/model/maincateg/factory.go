package maincateg

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory/db/esql"
)

type MainCategFactory struct {
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

func NewMainCategFactory(db *sql.DB) *MainCategFactory {
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

	return &MainCategFactory{
		MainCateg: efactory.New(MainCateg{}).SetConfig(categConfig).
			SetTrait("income", setIncomeType).
			SetTrait("expense", setExpenseType),
		User: efactory.New(user.User{}).SetConfig(userConfig),
		Icon: efactory.New(icon.Icon{}).SetConfig(iconConfig),
	}
}

func (mf *MainCategFactory) InsertUserAndIcon(userI int, iconI int) ([]user.User, []icon.Icon, error) {
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

func (mf *MainCategFactory) InsertMainCategWithAss(ow MainCateg) (MainCateg, user.User, icon.Icon, error) {
	user := &user.User{}
	icon := &icon.Icon{}

	maincateg, _, err := mf.MainCateg.Build().Overwrite(ow).WithOne(user).WithOne(icon).InsertWithAss()

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

	maincategList, _, err := mf.MainCateg.BuildList(i).
		WithTraits(traitName...).
		WithMany(userPtrList...).WithMany(iconPtrList...).
		InsertWithAss()
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
	mf.MainCateg.Reset()
	mf.User.Reset()
	mf.Icon.Reset()
}
