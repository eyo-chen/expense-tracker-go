package transaction

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
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
		transaction: gofacto.New(Transaction{}).WithDB(mysqlf.NewConfig(db)).WithBlueprint(blueprint),
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

func (tf *TransactionFactory) PrepareUserMainAndSubCateg(ctx context.Context) (user.User, maincateg.MainCateg, subcateg.SubCateg, error) {
	u := user.User{}
	m := maincateg.MainCateg{Type: domain.TransactionTypeExpense.ToModelValue(), IconType: domain.IconTypeDefault.ToModelValue()}

	s, err := tf.subcateg.Build(ctx).WithOne(&u, &m).Insert()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, err
	}

	return u, m, s, nil
}

func (tf *TransactionFactory) InsertTransactionsWithOneUser(ctx context.Context, i int, ow ...Transaction) ([]Transaction, user.User, []maincateg.MainCateg, []subcateg.SubCateg, error) {
	u := user.User{}
	maincategOW := maincateg.MainCateg{Type: domain.TransactionTypeExpense.ToModelValue(), IconType: domain.IconTypeDefault.ToModelValue()}
	maincategAnyList := typeconv.ToAnysWithOW[maincateg.MainCateg](i, &maincategOW)
	subcategAnyList := typeconv.ToAnysWithOW[subcateg.SubCateg](i, nil)

	transList, err := tf.transaction.
		BuildList(ctx, i).
		WithOne(&u).
		WithMany(maincategAnyList).
		WithMany(subcategAnyList).
		Overwrites(ow...).
		Insert()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, err
	}

	maincategList := typeconv.ToT[maincateg.MainCateg](maincategAnyList)
	subcategList := typeconv.ToT[subcateg.SubCateg](subcategAnyList)

	return transList, u, maincategList, subcategList, nil
}

// InsertMainCategList inserts a list of main categories
func (tf *TransactionFactory) InsertMainCategList(ctx context.Context, i int, ow ...maincateg.MainCateg) ([]maincateg.MainCateg, user.User, error) {
	u := user.User{}

	maincategList, err := tf.maincateg.BuildList(ctx, i).Overwrites(ow...).WithOne(&u).Insert()
	if err != nil {
		return nil, user.User{}, err
	}

	return maincategList, u, nil
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
