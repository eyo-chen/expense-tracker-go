package transaction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/transaction"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx = context.Background()
)

type TransactionSuite struct {
	suite.Suite
	transactionUC   *transaction.TransactionUC
	mockTransaction *mocks.TransactionModel
	mockMainCateg   *mocks.MainCategModel
	mockSubCateg    *mocks.SubCategModel
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) SetupTest() {
	s.mockTransaction = mocks.NewTransactionModel(s.T())
	s.mockMainCateg = mocks.NewMainCategModel(s.T())
	s.mockSubCateg = mocks.NewSubCategModel(s.T())
	s.transactionUC = transaction.NewTransactionUC(s.mockTransaction, s.mockMainCateg, s.mockSubCateg)
}

func (s *TransactionSuite) TearDownTest() {
	s.mockTransaction.AssertExpectations(s.T())
	s.mockMainCateg.AssertExpectations(s.T())
	s.mockSubCateg.AssertExpectations(s.T())
}

func (s *TransactionSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, delete successfully":       delete_NoError_DeleteSuccessfully,
		"when check permession fail, return error": delete_CheckPermessionFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func delete_NoError_DeleteSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	s.mockTransaction.
		On("GetByIDAndUserID", mockCtx, int64(1), user.ID).
		Return(domain.Transaction{}, nil).Once()

	s.mockTransaction.On("Delete", mockCtx, int64(1)).
		Return(nil).Once()

	err := s.transactionUC.Delete(mockCtx, int64(1), user)
	s.Require().NoError(err, desc)
}

func delete_CheckPermessionFail_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	s.mockTransaction.
		On("GetByIDAndUserID", mockCtx, int64(1), user.ID).
		Return(domain.Transaction{}, errors.New("error")).Once()

	err := s.transactionUC.Delete(mockCtx, int64(1), user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func (s *TransactionSuite) TestGetChartData() {
	tests := []struct {
		desc           string
		setupFun       func()
		chartType      domain.ChartType
		chartDateRange domain.ChartDateRange
		user           domain.User
		expResult      domain.ChartData
		expErr         error
	}{
		{
			desc: "when no error, return chart data",
			setupFun: func() {
				chartDataRange := domain.ChartDateRange{
					StartDate: "2024-03-17",
					EndDate:   "2024-03-23",
				}

				chartDataByWeekday := domain.ChartDataByWeekday{
					"Sun": 100,
					"Mon": 200,
					"Tue": 300,
					"Wed": 400,
					"Thu": 500,
					"Fri": 600,
					"Sat": 700,
				}

				s.mockTransaction.On("GetChartData", mockCtx, domain.ChartTypeBar, chartDataRange, int64(1)).
					Return(chartDataByWeekday, nil).Once()
			},
			chartType: domain.ChartTypeBar,
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			user: domain.User{
				ID: 1,
			},
			expResult: domain.ChartData{
				Labels:   []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
				Datasets: []float64{100, 200, 300, 400, 500, 600, 700},
			},
			expErr: nil,
		},
		{
			desc: "when chart data by weekday is not fully filled, still return chart data",
			setupFun: func() {
				chartDataRange := domain.ChartDateRange{
					StartDate: "2024-03-17",
					EndDate:   "2024-03-23",
				}

				// only have data for Sun, Mon, Tue
				chartDataByWeekday := domain.ChartDataByWeekday{
					"Sun": 100,
					"Mon": 200,
					"Tue": 300,
				}

				s.mockTransaction.On("GetChartData", mockCtx, domain.ChartTypeBar, chartDataRange, int64(1)).
					Return(chartDataByWeekday, nil).Once()
			},
			chartType: domain.ChartTypeBar,
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			user: domain.User{
				ID: 1,
			},
			expResult: domain.ChartData{
				Labels:   []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
				Datasets: []float64{100, 200, 300, 0, 0, 0, 0},
			},
			expErr: nil,
		},
		{
			desc: "when get chart data fail, return error",
			setupFun: func() {
				chartDataRange := domain.ChartDateRange{
					StartDate: "2024-03-17",
					EndDate:   "2024-03-23",
				}

				s.mockTransaction.On("GetChartData", mockCtx, domain.ChartTypeBar, chartDataRange, int64(1)).
					Return(nil, errors.New("error")).Once()
			},
			chartType: domain.ChartTypeBar,
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			user: domain.User{
				ID: 1,
			},
			expResult: domain.ChartData{},
			expErr:    errors.New("error"),
		},
	}

	for _, t := range tests {
		s.Run(t.desc, func() {
			s.SetupTest()
			t.setupFun()

			result, err := s.transactionUC.GetChartData(mockCtx, t.chartType, t.chartDateRange, t.user)
			s.Require().Equal(t.expResult, result)
			s.Require().Equal(t.expErr, err)

			s.TearDownTest()
		})
	}
}
