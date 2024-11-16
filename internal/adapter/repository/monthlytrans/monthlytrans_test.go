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
	s.db.Close()
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
	defer rows.Close()

	createdTrans := []domain.MonthlyAggregatedData{}
	for rows.Next() {
		var trans domain.MonthlyAggregatedData
		err := rows.Scan(&trans.UserID, &trans.TotalExpense, &trans.TotalIncome)
		s.Require().NoError(err)

		createdTrans = append(createdTrans, trans)
	}

	return createdTrans
}
