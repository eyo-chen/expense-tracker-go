package transaction

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory/db/esql"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory/utils"
)

type TransactionFactory struct {
	transaction *efactory.Factory[Transaction]
	user        *efactory.Factory[user.User]
	maincateg   *efactory.Factory[maincateg.MainCateg]
	subcateg    *efactory.Factory[subcateg.SubCateg]
}

func NewTransactionFactory(db *sql.DB) *TransactionFactory {
	return &TransactionFactory{
		transaction: efactory.New(Transaction{}).SetConfig(efactory.Config[Transaction]{
			DB:        &esql.Config{DB: db},
			BluePrint: BluePrint,
		}),
		user: efactory.New(user.User{}).SetConfig(efactory.Config[user.User]{
			DB: &esql.Config{DB: db},
		}),
		maincateg: efactory.New(maincateg.MainCateg{}).SetConfig(efactory.Config[maincateg.MainCateg]{
			DB:          &esql.Config{DB: db},
			StorageName: "main_categories",
			BluePrint:   maincateg.BluePrint,
		}),
		subcateg: efactory.New(subcateg.SubCateg{}).SetConfig(efactory.Config[subcateg.SubCateg]{
			DB:          &esql.Config{DB: db},
			StorageName: "sub_categories",
		}),
	}
}

func (tf *TransactionFactory) PrepareUserMainAndSubCateg() (user.User, maincateg.MainCateg, subcateg.SubCateg, icon.Icon, error) {
	u := user.User{}
	i := icon.Icon{}
	m, _, err := tf.maincateg.Build().WithOne(&u).WithOne(&i).InsertWithAss()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, icon.Icon{}, err
	}

	ow := subcateg.SubCateg{UserID: u.ID, MainCategID: m.ID}
	s, err := tf.subcateg.Build().Overwrite(ow).Insert()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, icon.Icon{}, err
	}

	return u, m, s, i, nil
}

func (tf *TransactionFactory) InsertTransactionsWithOneUser(i int, ow ...Transaction) ([]Transaction, user.User, []maincateg.MainCateg, []subcateg.SubCateg, []icon.Icon, error) {
	u := user.User{}

	iconPtrList := utils.CvtToAnysWithOW[icon.Icon](i, nil)

	maincategList, _, err := tf.maincateg.BuildList(i).WithOne(&u).WithMany(iconPtrList...).InsertWithAss()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	owSub := []subcateg.SubCateg{}
	for _, m := range maincategList {
		owSub = append(owSub, subcateg.SubCateg{UserID: m.UserID, MainCategID: m.ID})
	}

	subcategList, err := tf.subcateg.BuildList(i).Overwrites(owSub...).Insert()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	owTrans := []Transaction{}
	for k, m := range maincategList {
		owTrans = append(owTrans, Transaction{
			UserID:      m.UserID,
			MainCategID: m.ID,
			SubCategID:  subcategList[k].ID,
		})
	}

	transList, err := tf.transaction.BuildList(i).Overwrites(owTrans...).Overwrites(ow...).Insert()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	iconList := utils.CvtToT[icon.Icon](iconPtrList)
	return transList, u, maincategList, subcategList, iconList, nil
}

// InsertMainCategList inserts a list of main categories
func (tf *TransactionFactory) InsertMainCategList(i int, ow ...maincateg.MainCateg) ([]maincateg.MainCateg, user.User, []icon.Icon, error) {
	u := user.User{}

	iconPtrList := utils.CvtToAnysWithOW[icon.Icon](i, nil)
	maincategList, _, err := tf.maincateg.BuildList(i).Overwrites(ow...).WithOne(&u).WithMany(iconPtrList...).InsertWithAss()
	if err != nil {
		return nil, user.User{}, []icon.Icon{}, err
	}

	iconList := utils.CvtToT[icon.Icon](iconPtrList)
	return maincategList, u, iconList, nil
}

// InsertTransactionWithGivenUser inserts a transaction with a given user
// it assumes that the user has main category
func (tf *TransactionFactory) InsertTransactionWithGivenUser(i int, u user.User, ow ...Transaction) ([]Transaction, subcateg.SubCateg, error) {
	// create only one sub category
	owSub := subcateg.SubCateg{UserID: u.ID, MainCategID: ow[0].MainCategID}
	s, err := tf.subcateg.Build().Overwrite(owSub).Insert()
	if err != nil {
		return []Transaction{}, subcateg.SubCateg{}, err
	}

	owTrans := make([]Transaction, i)
	for k := 0; k < i; k++ {
		owTrans[k] = Transaction{
			UserID:     u.ID,
			SubCategID: s.ID,
		}
	}

	transList, err := tf.transaction.BuildList(i).Overwrites(owTrans...).Overwrites(ow...).Insert()
	if err != nil {
		return []Transaction{}, subcateg.SubCateg{}, err
	}

	return transList, s, nil
}

func (tf *TransactionFactory) Reset() {
	tf.transaction.Reset()
	tf.user.Reset()
	tf.maincateg.Reset()
	tf.subcateg.Reset()
}
