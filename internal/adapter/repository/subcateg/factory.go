package subcateg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/gofacto"
	"github.com/eyo-chen/gofacto/db/mysqlf"
	"github.com/eyo-chen/gofacto/typeconv"
)

type factory struct {
	subCateg  *gofacto.Factory[SubCateg]
	maincateg *gofacto.Factory[maincateg.MainCateg]
}

func newFactory(db *sql.DB) *factory {
	return &factory{
		subCateg: gofacto.New(SubCateg{}).
			WithDB(mysqlf.NewConfig(db)).
			WithStorageName("sub_categories"),
		maincateg: gofacto.New(maincateg.MainCateg{}).
			WithDB(mysqlf.NewConfig(db)).
			WithStorageName("main_categories").
			WithBlueprint(maincateg.Blueprint),
	}
}

// InsertUserAndMaincateg inserts one user and one main category.
func (f *factory) InsertUserAndMaincateg(ctx context.Context) (user.User, maincateg.MainCateg, error) {
	u := user.User{}
	ow := maincateg.MainCateg{Type: domain.TransactionTypeExpense.ToModelValue()}

	m, err := f.maincateg.Build(ctx).WithOne(&u).WithOne(&icon.Icon{}).Overwrite(ow).Insert()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, err
	}

	return u, m, nil
}

// InsertSubcategsWithOneOrManyMainCateg inserts n subcategories with one or many main categories and one user.
// The maincategIndex is the number of main categories to insert.
// The subcategIndexes is the number of subcategories to insert for each main category.
func (f *factory) InsertSubcategsWithOneOrManyMainCateg(ctx context.Context, maincategIndex int, subcategIndexes []int, ows ...maincateg.MainCateg) (map[int64][]SubCateg, []maincateg.MainCateg, user.User, error) {
	// validate input
	// maincategIndex and subcategIndexes must have the same length
	// because the value of subcategIndexes[i] is the number of subcategories to insert for the main category with index i.
	if maincategIndex != len(subcategIndexes) {
		return nil, nil, user.User{}, errors.New("maincategIndex and subcategIndexes must have the same length")
	}

	// prepare main categories overwrites
	maincategOWs := make([]maincateg.MainCateg, maincategIndex)
	for i := 0; i < maincategIndex; i++ {
		if ows != nil && i < len(ows) {
			maincategOWs[i] = ows[i]
		}

		// set default type to expense if not set
		if maincategOWs[i].Type == "" {
			maincategOWs[i].Type = domain.TransactionTypeExpense.ToModelValue()
		}
	}

	// prepare associations data
	u := user.User{}
	iconPtrList := typeconv.ToAnysWithOW[icon.Icon](maincategIndex, nil)

	// insert main categories
	ms, err := f.maincateg.BuildList(ctx, maincategIndex).Overwrites(maincategOWs...).WithOne(&u).WithMany(iconPtrList).Insert()
	if err != nil {
		return nil, nil, user.User{}, err
	}

	// prepare sub categories overwrites
	subcategOWs := make([][]SubCateg, len(subcategIndexes))
	for i := 0; i < len(subcategIndexes); i++ {
		subcategOWs[i] = make([]SubCateg, subcategIndexes[i])
		for j := 0; j < subcategIndexes[i]; j++ {
			subcategOWs[i][j].MainCategID = ms[i].ID
			subcategOWs[i][j].UserID = u.ID
		}
	}

	// insert sub categories, and prepare result
	result := map[int64][]SubCateg{}
	for index, amount := range subcategIndexes {
		ows := subcategOWs[index]
		ss, err := f.subCateg.BuildList(ctx, amount).Overwrites(ows...).Insert()
		if err != nil {
			return nil, nil, user.User{}, err
		}

		mainCateg := ms[index]
		result[mainCateg.ID] = ss
	}

	return result, ms, u, nil
}

func (f *factory) Reset() {
	f.subCateg.Reset()
	f.maincateg.Reset()
}
