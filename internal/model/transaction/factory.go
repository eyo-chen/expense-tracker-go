package transaction

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
)

type TransactionFactory struct {
	transaction *testutil.Factory[Transaction]
	user        *testutil.Factory[user.User]
	maincateg   *testutil.Factory[maincateg.MainCateg]
	subcateg    *testutil.Factory[subcateg.SubCateg]
}

func NewTransactionFactory(db *sql.DB) *TransactionFactory {
	return &TransactionFactory{
		transaction: testutil.NewFactory(db, Transaction{}, BluePrint, Inserter),
		user:        testutil.NewFactory(db, user.User{}, user.Blueprint, user.Inserter),
		maincateg:   testutil.NewFactory(db, maincateg.MainCateg{}, maincateg.BluePrint, maincateg.Inserter),
		subcateg:    testutil.NewFactory(db, subcateg.SubCateg{}, subcateg.Blueprint, subcateg.Inserter),
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

	iconPtrList := make([]interface{}, 0, i)
	for k := 0; k < i; k++ {
		iconPtrList = append(iconPtrList, &icon.Icon{})
	}

	maincategList, _, err := tf.maincateg.BuildList(i).WithOne(&u).WithMany(i, iconPtrList...).InsertListWithAss()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	owSub := []subcateg.SubCateg{}
	for _, m := range maincategList {
		owSub = append(owSub, subcateg.SubCateg{UserID: m.UserID, MainCategID: m.ID})
	}

	subcategList, err := tf.subcateg.BuildList(i).Overwrites(owSub).InsertList()
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

	transList, err := tf.transaction.BuildList(i).Overwrites(owTrans).Overwrites(ow).InsertList()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	iconList := make([]icon.Icon, 0, i)
	for _, v := range iconPtrList {
		iconList = append(iconList, *v.(*icon.Icon))
	}

	return transList, u, maincategList, subcategList, iconList, nil
}

func (tf *TransactionFactory) Reset() {
	tf.transaction.Reset()
	tf.user.Reset()
	tf.maincateg.Reset()
	tf.subcateg.Reset()
}
