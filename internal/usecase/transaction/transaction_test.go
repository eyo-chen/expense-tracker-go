package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/transaction"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx     = context.Background()
	mockLoc, _  = time.LoadLocation("")
	mockTimeNow = time.Unix(1629446406, 0).Truncate(24 * time.Hour).In(mockLoc)
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

func (s *TransactionSuite) SetupSuite() {
	logger.Register()
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

func (s *TransactionSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, update successfully":                                                      update_NoError_UpdateSuccessfully,
		"when get main category fail, return error":                                               update_GetMainCategFail_ReturnError,
		"when type of main category not match transaction type, return error":                     update_TypeNotMatch_ReturnError,
		"when get sub category fail, return error":                                                update_GetSubCategFail_ReturnError,
		"when main category of sub category not match main category of transaction, return error": update_MainCategNotMatch_ReturnError,
		"when get transaction fail, return error":                                                 update_GetTransFail_UpdateSuccessfully,
		"when update fail, return error":                                                          update_UpdateFail_UpdateSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoError_UpdateSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCateg.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCateg.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	s.mockTransaction.On("GetByIDAndUserID", mockCtx, trans.ID, user.ID).
		Return(domain.Transaction{}, nil).Once()

	s.mockTransaction.On("Update", mockCtx, trans).
		Return(nil).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().NoError(err, desc)
}

func update_GetMainCategFail_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCateg.On("GetByID", trans.MainCategID, user.ID).
		Return(nil, errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func update_TypeNotMatch_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeIncome} // set type to income
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCateg.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(domain.ErrTypeNotConsistent, err, desc)
}

func update_GetSubCategFail_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCateg.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCateg.On("GetByID", trans.SubCategID, user.ID).
		Return(nil, errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func update_MainCategNotMatch_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 2}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCateg.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCateg.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(domain.ErrMainCategNotConsistent, err, desc)
}

