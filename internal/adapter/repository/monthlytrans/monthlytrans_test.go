package monthlytrans

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/dockerutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX = context.Background()
)

type MonthlyTransSuite struct {
	suite.Suite
	dk      *dockerutil.Container
	db      *sql.DB
	migrate *migrate.Migrate
	repo    *Repo
	f       *factory
}

func TestMonthlyTransSuite(t *testing.T) {
	suite.Run(t, new(MonthlyTransSuite))
}

func (s *MonthlyTransSuite) SetupSuite() {
	s.dk = dockerutil.RunDocker(dockerutil.ImageMySQL)
	db, migrate := testutil.ConnToDB(s.dk.Port)
	logger.Register()

	s.db = db
	s.migrate = migrate
	s.repo = New(s.db)
	s.f = newFactory(s.db)
}

func (s *MonthlyTransSuite) TearDownSuite() {
	if err := s.db.Close(); err != nil {
		logger.Error("Unable to close mysql database", "error", err)
	}
	s.migrate.Close()
	s.dk.PurgeDocker()
}

func (s *MonthlyTransSuite) SetupTest() {
	s.repo = New(s.db)
	s.f = newFactory(s.db)
}

func (s *MonthlyTransSuite) TearDownTest() {
	tx, err := s.db.Begin()
	if err != nil {
		s.Require().NoError(err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.Require().NoError(err)
		}
	}()

	if _, err := tx.Exec("DELETE FROM monthly_transactions"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM users"); err != nil {
		s.Require().NoError(err)
	}

	s.Require().NoError(tx.Commit())
	s.f.Reset()
}

func (s *MonthlyTransSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *MonthlyTransSuite, desc string){
		"when one data, insert successfully":      create_OneData_InsertSuccessfully,
		"when multiple data, insert successfully": create_MultipleData_InsertSuccessfully,
		"when same user and date, return error":   create_SameUserAndDate_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_OneData_InsertSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	users, err := s.f.InsertUsers(mockCTX, 1)
	s.Require().NoError(err, desc)

	// prepare input data
	date, err := time.Parse(time.DateOnly, "2024-10-04")
	s.Require().NoError(err, desc)
	trans := []domain.MonthlyAggregatedData{
		{
			UserID:       users[0].ID,
			TotalExpense: 100,
			TotalIncome:  200,
		},
	}

	// action
	err = s.repo.Create(mockCTX, date, trans)
	s.Require().NoError(err, desc)

	// assertion
	createdTrans := getMonthlyAggTrans(s, date)
	s.Require().Equal(trans, createdTrans, desc)
}

func create_MultipleData_InsertSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	users, err := s.f.InsertUsers(mockCTX, 3)
	s.Require().NoError(err, desc)

	// prepare input data
	date, err := time.Parse(time.DateOnly, "2024-10-04")
	s.Require().NoError(err, desc)
	trans := []domain.MonthlyAggregatedData{
		{
			UserID:       users[0].ID,
			TotalExpense: 100,
			TotalIncome:  200,
		},
		{
			UserID:       users[1].ID,
			TotalExpense: 300,
			TotalIncome:  400,
		},
		{
			UserID:       users[2].ID,
			TotalExpense: 500,
			TotalIncome:  600,
		},
	}

	// action
	err = s.repo.Create(mockCTX, date, trans)
	s.Require().NoError(err, desc)

	// assertion
	createdTrans := getMonthlyAggTrans(s, date)
	s.Require().Equal(trans, createdTrans, desc)
}

func create_SameUserAndDate_ReturnError(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	users, err := s.f.InsertUsers(mockCTX, 1)
	s.Require().NoError(err, desc)

	// prepare input data
	date, err := time.Parse(time.DateOnly, "2024-10-04")
	s.Require().NoError(err, desc)

	trans := []domain.MonthlyAggregatedData{
		{
			UserID:       users[0].ID,
			TotalExpense: 100,
			TotalIncome:  200,
		},
		{
			UserID:       users[0].ID,
			TotalExpense: 100,
			TotalIncome:  200,
		},
	}

	// action
	err = s.repo.Create(mockCTX, date, trans)
	s.Require().ErrorIs(err, domain.ErrUniqueUserDate, desc)

	// assertion
	createdTrans := getMonthlyAggTrans(s, date)
	s.Require().Empty(createdTrans, desc)
}

