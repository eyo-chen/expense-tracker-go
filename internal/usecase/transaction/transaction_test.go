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
				s.mockTransaction.On("GetChartData", mockCtx, domain.ChartTypeBar, domain.ChartDateRange{
					StartDate: "2021-01-01",
					EndDate:   "2021-01-31",
				}, int64(1)).
					Return(domain.ChartData{
						Labels:   []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
						Datasets: []float64{100, 200, 300, 400, 500, 600, 700},
					}, nil).Once()
			},
			chartType: domain.ChartTypeBar,
			chartDateRange: domain.ChartDateRange{
				StartDate: "2021-01-01",
				EndDate:   "2021-01-31",
			},
			user: domain.User{
				ID: 1,
			},
			expResult: domain.ChartData{
				Labels:   []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
				Datasets: []float64{100, 200, 300, 400, 500, 600, 700},
			},
			expErr: nil,
		},
		{
			desc: "when get chart data fail, return error",
			setupFun: func() {
				s.mockTransaction.On("GetChartData", mockCtx, domain.ChartTypeBar, domain.ChartDateRange{
					StartDate: "2021-01-01",
					EndDate:   "2021-01-31",
				}, int64(1)).
					Return(domain.ChartData{}, errors.New("error")).Once()
			},
			chartType: domain.ChartTypeBar,
			chartDateRange: domain.ChartDateRange{
				StartDate: "2021-01-01",
				EndDate:   "2021-01-31",
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
