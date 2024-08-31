package maincateg

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/model/user"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory/utils"
	"github.com/eyo-chen/gofacto"
	"github.com/eyo-chen/gofacto/db/mysqlf"
	"github.com/eyo-chen/gofacto/typeconv"
)

type factory struct {
	MainCateg *gofacto.Factory[MainCateg]
	User      *gofacto.Factory[user.User]
	Icon      *gofacto.Factory[icon.Icon]
}

func setIncomeType(m *MainCateg) {
	m.Type = domain.TransactionTypeIncome.ToModelValue()
}

func setExpenseType(m *MainCateg) {
	m.Type = domain.TransactionTypeExpense.ToModelValue()
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		MainCateg: gofacto.New(MainCateg{}).WithDB(mysqlf.NewConfig(db)).
			WithStorageName("main_categories").
			WithTrait("income", setIncomeType).
			WithTrait("expense", setExpenseType),
		User: gofacto.New(user.User{}).WithDB(mysqlf.NewConfig(db)),
		Icon: gofacto.New(icon.Icon{}).WithDB(mysqlf.NewConfig(db)),
	}
}

// InsertUsersAndIcons inserts many users and icons
func (mf *factory) InsertUsersAndIcons(ctx context.Context, userI int, iconI int) ([]user.User, []icon.Icon, error) {
	users, err := mf.User.BuildList(ctx, userI).Insert()
	if err != nil {
		return nil, nil, err
	}

	icons, err := mf.Icon.BuildList(ctx, iconI).Insert()
	if err != nil {
		return nil, nil, err
	}

	return users, icons, nil
}

// InsertMainCateg inserts a main category
func (mf *factory) InsertMainCategWithAss(ctx context.Context, ow MainCateg) (MainCateg, user.User, icon.Icon, error) {
	user := &user.User{}
	icon := &icon.Icon{}

	maincateg, err := mf.MainCateg.Build(ctx).
		Overwrite(MainCateg{Type: domain.TransactionTypeIncome.ToModelValue()}). // set default type to income
		Overwrite(ow).
		WithOne(user).
		WithOne(icon).
		Insert()

	return maincateg, *user, *icon, err
}

// InsertMainCategList inserts many main categories with associations and traits
func (mf *factory) InsertMainCategListWithAss(ctx context.Context, i int, userIdx int, iconIdx int, traitName ...string) ([]MainCateg, []user.User, []icon.Icon, error) {
	iconPtrList := typeconv.ToAnysWithOW[icon.Icon](iconIdx, nil)
	userPtrList := typeconv.ToAnysWithOW[user.User](userIdx, nil)

	maincategList, err := mf.MainCateg.BuildList(ctx, i).
		Overwrite(MainCateg{Type: domain.TransactionTypeIncome.ToModelValue()}). // set default type to income
		SetTraits(traitName...).
		WithMany(userPtrList).
		WithMany(iconPtrList).
		Insert()
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