func getMonthlyAggTrans(s *MonthlyTransSuite, date time.Time) []domain.MonthlyAggregatedData {
	stmt := `
		SELECT user_id, total_expense, total_income
		FROM monthly_transactions
		WHERE month_date = ?
	`
	rows, err := s.db.QueryContext(mockCTX, stmt, date)
	s.Require().NoError(err)
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Unable to close rows", "package", packageName, "err", err)
		}
	}()

	createdTrans := []domain.MonthlyAggregatedData{}
	for rows.Next() {
		var trans domain.MonthlyAggregatedData
		err := rows.Scan(&trans.UserID, &trans.TotalExpense, &trans.TotalIncome)
		s.Require().NoError(err)

		createdTrans = append(createdTrans, trans)
	}

	return createdTrans
}

func (s *MonthlyTransSuite) TestGetByUserIDAndMonthDate() {
	for scenario, fn := range map[string]func(s *MonthlyTransSuite, desc string){
		"when one data, return successfully":  getByUserIDAndMonthDate_OneData_ReturnSuccessfully,
		"when many data, return successfully": getByUserIDAndMonthDate_ManyData_ReturnSuccessfully,
		"when many user, return successfully": getByUserIDAndMonthDate_ManyUser_ReturnSuccessfully,
		"when data not found, return error":   getByUserIDAndMonthDate_DataNotFound_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByUserIDAndMonthDate_OneData_ReturnSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate := "2024-10-01"
	monthDate, err := time.Parse(time.DateOnly, startDate)
	s.Require().NoError(err, desc)

	ows := []MonthlyTrans{{MonthDate: monthDate, TotalExpense: 100, TotalIncome: 200}}
	user, monthlyTrans, err := s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 1, ows)
	s.Require().NoError(err, desc)

	// prepare expected result
	expRes := domain.AccInfo{
		TotalExpense: monthlyTrans[0].TotalExpense,
		TotalIncome:  monthlyTrans[0].TotalIncome,
		TotalBalance: monthlyTrans[0].TotalIncome - monthlyTrans[0].TotalExpense,
	}

	// action
	accInfo, err := s.repo.GetByUserIDAndMonthDate(mockCTX, user.ID, monthDate)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, accInfo, desc)
}

func getByUserIDAndMonthDate_ManyData_ReturnSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate := "2024-10-01"
	monthDate, err := time.Parse(time.DateOnly, startDate)
	s.Require().NoError(err, desc)
	startDate1 := "2024-11-01"
	monthDate1, err := time.Parse(time.DateOnly, startDate1)
	s.Require().NoError(err, desc)

	ows := []MonthlyTrans{
		{MonthDate: monthDate, TotalExpense: 100, TotalIncome: 200},
		{MonthDate: monthDate1, TotalExpense: 300, TotalIncome: 400},
	}
	user, monthlyTrans, err := s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 2, ows)
	s.Require().NoError(err, desc)

	// prepare expected result
	expRes := domain.AccInfo{
		TotalExpense: monthlyTrans[0].TotalExpense,
		TotalIncome:  monthlyTrans[0].TotalIncome,
		TotalBalance: monthlyTrans[0].TotalIncome - monthlyTrans[0].TotalExpense,
	}

	// action
	accInfo, err := s.repo.GetByUserIDAndMonthDate(mockCTX, user.ID, monthDate)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, accInfo, desc)
}

