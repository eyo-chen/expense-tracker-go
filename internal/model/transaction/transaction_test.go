package transaction_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/transaction"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx     = context.Background()
	mockLoc, _  = time.LoadLocation("")
	mockTimeNow = time.Unix(1629446406, 0).Truncate(24 * time.Hour).In(mockLoc)
)

type TransactionSuite struct {
	suite.Suite
	db               *sql.DB
	migrate          *migrate.Migrate
	transactionModel *transaction.TransactionModel
	f                *transaction.TransactionFactory
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db, migrate := testutil.ConnToDB(port)
	logger.Register()

	s.db = db
	s.migrate = migrate
	s.transactionModel = transaction.NewTransactionModel(s.db)
	s.f = transaction.NewTransactionFactory(db)
}

func (s *TransactionSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *TransactionSuite) SetupTest() {
	s.transactionModel = transaction.NewTransactionModel(s.db)
	s.f = transaction.NewTransactionFactory(s.db)
}

func (s *TransactionSuite) TearDownTest() {
	tx, err := s.db.Begin()
	if err != nil {
		s.Require().NoError(err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM transactions"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM icons"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM sub_categories"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM main_categories"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM users"); err != nil {
		s.Require().NoError(err)
	}

	s.Require().NoError(tx.Commit())
	s.f.Reset()
}

func (s *TransactionSuite) TestCreate() {
	user, main, sub, _, err := s.f.PrepareUserMainAndSubCateg()
	s.Require().NoError(err)

	t := domain.CreateTransactionInput{
		UserID:      user.ID,
		Type:        domain.CvtToTransactionType(main.Type),
		MainCategID: main.ID,
		SubCategID:  sub.ID,
		Price:       100,
		Note:        "test",
		Date:        mockTimeNow,
	}

	err = s.transactionModel.Create(mockCtx, t)
	s.Require().NoError(err)

	var checkT transaction.Transaction
	stmt := "SELECT user_id, type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE user_id = ?"
	err = s.db.QueryRow(stmt, user.ID).Scan(&checkT.UserID, &checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Equal(t.UserID, checkT.UserID)
	s.Equal(t.Type.ToModelValue(), checkT.Type)
	s.Equal(t.MainCategID, checkT.MainCategID)
	s.Equal(t.SubCategID, checkT.SubCategID)
	s.Equal(t.Price, checkT.Price)
	s.Equal(t.Note, checkT.Note)
}

func (s *TransactionSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return successfully":            getAll_NoError_ReturnSuccessfully,
		"when with multiple users, return successfully": getAll_WithMultipleUsers_ReturnSuccessfully,
		"when with diff date, return successfully":      getAll_WithDiffDate_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAll_NoError_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 0)

	trans, err := s.transactionModel.GetAll(mockCtx, nil, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 0)

	trans, err := s.transactionModel.GetAll(mockCtx, nil, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_WithDiffDate_ReturnSuccessfully(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(3, ow1, ow2, ow3)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(3, ow1, ow2, ow3)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 0, 1)

	getQuery := &domain.GetQuery{
		StartDate: mockTimeNow.AddDate(0, 0, -3).Format("2006-01-02"),
		EndDate:   mockTimeNow.AddDate(0, 0, -2).Format("2006-01-02"),
	}

	trans, err := s.transactionModel.GetAll(mockCtx, getQuery, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)

}
