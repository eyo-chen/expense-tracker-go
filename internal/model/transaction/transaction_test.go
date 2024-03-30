package transaction_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/transaction"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx       = context.Background()
	mockLoc, _    = time.LoadLocation("")
	mockTimeNow   = time.Unix(1629446406, 0).Truncate(24 * time.Hour).In(mockLoc)
	weekDayFormat = "Mon"
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
		"when no error, return successfully":                             getAll_NoError_ReturnSuccessfully,
		"when with multiple users, return successfully":                  getAll_WithMultipleUsers_ReturnSuccessfully,
		"when with many transactions, return all transactions":           getAll_WithManyTransaction_ReturnSuccessfully,
		"when query start date, return data after start date":            getAll_QueryStartDate_ReturnDataAfterStartDate,
		"when query end date, return data before end date":               getAll_QueryEndDate_ReturnDataBeforeEndDate,
		"when query start and end date, return data between them":        getAll_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate,
		"when query main category id, return data with main category id": getAll_QueryMainCategID_ReturnDataWithMainCategID,
		"when query sub category id, return data with sub category id":   getAll_QuerySubCategID_ReturnDataWithSubCategID,
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

	startDate := mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly)
	getQuery := domain.GetQuery{
		StartDate: &startDate,
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

	endDate := mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly)
	getQuery := domain.GetQuery{
		EndDate: &endDate,
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

	startDate := mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly)
	endDate := mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly)
	getQuery := domain.GetQuery{
		StartDate: &startDate,
		EndDate:   &endDate,
	}
	trans, err := s.transactionModel.GetAll(mockCtx, getQuery, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_QueryMainCategID_ReturnDataWithMainCategID(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 0)

	getQuery := domain.GetQuery{
		MainCategID: &mainList[0].ID,
	}
	trans, err := s.transactionModel.GetAll(mockCtx, getQuery, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getAll_QuerySubCategID_ReturnDataWithSubCategID(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, iconList, err := s.f.InsertTransactionsWithOneUser(4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := transaction.GetAll_GenExpResult(transactionList, user, mainList, subList, iconList, 1)

	getQuery := domain.GetQuery{
		SubCategID: &subList[1].ID,
	}
	trans, err := s.transactionModel.GetAll(mockCtx, getQuery, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func (s *TransactionSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, update successfully":       update_WithOneData_UpdateSuccessfully,
		"when with multiple data, update successfully":  update_WithMultipleData_UpdateSuccessfully,
		"when with multiple users, update successfully": update_WithMultipleUsers_UpdateSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_WithOneData_UpdateSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	t := domain.UpdateTransactionInput{
		ID:          transactions[0].ID,
		Type:        domain.CvtToTransactionType(transactions[0].Type),
		MainCategID: transactions[0].MainCategID,
		SubCategID:  transactions[0].SubCategID,
		Price:       999,
		Note:        "update note",
		Date:        mockTimeNow.AddDate(0, 0, -1),
	}

	err = s.transactionModel.Update(mockCtx, t)
	s.Require().NoError(err, desc)

	var checkT transaction.Transaction
	stmt := "SELECT type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Require().Equal(t.Type.ToModelValue(), checkT.Type)
	s.Require().Equal(t.MainCategID, checkT.MainCategID)
	s.Require().Equal(t.SubCategID, checkT.SubCategID)
	s.Require().Equal(t.Price, checkT.Price)
	s.Require().Equal(t.Note, checkT.Note)
	s.Require().Equal(t.Date, checkT.Date)
}

func update_WithMultipleData_UpdateSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(2)
	s.Require().NoError(err, desc)

	t := domain.UpdateTransactionInput{
		ID:          transactions[0].ID,
		Type:        domain.CvtToTransactionType(transactions[0].Type),
		MainCategID: transactions[0].MainCategID,
		SubCategID:  transactions[0].SubCategID,
		Price:       999,
		Note:        "update note",
		Date:        mockTimeNow.AddDate(0, 0, -1),
	}

	err = s.transactionModel.Update(mockCtx, t)
	s.Require().NoError(err, desc)

	var checkT transaction.Transaction
	stmt := "SELECT type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Require().Equal(t.Type.ToModelValue(), checkT.Type)
	s.Require().Equal(t.Price, checkT.Price)
	s.Require().Equal(t.Note, checkT.Note)
	s.Require().Equal(t.Date, checkT.Date)

	// check if other data is not updated
	var checkT2 transaction.Transaction
	err = s.db.QueryRow(stmt, transactions[1].ID).Scan(&checkT2.Type, &checkT2.MainCategID, &checkT2.SubCategID, &checkT2.Price, &checkT2.Note, &checkT2.Date)
	s.Require().NoError(err)
	s.Require().Equal(transactions[1].Type, checkT2.Type)
	s.Require().Equal(transactions[1].Price, checkT2.Price)
	s.Require().Equal(transactions[1].Note, checkT2.Note)
	s.Require().Equal(transactions[1].Date, checkT2.Date)
}

func update_WithMultipleUsers_UpdateSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	// prepare more users
	transactions2, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	t := domain.UpdateTransactionInput{
		ID:          transactions[0].ID,
		Type:        domain.CvtToTransactionType(transactions[0].Type),
		MainCategID: transactions[0].MainCategID,
		SubCategID:  transactions[0].SubCategID,
		Price:       999,
		Note:        "update note",
		Date:        mockTimeNow,
	}

	err = s.transactionModel.Update(mockCtx, t)
	s.Require().NoError(err, desc)

	var checkT transaction.Transaction
	stmt := "SELECT type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Require().Equal(t.Type.ToModelValue(), checkT.Type)
	s.Require().Equal(t.Price, checkT.Price)
	s.Require().Equal(t.Note, checkT.Note)
	s.Require().Equal(t.Date, checkT.Date)

	// check if other user's data is not updated
	var checkT2 transaction.Transaction
	err = s.db.QueryRow(stmt, transactions2[0].ID).Scan(&checkT2.Type, &checkT2.MainCategID, &checkT2.SubCategID, &checkT2.Price, &checkT2.Note, &checkT2.Date)
	s.Require().NoError(err)
	s.Require().Equal(transactions2[0].Type, checkT2.Type)
	s.Require().Equal(transactions2[0].Price, checkT2.Price)
	s.Require().Equal(transactions2[0].Note, checkT2.Note)
	s.Require().Equal(transactions2[0].Date, checkT2.Date)
}

func (s *TransactionSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, delete successfully":       delete_WithOneData_DeleteSuccessfully,
		"when with multiple data, delete successfully":  delete_WithMultipleData_DeleteSuccessfully,
		"when with multiple users, delete successfully": delete_WithMultipleUsers_DeleteSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func delete_WithOneData_DeleteSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	err = s.transactionModel.Delete(mockCtx, transactions[0].ID)
	s.Require().NoError(err, desc)

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	s.Require().NoError(err, desc)
	s.Require().Equal(0, count, desc)

	// check if data exists
	var checkT transaction.Transaction
	stmt := "SELECT id FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.ID)
	s.Require().Equal(sql.ErrNoRows, err, desc)
}

func delete_WithMultipleData_DeleteSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(3)
	s.Require().NoError(err, desc)

	err = s.transactionModel.Delete(mockCtx, transactions[0].ID)
	s.Require().NoError(err, desc)

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	s.Require().NoError(err, desc)
	s.Require().Equal(2, count, desc)

	// check if data exists
	var checkT transaction.Transaction
	stmt := "SELECT id FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.ID)
	s.Require().Equal(sql.ErrNoRows, err, desc)
}

func delete_WithMultipleUsers_DeleteSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(3)
	s.Require().NoError(err, desc)

	// prepare more users
	transactions2, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	err = s.transactionModel.Delete(mockCtx, transactions[0].ID)
	s.Require().NoError(err, desc)

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	s.Require().NoError(err, desc)
	s.Require().Equal(3, count, desc)

	// check if data exists
	var checkT transaction.Transaction
	stmt := "SELECT id FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.ID)
	s.Require().Equal(sql.ErrNoRows, err, desc)

	// check if other user's data still exists
	var countUser2 int
	stmt = "SELECT COUNT(*) FROM transactions WHERE user_id = ?"
	err = s.db.QueryRow(stmt, transactions2[0].UserID).Scan(&countUser2)
	s.Require().NoError(err, desc)
	s.Require().Equal(1, countUser2, desc)
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
	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue()}
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
	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue()}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(1, ow1)
	s.Require().NoError(err, desc)

	ow2 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue()}
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
	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue()}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue()}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue()}
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
	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
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

	startDate := mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly)
	query := domain.GetAccInfoQuery{
		StartDate: &startDate,
	}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryEndDate_ReturnDataBeforeEndDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
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

	endDate := mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly)
	query := domain.GetAccInfoQuery{
		EndDate: &endDate,
	}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := transaction.Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
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

	startDate := mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly)
	endDate := mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly)
	query := domain.GetAccInfoQuery{
		StartDate: &startDate,
		EndDate:   &endDate,
	}
	accInfo, err := s.transactionModel.GetAccInfo(mockCtx, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func (s *TransactionSuite) TestGetByIDAndUserID() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when only one data, return successfully":       getByIDAndUserID_OnlyOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":  getByIDAndUserID_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getByIDAndUserID_WithMultipleUsers_ReturnSuccessfully,
		"when id not found, return error":               getByIDAndUserID_IDNotFound_ReturnError,
		"when user id not found, return error":          getByIDAndUserID_UserIDNotFound_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByIDAndUserID_OnlyOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := domain.Transaction{
		ID:     transactions[0].ID,
		UserID: transactions[0].UserID,
		Type:   domain.CvtToTransactionType(transactions[0].Type),
		Price:  transactions[0].Price,
		Note:   transactions[0].Note,
		Date:   transactions[0].Date,
	}

	trans, err := s.transactionModel.GetByIDAndUserID(mockCtx, transactions[0].ID, transactions[0].UserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getByIDAndUserID_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(3)
	s.Require().NoError(err, desc)

	expResult := domain.Transaction{
		ID:     transactions[0].ID,
		UserID: transactions[0].UserID,
		Type:   domain.CvtToTransactionType(transactions[0].Type),
		Price:  transactions[0].Price,
		Note:   transactions[0].Note,
		Date:   transactions[0].Date,
	}

	trans, err := s.transactionModel.GetByIDAndUserID(mockCtx, transactions[0].ID, transactions[0].UserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, err)
}

func getByIDAndUserID_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(3)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	expResult := domain.Transaction{
		ID:     transactions[0].ID,
		UserID: transactions[0].UserID,
		Type:   domain.CvtToTransactionType(transactions[0].Type),
		Price:  transactions[0].Price,
		Note:   transactions[0].Note,
		Date:   transactions[0].Date,
	}

	trans, err := s.transactionModel.GetByIDAndUserID(mockCtx, transactions[0].ID, transactions[0].UserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, err)
}