func getByUserIDAndMonthDate_ManyUser_ReturnSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate := "2024-10-01"
	monthDate, err := time.Parse(time.DateOnly, startDate)
	s.Require().NoError(err, desc)
	startDate1 := "2024-11-01"
	monthDate1, err := time.Parse(time.DateOnly, startDate1)
	s.Require().NoError(err, desc)

	ows := []MonthlyTrans{
		{MonthDate: monthDate, TotalExpense: 100, TotalIncome: 200},
		{MonthDate: monthDate1, TotalExpense: 300, TotalIncome: 400},
	}
	user, monthlyTrans, err := s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 2, ows)
	s.Require().NoError(err, desc)

	// prepare another user
	ows1 := []MonthlyTrans{
		{MonthDate: monthDate, TotalExpense: 500, TotalIncome: 600},
		{MonthDate: monthDate1, TotalExpense: 700, TotalIncome: 800},
	}
	_, _, err = s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 2, ows1)
	s.Require().NoError(err, desc)

	// prepare expected result
	expRes := domain.AccInfo{
		TotalExpense: monthlyTrans[0].TotalExpense,
		TotalIncome:  monthlyTrans[0].TotalIncome,
		TotalBalance: monthlyTrans[0].TotalIncome - monthlyTrans[0].TotalExpense,
	}

	// action
	accInfo, err := s.repo.GetByUserIDAndMonthDate(mockCTX, user.ID, monthDate)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, accInfo, desc)
}

func getByUserIDAndMonthDate_DataNotFound_ReturnError(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate := "2024-10-01"
	monthDate, err := time.Parse(time.DateOnly, startDate)
	s.Require().NoError(err, desc)
	startDate1 := "2024-11-01"
	monthDate1, err := time.Parse(time.DateOnly, startDate1)
	s.Require().NoError(err, desc)
	startDate2 := "2024-12-01"
	monthDate2, err := time.Parse(time.DateOnly, startDate2)
	s.Require().NoError(err, desc)

	ows := []MonthlyTrans{
		{MonthDate: monthDate, TotalExpense: 100, TotalIncome: 200},
		{MonthDate: monthDate1, TotalExpense: 300, TotalIncome: 400},
	}
	user, _, err := s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 2, ows)
	s.Require().NoError(err, desc)

	// prepare another user
	ows1 := []MonthlyTrans{
		{MonthDate: monthDate, TotalExpense: 500, TotalIncome: 600},
		{MonthDate: monthDate1, TotalExpense: 700, TotalIncome: 800},
	}
	_, _, err = s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 2, ows1)
	s.Require().NoError(err, desc)

	// action
	accInfo, err := s.repo.GetByUserIDAndMonthDate(mockCTX, user.ID, monthDate2)

	// assertion
	s.Require().ErrorIs(err, domain.ErrDataNotFound, desc)
	s.Require().Empty(accInfo, desc)
}

func (s *MonthlyTransSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *MonthlyTransSuite, desc string){
		"when one data, update successfully":      update_OneData_UpdateSuccessfully,
		"when multiple date, update successfully": update_MultipleDate_UpdateSuccessfully,
		"when multiple user, update successfully": update_MultipleUser_UpdateSuccessfully,
		"when invalid trans type, return error":   update_InvalidTransType_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_OneData_UpdateSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate := "2024-12-01"
	monthDate, err := time.Parse(time.DateOnly, startDate)
	s.Require().NoError(err, desc)

	ows := []MonthlyTrans{{MonthDate: monthDate, TotalExpense: 100, TotalIncome: 200}}
	user, _, err := s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 1, ows)
	s.Require().NoError(err, desc)

	// prepare expected result
	expRes := MonthlyTrans{
		UserID:       user.ID,
		MonthDate:    monthDate,
		TotalExpense: 100,
		TotalIncome:  300,
	}

	// action
	err = s.repo.Update(mockCTX, user.ID, monthDate, domain.TransactionTypeIncome, 100)

	// assertion
	s.Require().NoError(err, desc)
	updatedTrans := getMonthlyTrans(s, user.ID, monthDate)
	s.Require().Equal(updatedTrans, expRes, desc)
}

