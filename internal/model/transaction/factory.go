package transaction

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/model/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/model/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/model/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/model/user"
	"github.com/eyo-chen/gofacto"
	"github.com/eyo-chen/gofacto/db/mysqlf"
	"github.com/eyo-chen/gofacto/typeconv"
)

type TransactionFactory struct {
	transaction *gofacto.Factory[Transaction]
	user        *gofacto.Factory[user.User]
	maincateg   *gofacto.Factory[maincateg.MainCateg]
	subcateg    *gofacto.Factory[subcateg.SubCateg]
}

func NewTransactionFactory(db *sql.DB) *TransactionFactory {
	return &TransactionFactory{
		transaction: gofacto.New(Transaction{}).WithDB(mysqlf.NewConfig(db)).WithBlueprint(BluePrint),
		user:        gofacto.New(user.User{}).WithDB(mysqlf.NewConfig(db)),
		maincateg: gofacto.New(maincateg.MainCateg{}).
			WithDB(mysqlf.NewConfig(db)).
			WithStorageName("main_categories").
			WithBlueprint(maincateg.Blueprint),
		subcateg: gofacto.New(subcateg.SubCateg{}).
			WithDB(mysqlf.NewConfig(db)).
			WithStorageName("sub_categories"),
	}
}

func (tf *TransactionFactory) PrepareUserMainAndSubCateg(ctx context.Context) (user.User, maincateg.MainCateg, subcateg.SubCateg, icon.Icon, error) {
	u := user.User{}
	i := icon.Icon{}
	m, err := tf.maincateg.Build(ctx).WithOne(&u).WithOne(&i).Insert()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, icon.Icon{}, err
	}

	ow := subcateg.SubCateg{UserID: u.ID, MainCategID: m.ID}
	s, err := tf.subcateg.Build(ctx).Overwrite(ow).Insert()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, icon.Icon{}, err
	}

	return u, m, s, i, nil
}

func (tf *TransactionFactory) InsertTransactionsWithOneUser(ctx context.Context, i int, ow ...Transaction) ([]Transaction, user.User, []maincateg.MainCateg, []subcateg.SubCateg, []icon.Icon, error) {
	u := user.User{}

	iconPtrList := typeconv.ToAnysWithOW[icon.Icon](i, nil)

	maincategList, err := tf.maincateg.BuildList(ctx, i).WithOne(&u).WithMany(iconPtrList).Insert()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	owSub := []subcateg.SubCateg{}
	for _, m := range maincategList {
		owSub = append(owSub, subcateg.SubCateg{UserID: m.UserID, MainCategID: m.ID})
	}

	subcategList, err := tf.subcateg.BuildList(ctx, i).Overwrites(owSub...).Insert()
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

	transList, err := tf.transaction.BuildList(ctx, i).Overwrites(owTrans...).Overwrites(ow...).Insert()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	iconList := typeconv.ToT[icon.Icon](iconPtrList)
	return transList, u, maincategList, subcategList, iconList, nil
}

// InsertMainCategList inserts a list of main categories
func (tf *TransactionFactory) InsertMainCategList(ctx context.Context, i int, ow ...maincateg.MainCateg) ([]maincateg.MainCateg, user.User, []icon.Icon, error) {
	u := user.User{}

	iconPtrList := typeconv.ToAnysWithOW[icon.Icon](i, nil)
	maincategList, err := tf.maincateg.BuildList(ctx, i).Overwrites(ow...).WithOne(&u).WithMany(iconPtrList).Insert()
	if err != nil {
		return nil, user.User{}, []icon.Icon{}, err
	}

	iconList := typeconv.ToT[icon.Icon](iconPtrList)
	return maincategList, u, iconList, nil
}

// InsertTransactionWithGivenUser inserts a transaction with a given user
// it assumes that the user has main category
func (tf *TransactionFactory) InsertTransactionWithGivenUser(ctx context.Context, i int, u user.User, ow ...Transaction) ([]Transaction, subcateg.SubCateg, error) {
	// create only one sub category
	owSub := subcateg.SubCateg{UserID: u.ID, MainCategID: ow[0].MainCategID}
	s, err := tf.subcateg.Build(ctx).Overwrite(owSub).Insert()
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

	transList, err := tf.transaction.BuildList(ctx, i).Overwrites(owTrans...).Overwrites(ow...).Insert()
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
