package maincateg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
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

func Blueprint(i int) MainCateg {
	return MainCateg{
		Name:     fmt.Sprintf("test%d", i),
		Type:     domain.TransactionTypeIncome.ToModelValue(),
		IconType: domain.IconTypeDefault.ToModelValue(),
		IconData: "url",
	}
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		MainCateg: gofacto.New(MainCateg{}).WithDB(mysqlf.NewConfig(db)).
			WithStorageName("main_categories").
			WithTrait("income", setIncomeType).
			WithTrait("expense", setExpenseType).
			WithBlueprint(Blueprint),
		User: gofacto.New(user.User{}).WithDB(mysqlf.NewConfig(db)),
		Icon: gofacto.New(icon.Icon{}).WithDB(mysqlf.NewConfig(db)),
	}
}

// InsertUsers inserts many users
func (mf *factory) InsertUsers(ctx context.Context, userI int) ([]user.User, error) {
	users, err := mf.User.BuildList(ctx, userI).Insert()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// InsertMainCateg inserts a main category
func (mf *factory) InsertMainCategWithAss(ctx context.Context, ow MainCateg) (MainCateg, user.User, error) {
	user := &user.User{}

	maincateg, err := mf.MainCateg.Build(ctx).
		Overwrite(ow).
		WithOne(user).
		Insert()

	return maincateg, *user, err
}

// InsertMainCategList inserts many main categories with associations and traits
func (mf *factory) InsertMainCategListWithAss(ctx context.Context, i int, userIdx int, iconIdx int, traitName ...string) ([]MainCateg, []user.User, error) {
	userPtrList := typeconv.ToAnysWithOW[user.User](userIdx, nil)

	maincategList, err := mf.MainCateg.BuildList(ctx, i).
		SetTraits(traitName...).
		WithMany(userPtrList).
		Insert()
	if err != nil {
		return nil, nil, err
	}

	users := typeconv.ToT[user.User](userPtrList)
	return maincategList, users, nil
}

// Reset resets the factory
func (mf *factory) Reset() {
	mf.MainCateg.Reset()
	mf.User.Reset()
	mf.Icon.Reset()
}
