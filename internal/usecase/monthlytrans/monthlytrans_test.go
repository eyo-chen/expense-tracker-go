package monthlytrans

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx = context.Background()
)

type MonthlyTransSuite struct {
	suite.Suite
	uc                   *UC
	mockMonthlyTransRepo *mocks.MonthlyTransRepo
	mockTransactionRepo  *mocks.TransactionRepo
}

func TestMonthlyTransSuite(t *testing.T) {
	suite.Run(t, new(MonthlyTransSuite))
}

func (s *MonthlyTransSuite) SetupTest() {
	s.mockMonthlyTransRepo = mocks.NewMonthlyTransRepo(s.T())
	s.mockTransactionRepo = mocks.NewTransactionRepo(s.T())
	s.uc = New(s.mockMonthlyTransRepo, s.mockTransactionRepo)
}

func (s *MonthlyTransSuite) TearDownTest() {
	s.mockMonthlyTransRepo.AssertExpectations(s.T())
	s.mockTransactionRepo.AssertExpectations(s.T())
}

func (s *MonthlyTransSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *MonthlyTransSuite, desc string){
		"when no error, create successfully": create_NoError_CreateSuccessfully,
		"when get failed, return error":      create_GetFailed_ReturnError,
		"when create failed, return error":   create_CreateFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_CreateSuccessfully(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	mockDate, err := time.Parse(time.DateOnly, "2024-10-04")
	s.Require().NoError(err, desc)
	mockTrans := []domain.MonthlyAggregatedData{
		{
			UserID:       1,
			TotalExpense: 100,
			TotalIncome:  200,
		},
	}

	// prepare mock service
	s.mockTransactionRepo.On("GetMonthlyAggregatedData", mockCtx, mockDate).Return(mockTrans, nil)
	s.mockMonthlyTransRepo.On("Create", mockCtx, mockDate, mockTrans).Return(nil)

	// action
	err = s.uc.Create(mockCtx, mockDate)
	s.Require().NoError(err, desc)
}

func create_GetFailed_ReturnError(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	mockDate, err := time.Parse(time.DateOnly, "2024-10-04")
	s.Require().NoError(err, desc)
	mockErr := errors.New("failed to get monthly aggregated data")

	// prepare mock service
	s.mockTransactionRepo.On("GetMonthlyAggregatedData", mockCtx, mockDate).Return(nil, mockErr)

	// action
	err = s.uc.Create(mockCtx, mockDate)
	s.Require().ErrorIs(err, mockErr, desc)
}

func create_CreateFailed_ReturnError(s *MonthlyTransSuite, desc string) {
	// prepare mock data
	mockDate, err := time.Parse(time.DateOnly, "2024-10-04")
	s.Require().NoError(err, desc)
	mockTrans := []domain.MonthlyAggregatedData{
		{
			UserID:       1,
			TotalExpense: 100,
			TotalIncome:  200,
		},
	}
	mockErr := errors.New("failed to create monthly transactions")

	// prepare mock service
	s.mockTransactionRepo.On("GetMonthlyAggregatedData", mockCtx, mockDate).Return(mockTrans, nil)
	s.mockMonthlyTransRepo.On("Create", mockCtx, mockDate, mockTrans).Return(mockErr)

	// action
	err = s.uc.Create(mockCtx, mockDate)
	s.Require().ErrorIs(err, mockErr, desc)
}