func getByIDAndUserID_IDNotFound_ReturnError(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	_, err = s.transactionModel.GetByIDAndUserID(mockCtx, transactions[0].ID+1, transactions[0].UserID)
	s.Require().Error(err, desc)
}

func getByIDAndUserID_UserIDNotFound_ReturnError(s *TransactionSuite, desc string) {
	transactions, _, _, _, _, err := s.f.InsertTransactionsWithOneUser(1)
	s.Require().NoError(err, desc)

	_, err = s.transactionModel.GetByIDAndUserID(mockCtx, transactions[0].ID, transactions[0].UserID+1)
	s.Require().Error(err, desc)
}

func (s *TransactionSuite) TestGetBarChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, return successfully":       getBarChartData_WithOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":  getBarChartData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getBarChartData_WithMultipleUsers_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getBarChartData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{
		Price: 999,
		Type:  domain.TransactionTypeExpense.ToModelValue(),
		Date:  mockTimeNow.AddDate(0, 0, -4),
	}
	transactions, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(1, ow1)
	s.Require().NoError(err, desc)

	expResult := domain.ChartDataByWeekday{
		ow1.Date.Format(weekDayFormat): 999,
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		StartDate: transactions[0].Date.Format(time.DateOnly),
		EndDate:   transactions[0].Date.Format(time.DateOnly),
	}
	chartData, err := s.transactionModel.GetBarChartData(mockCtx, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getBarChartData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	date, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)

	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 0)}
	ow2 := transaction.Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 0)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 1)}
	ow4 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 1)}
	ow5 := transaction.Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 2)}
	ow6 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 2)}
	ow7 := transaction.Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 3)}
	ow8 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 3)}
	ow9 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 4)}
	ow10 := transaction.Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 4)}
	transactions, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(10, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	expResult := domain.ChartDataByWeekday{
		transactions[0].Date.Format(weekDayFormat): 1000,
		transactions[2].Date.Format(weekDayFormat): 1000,
		transactions[4].Date.Format(weekDayFormat): 999,
		transactions[6].Date.Format(weekDayFormat): 1001,
		transactions[8].Date.Format(weekDayFormat): 2000,
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		StartDate: transactions[0].Date.Format(time.DateOnly),
		EndDate:   transactions[8].Date.Format(time.DateOnly),
	}
	chartData, err := s.transactionModel.GetBarChartData(mockCtx, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getBarChartData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	date, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)

	ow1 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 0)}
	ow2 := transaction.Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 0)}
	ow3 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 1)}
	ow4 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 1)}
	ow5 := transaction.Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 2)}
	ow6 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 2)}
	ow7 := transaction.Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 3)}
	ow8 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 3)}
	ow9 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 4)}
	ow10 := transaction.Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 4)}
	transactions, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(10, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	// prepare more users
	ow11 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 1)}
	ow12 := transaction.Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 2)}
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(2, ow11, ow12)
	s.Require().NoError(err, desc)

	ow13 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 3)}
	ow14 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 4)}
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(2, ow13, ow14)
	s.Require().NoError(err, desc)

	expResult := domain.ChartDataByWeekday{
		transactions[0].Date.Format(weekDayFormat): 1000,
		transactions[2].Date.Format(weekDayFormat): 1000,
		transactions[4].Date.Format(weekDayFormat): 999,
		transactions[6].Date.Format(weekDayFormat): 1001,
		transactions[8].Date.Format(weekDayFormat): 2000,
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		StartDate: transactions[0].Date.Format(time.DateOnly),
		EndDate:   transactions[8].Date.Format(time.DateOnly),
	}
	chartData, err := s.transactionModel.GetBarChartData(mockCtx, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func (s *TransactionSuite) TestGetPieChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when pie chart with one data, return successfully":      getPieChartData_WithOneData_ReturnSuccessfully,
		"when pie chart with multiple data, return successfully": getPieChartData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully":          getPieChartData_WithMultipleUsers_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getPieChartData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	ow1 := transaction.Transaction{
		Price: 999,
		Type:  domain.TransactionTypeExpense.ToModelValue(),
		Date:  mockTimeNow.AddDate(0, 0, -4),
	}
	transactions, user, mainCategs, _, _, err := s.f.InsertTransactionsWithOneUser(1, ow1)
	s.Require().NoError(err, desc)

	expResult := domain.ChartData{
		Labels:   []string{mainCategs[0].Name},
		Datasets: []float64{999},
	}

	dataRange := domain.ChartDateRange{
		StartDate: transactions[0].Date.Format(time.DateOnly),
		EndDate:   transactions[0].Date.Format(time.DateOnly),
	}
	transactionType := domain.TransactionTypeExpense

	chartData, err := s.transactionModel.GetPieChartData(mockCtx, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getPieChartData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, _, err := s.f.InsertMainCategList(5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	date, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)

	ow1 := transaction.Transaction{Price: 999, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: date}
	ow2 := transaction.Transaction{Price: 1, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: date}
	ow3 := transaction.Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: date}
	ow4 := transaction.Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: date}
	ow5 := transaction.Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: date}
	ow6 := transaction.Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: date}

	// income typ data
	ow7 := transaction.Transaction{Price: 999, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: date}
	ow8 := transaction.Transaction{Price: 1, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: date}

	// data out of date range
	ow9 := transaction.Transaction{Price: 1000, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: date.AddDate(0, 0, 10)}
	ow10 := transaction.Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: date.AddDate(0, 0, 10)}

	_, _, err = s.f.InsertTransactionWithGivenUser(10, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	expResult := domain.ChartData{
		Labels:   []string{mainCategList[0].Name, mainCategList[1].Name, mainCategList[2].Name},
		Datasets: []float64{1000, 2000, 4000},
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		StartDate: date.Format(time.DateOnly),
		EndDate:   date.AddDate(0, 0, 1).Format(time.DateOnly),
	}
	chartData, err := s.transactionModel.GetPieChartData(mockCtx, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getPieChartData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, _, err := s.f.InsertMainCategList(5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	date, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)

	ow1 := transaction.Transaction{Price: 999, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: date}
	ow2 := transaction.Transaction{Price: 1, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: date}
	ow3 := transaction.Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: date}
	ow4 := transaction.Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: date}
	ow5 := transaction.Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: date}
	ow6 := transaction.Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: date}

	// income typ data
	ow7 := transaction.Transaction{Price: 999, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: date}
	ow8 := transaction.Transaction{Price: 1, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: date}

	// data out of date range
	ow9 := transaction.Transaction{Price: 1000, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: date.AddDate(0, 0, 10)}
	ow10 := transaction.Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: date.AddDate(0, 0, 10)}

	_, _, err = s.f.InsertTransactionWithGivenUser(10, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	// prepare more user
	ow11 := transaction.Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 1)}
	ow12 := transaction.Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 2)}
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(2, ow11, ow12)
	s.Require().NoError(err, desc)

	ow13 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 3)}
	ow14 := transaction.Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 4)}
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(2, ow13, ow14)
	s.Require().NoError(err, desc)

	expResult := domain.ChartData{
		Labels:   []string{mainCategList[0].Name, mainCategList[1].Name, mainCategList[2].Name},
		Datasets: []float64{1000, 2000, 4000},
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		StartDate: date.Format(time.DateOnly),
		EndDate:   date.AddDate(0, 0, 1).Format(time.DateOnly),
	}
	chartData, err := s.transactionModel.GetPieChartData(mockCtx, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func (s *TransactionSuite) TestGetMonthlyData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, return successfully":       getMonthlyData_WithOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":  getMonthlyData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getMonthlyData_WithMultipleUsers_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getMonthlyData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err, desc)

	ow1 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 3)}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(1, ow1)
	s.Require().NoError(err, desc)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expResult := domain.MonthDayToTransactionType{
		"4": domain.TransactionTypeExpense, // startDate.AddDate(0, 0, 3)
	}

	monthlyData, err := s.transactionModel.GetMonthlyData(mockCtx, dateRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyData, desc)
}

