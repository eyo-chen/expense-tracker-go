package subcateg

import (
	"database/sql"
	"errors"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/model/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/model/user"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory/db/esql"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory/utils"
)

type factory struct {
	subCateg  *efactory.Factory[SubCateg]
	maincateg *efactory.Factory[maincateg.MainCateg]
}

func newFactory(db *sql.DB) *factory {
	subcategConfig := efactory.Config[SubCateg]{
		DB:          &esql.Config{DB: db},
		StorageName: "sub_categories",
	}

	maincategConfig := efactory.Config[maincateg.MainCateg]{
		DB:          &esql.Config{DB: db},
		StorageName: "main_categories",
	}

	return &factory{
		subCateg:  efactory.New(SubCateg{}).SetConfig(subcategConfig),
		maincateg: efactory.New(maincateg.MainCateg{}).SetConfig(maincategConfig),
	}
}

// InsertUserAndMaincateg inserts one user and one main category.
func (f *factory) InsertUserAndMaincateg() (user.User, maincateg.MainCateg, error) {
	u := user.User{}
	ow := maincateg.MainCateg{Type: domain.TransactionTypeExpense.ToModelValue()}

	m, _, err := f.maincateg.Build().WithOne(&u).WithOne(&icon.Icon{}).Overwrite(ow).InsertWithAss()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, err
	}

	return u, m, nil
}

// InsertSubcategsWithOneOrManyMainCateg inserts n subcategories with one or many main categories and one user.
// The maincategIndex is the number of main categories to insert.
// The subcategIndexes is the number of subcategories to insert for each main category.
func (f *factory) InsertSubcategsWithOneOrManyMainCateg(maincategIndex int, subcategIndexes []int, ows ...maincateg.MainCateg) (map[int64][]SubCateg, []maincateg.MainCateg, user.User, error) {
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
	iconPtrList := utils.CvtToAnysWithOW[icon.Icon](maincategIndex, nil)

	// insert main categories
	ms, _, err := f.maincateg.BuildList(maincategIndex).Overwrites(maincategOWs...).WithOne(&u).WithMany(iconPtrList...).InsertWithAss()
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
		ss, err := f.subCateg.BuildList(amount).Overwrites(ows...).Insert()
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