func update_MultipleDate_UpdateSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate1 := "2024-10-01"
	monthDate1, err := time.Parse(time.DateOnly, startDate1)
	s.Require().NoError(err, desc)
	startDate2 := "2024-11-01"
	monthDate2, err := time.Parse(time.DateOnly, startDate2)
	s.Require().NoError(err, desc)
	startDate3 := "2024-12-01"
	monthDate3, err := time.Parse(time.DateOnly, startDate3)
	s.Require().NoError(err, desc)

	ows := []MonthlyTrans{
		{MonthDate: monthDate1, TotalExpense: 100, TotalIncome: 200},
		{MonthDate: monthDate2, TotalExpense: 300, TotalIncome: 400},
		{MonthDate: monthDate3, TotalExpense: 500, TotalIncome: 600},
	}
	user, _, err := s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 3, ows)
	s.Require().NoError(err, desc)

	// prepare expected result
	expRes := MonthlyTrans{
		UserID:       user.ID,
		MonthDate:    monthDate1,
		TotalExpense: 200,
		TotalIncome:  200,
	}

	// action
	err = s.repo.Update(mockCTX, user.ID, monthDate1, domain.TransactionTypeExpense, 100)

	// assertion
	s.Require().NoError(err, desc)
	updatedTrans := getMonthlyTrans(s, user.ID, monthDate1)
	s.Require().Equal(updatedTrans, expRes, desc)
}

func update_MultipleUser_UpdateSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate1 := "2024-10-01"
	monthDate1, err := time.Parse(time.DateOnly, startDate1)
	s.Require().NoError(err, desc)
	startDate2 := "2024-11-01"
	monthDate2, err := time.Parse(time.DateOnly, startDate2)
	s.Require().NoError(err, desc)
	startDate3 := "2024-12-01"
	monthDate3, err := time.Parse(time.DateOnly, startDate3)
	s.Require().NoError(err, desc)

	ows1 := []MonthlyTrans{
		{MonthDate: monthDate1, TotalExpense: 100, TotalIncome: 200},
		{MonthDate: monthDate2, TotalExpense: 300, TotalIncome: 400},
		{MonthDate: monthDate3, TotalExpense: 500, TotalIncome: 600},
	}
	_, _, err = s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 3, ows1)
	s.Require().NoError(err, desc)
	ows2 := []MonthlyTrans{
		{MonthDate: monthDate1, TotalExpense: 800, TotalIncome: 900},
		{MonthDate: monthDate2, TotalExpense: 1000, TotalIncome: 1100},
		{MonthDate: monthDate3, TotalExpense: 1200, TotalIncome: 1300},
	}
	user, _, err := s.f.InsertManyMonthlyTransWithOneUser(mockCTX, 3, ows2)
	s.Require().NoError(err, desc)

	// prepare expected result
	expRes := MonthlyTrans{
		UserID:       user.ID,
		MonthDate:    monthDate2,
		TotalExpense: 1500,
		TotalIncome:  1100,
	}

	// action
	err = s.repo.Update(mockCTX, user.ID, monthDate2, domain.TransactionTypeExpense, 500)

	// assertion
	s.Require().NoError(err, desc)
	updatedTrans := getMonthlyTrans(s, user.ID, monthDate2)
	s.Require().Equal(updatedTrans, expRes, desc)
}

func update_InvalidTransType_ReturnError(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	startDate := "2024-10-01"
	monthDate, err := time.Parse(time.DateOnly, startDate)
	s.Require().NoError(err, desc)

	// action
	err = s.repo.Update(mockCTX, 1, monthDate, domain.TransactionTypeBoth, 100)

	// assertion
	s.Require().ErrorIs(err, domain.ErrInvalidTransType, desc)
}

func getMonthlyTrans(s *MonthlyTransSuite, userID int64, date time.Time) MonthlyTrans {
	stmt := `
		SELECT user_id, month_date, total_expense, total_income
		FROM monthly_transactions
		WHERE user_id = ? AND month_date = ?
	`

	row := s.db.QueryRowContext(mockCTX, stmt, userID, date)
	var trans MonthlyTrans
	err := row.Scan(&trans.UserID, &trans.MonthDate, &trans.TotalExpense, &trans.TotalIncome)
	s.Require().NoError(err)

	return trans
}
