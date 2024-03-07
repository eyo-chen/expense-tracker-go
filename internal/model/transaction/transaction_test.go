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

	s.TearDownTest()
}

func (s *TransactionSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return successfully":                      getAll_NoError_ReturnSuccessfully,
		"when with multiple users, return successfully":           getAll_WithMultipleUsers_ReturnSuccessfully,
		"when with many transactions, return all transactions":    getAll_WithManyTransaction_ReturnSuccessfully,
		"when query start date, return data after start date":     getAll_QueryStartDate_ReturnDataAfterStartDate,
		"when query end date, return data before end date":        getAll_QueryEndDate_ReturnDataBeforeEndDate,
		"when query start and end date, return data between them": getAll_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate,
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

	query := domain.GetQuery{}
	trans, err := s.transactionModel.GetAll(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 0)

	query := domain.GetQuery{}
	trans, err := s.transactionModel.GetAll(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_WithManyTransaction_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(3)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 0, 1, 2)

	query := domain.GetQuery{}

	trans, err := s.transactionModel.GetAll(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_QueryStartDate_ReturnDataAfterStartDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 1, 2, 3)

	getQuery := domain.GetQuery{
		StartDate: mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly),
	}
	trans, err := s.transactionModel.GetAll(mockCtx, getQuery, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_QueryEndDate_ReturnDataBeforeEndDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 0, 1)

	getQuery := domain.GetQuery{
		EndDate: mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly),
	}
	trans, err := s.transactionModel.GetAll(mockCtx, getQuery, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 1, 2)

	getQuery := domain.GetQuery{
		StartDate: mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly),
		EndDate:   mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly),
	}
	trans, err := s.transactionModel.GetAll(mockCtx, getQuery, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func (s *TransactionSuite) TestGetAccInfo() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return successfully":                                  getAccInfo_NoError_ReturnSuccessfully,
		"when with multiple users, return data only with one user":            getAccInfo_WithMultipleUsers_ReturnDataOnlyWithOneUser,
		"when with many transaction, return correct calculation":              getAccInfo_WithManyTransaction_ReturnCorrectCalculation,
		"when query start date, return accumulated data after start date":     getAccInfo_QueryStartDate_ReturnDataAfterStartDate,
		"when query end date, return accumulated data before end date":        getAccInfo_QueryEndDate_ReturnDataBeforeEndDate,
		"when query start and end date, return accumulated data between them": getAccInfo_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAccInfo_NoError_ReturnSuccessfully(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.Expense.ToModelValue()}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(1, ow1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 999,
		TotalIncome:  0,
		TotalBalance: -999,
	}

	query := domain.GetAccInfoQuery{}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_WithMultipleUsers_ReturnDataOnlyWithOneUser(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.Expense.ToModelValue()}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(1, ow1)
	s.Require().NoError(err, desc)

	ow2 := transaction.Transaction{Price: 1000, Type: domain.Expense.ToModelValue()}
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1, ow2)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 999,
		TotalIncome:  0,
		TotalBalance: -999,
	}

	query := domain.GetAccInfoQuery{}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_WithManyTransaction_ReturnCorrectCalculation(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.Expense.ToModelValue()}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.Expense.ToModelValue()}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.Income.ToModelValue()}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(3, ow1, ow2, ow3)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 1999,
		TotalIncome:  1000,
		TotalBalance: -999,
	}

	query := domain.GetAccInfoQuery{}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryStartDate_ReturnDataAfterStartDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.Expense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.Expense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.Income.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Price: 2000, Type: domain.Income.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 0,
		TotalIncome:  3000,
		TotalBalance: 3000,
	}

	query := domain.GetAccInfoQuery{
		StartDate: mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly),
	}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryEndDate_ReturnDataBeforeEndDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.Expense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.Expense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.Income.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Price: 2000, Type: domain.Income.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 1999,
		TotalIncome:  1000,
		TotalBalance: -999,
	}

	query := domain.GetAccInfoQuery{
		EndDate: mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly),
	}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.Expense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.Expense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.Income.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Price: 2000, Type: domain.Income.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 1000,
		TotalIncome:  1000,
		TotalBalance: 0,
	}

	query := domain.GetAccInfoQuery{
		StartDate: mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly),
		EndDate:   mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly),
	}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}