func update_GetTransFail_UpdateSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCateg.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCateg.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	s.mockTransaction.On("GetByIDAndUserID", mockCtx, trans.ID, user.ID).
		Return(domain.Transaction{}, errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func update_UpdateFail_UpdateSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCateg.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCateg.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	s.mockTransaction.On("GetByIDAndUserID", mockCtx, trans.ID, user.ID).
		Return(domain.Transaction{}, nil).Once()

	s.mockTransaction.On("Update", mockCtx, trans).
		Return(errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
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

func (s *TransactionSuite) TestGetBarChartData() {
	tests := []struct {
		desc            string
		setupFun        func()
		chartDateRange  domain.ChartDateRange
		transactionType domain.TransactionType
		user            domain.User
		expResult       domain.ChartData
		expErr          error
	}{
		{
			desc: "when no error, return chart data",
			setupFun: func() {
				chartDataRange := domain.ChartDateRange{
					StartDate: "2024-03-17",
					EndDate:   "2024-03-23",
				}

				DateToChartData := domain.DateToChartData{
					"2024-03-17": 100,
					"2024-03-18": 200,
					"2024-03-19": 300,
					"2024-03-20": 400,
					"2024-03-21": 500,
					"2024-03-22": 600,
					"2024-03-23": 700,
				}

				s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
					Return(DateToChartData, nil).Once()
			},
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			transactionType: domain.TransactionTypeExpense,
			user: domain.User{
				ID: 1,
			},
			expResult: domain.ChartData{
				Labels:   []string{"2024-03-17", "2024-03-18", "2024-03-19", "2024-03-20", "2024-03-21", "2024-03-22", "2024-03-23"},
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
				DateToChartData := domain.DateToChartData{
					"2024-03-17": 100,
					"2024-03-18": 200,
					"2024-03-19": 300,
				}

				s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
					Return(DateToChartData, nil).Once()
			},
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			transactionType: domain.TransactionTypeExpense,
			user: domain.User{
				ID: 1,
			},
			expResult: domain.ChartData{
				Labels:   []string{"2024-03-17", "2024-03-18", "2024-03-19", "2024-03-20", "2024-03-21", "2024-03-22", "2024-03-23"},
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

				s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
					Return(nil, errors.New("error")).Once()
			},
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			transactionType: domain.TransactionTypeExpense,
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

			result, err := s.transactionUC.GetBarChartData(mockCtx, t.chartDateRange, t.transactionType, t.user)
			s.Require().Equal(t.expResult, result)
			s.Require().Equal(t.expErr, err)

			s.TearDownTest()
		})
	}
}

func (s *TransactionSuite) TestGetPieChartData() {
	tests := []struct {
		desc            string
		setupFun        func()
		chartDateRange  domain.ChartDateRange
		transactionType domain.TransactionType
		user            domain.User
		expResult       domain.ChartData
		expErr          error
	}{
		{
			desc: "when no error, return chart data",
			setupFun: func() {
				chartDataRange := domain.ChartDateRange{
					StartDate: "2024-03-17",
					EndDate:   "2024-03-23",
				}

				chartData := domain.ChartData{
					Labels:   []string{"label1", "label2"},
					Datasets: []float64{100, 200},
				}

				s.mockTransaction.On("GetPieChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
					Return(chartData, nil).Once()
			},
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			transactionType: domain.TransactionTypeExpense,
			user: domain.User{
				ID: 1,
			},
			expResult: domain.ChartData{
				Labels:   []string{"label1", "label2"},
				Datasets: []float64{100, 200},
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

				s.mockTransaction.On("GetPieChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
					Return(domain.ChartData{}, errors.New("error")).Once()
			},
			chartDateRange: domain.ChartDateRange{
				StartDate: "2024-03-17",
				EndDate:   "2024-03-23",
			},
			transactionType: domain.TransactionTypeExpense,
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

			result, err := s.transactionUC.GetPieChartData(mockCtx, t.chartDateRange, t.transactionType, t.user)
			s.Require().Equal(t.expResult, result)
			s.Require().Equal(t.expErr, err)

			s.TearDownTest()
		})
	}
}

func (s *TransactionSuite) TestGetMonthlyData() {
	tests := []struct {
		desc     string
		setupFun func() (domain.GetMonthlyDateRange, domain.MonthDayToTransactionType)
		user     domain.User
		expErr   error
	}{
		{
			desc: "when it's 31 day in a month, return monthly data",
			setupFun: func() (domain.GetMonthlyDateRange, domain.MonthDayToTransactionType) {
				startDate, err := time.Parse(time.DateOnly, "2024-03-01")
				s.Require().NoError(err)
				endDate, err := time.Parse(time.DateOnly, "2024-03-31")
				s.Require().NoError(err)

				dateRange := domain.GetMonthlyDateRange{
					StartDate: startDate,
					EndDate:   endDate,
				}

				monthlyData := domain.MonthDayToTransactionType{
					5:  domain.TransactionTypeExpense,
					10: domain.TransactionTypeIncome,
					20: domain.TransactionTypeBoth,
				}

				s.mockTransaction.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
					Return(monthlyData, nil).Once()

				return dateRange, monthlyData
			},
			user: domain.User{
				ID: 1,
			},
		},
		{
			desc: "when it's 30 day in a month, return monthly data",
			setupFun: func() (domain.GetMonthlyDateRange, domain.MonthDayToTransactionType) {
				startDate, err := time.Parse(time.DateOnly, "2024-04-01")
				s.Require().NoError(err)
				endDate, err := time.Parse(time.DateOnly, "2024-04-30")
				s.Require().NoError(err)

				dateRange := domain.GetMonthlyDateRange{
					StartDate: startDate,
					EndDate:   endDate,
				}

				monthlyData := domain.MonthDayToTransactionType{
					3:  domain.TransactionTypeExpense,
					8:  domain.TransactionTypeIncome,
					10: domain.TransactionTypeBoth,
				}

				s.mockTransaction.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
					Return(monthlyData, nil).Once()

				return dateRange, monthlyData
			},
			user: domain.User{
				ID: 1,
			},
		},
		{
			desc: "when it's 29 day in a month, return monthly data",
			setupFun: func() (domain.GetMonthlyDateRange, domain.MonthDayToTransactionType) {
				startDate, err := time.Parse(time.DateOnly, "2024-02-01")
				s.Require().NoError(err)
				endDate, err := time.Parse(time.DateOnly, "2024-02-29")
				s.Require().NoError(err)

				dateRange := domain.GetMonthlyDateRange{
					StartDate: startDate,
					EndDate:   endDate,
				}

				monthlyData := domain.MonthDayToTransactionType{
					3:  domain.TransactionTypeExpense,
					8:  domain.TransactionTypeIncome,
					10: domain.TransactionTypeBoth,
				}

				s.mockTransaction.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
					Return(monthlyData, nil).Once()

				return dateRange, monthlyData
			},
			user: domain.User{
				ID: 1,
			},
		},
		{
			desc: "when get monthly data fail, return error",
			setupFun: func() (domain.GetMonthlyDateRange, domain.MonthDayToTransactionType) {
				startDate, err := time.Parse(time.DateOnly, "2024-05-01")
				s.Require().NoError(err)
				endDate, err := time.Parse(time.DateOnly, "2024-05-31")
				s.Require().NoError(err)

				dateRange := domain.GetMonthlyDateRange{
					StartDate: startDate,
					EndDate:   endDate,
				}

				s.mockTransaction.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
					Return(nil, errors.New("error")).Once()

				return dateRange, nil

			},
			user: domain.User{
				ID: 1,
			},
			expErr: errors.New("error"),
		},
	}

	for _, t := range tests {
		s.Run(t.desc, func() {
			s.SetupTest()
			dateRange, monthlyData := t.setupFun()

			result, err := s.transactionUC.GetMonthlyData(mockCtx, dateRange, t.user)
			expResult := transaction.GetMonthlyData_GenExpResult(monthlyData, dateRange.EndDate.Day(), err)
			s.Require().Equal(t.expErr, err, t.desc)
			s.Require().Equal(expResult, result, t.desc)

			s.TearDownTest()
		})
	}
}
