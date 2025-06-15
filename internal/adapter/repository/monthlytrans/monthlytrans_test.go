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
	createdTrans := getMonthlyTrans(s, date)
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
	createdTrans := getMonthlyTrans(s, date)
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
	createdTrans := getMonthlyTrans(s, date)
	s.Require().Empty(createdTrans, desc)
}

func getMonthlyTrans(s *MonthlyTransSuite, date time.Time) []domain.MonthlyAggregatedData {
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