func getMonthlyData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err, desc)

	ow1 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow2 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow3 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow4 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow5 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow6 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow7 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow8 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow9 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	ow10 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(10, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expResult := domain.MonthDayToTransactionType{
		"2": domain.TransactionTypeExpense, // startDate.AddDate(0, 0, 1)
		"5": domain.TransactionTypeIncome,  // startDate.AddDate(0, 0, 4)
		"6": domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 5)
		"7": domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 6)
	}

	monthlyData, err := s.transactionModel.GetMonthlyData(mockCtx, dateRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyData, desc)
}

func getMonthlyData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err, desc)

	ow1 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow2 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow3 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow4 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow5 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow6 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow7 := transaction.Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow8 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow9 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	ow10 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	_, user, _, _, _, err := s.f.InsertTransactionsWithOneUser(10, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	// prepare more users
	ow11 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1, ow11)
	s.Require().NoError(err, desc)
	ow12 := transaction.Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 2)}
	_, _, _, _, _, err = s.f.InsertTransactionsWithOneUser(1, ow12)
	s.Require().NoError(err, desc)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expResult := domain.MonthDayToTransactionType{
		"2": domain.TransactionTypeExpense, // startDate.AddDate(0, 0, 1)
		"5": domain.TransactionTypeIncome,  // startDate.AddDate(0, 0, 4)
		"6": domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 5)
		"7": domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 6)
	}

	monthlyData, err := s.transactionModel.GetMonthlyData(mockCtx, dateRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyData, desc)
}
