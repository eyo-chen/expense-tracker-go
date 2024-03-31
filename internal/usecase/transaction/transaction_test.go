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
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when time range type is one week day, return week day data": getBarChartData_WithTimeRangeTypeOneWeekDay_ReturnWeekDayData,
		"when time range type is one week, return date data":         getBarChartData_WithTimeRangeTypeOneWeek_ReturnDateData,
		"when time range type is two weeks, return date data":        getBarChartData_WithTimeRangeTypeTwoWeeks_ReturnDateData,
		"when time range type is one month, return date data":        getBarChartData_WithTimeRangeTypeOneMonth_ReturnDateData,
		"when get chart data fail, return error":                     getBarChartData_GetChartDataFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getBarChartData_WithTimeRangeTypeOneWeekDay_ReturnWeekDayData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": 200,
		"2024-03-21": 500,
		"2024-03-22": 600,
	}

	s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeekDay, domain.TransactionTypeExpense, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithTimeRangeTypeOneWeek_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": 200,
		"2024-03-21": 500,
		"2024-03-22": 600,
	}

	s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeek, domain.TransactionTypeExpense, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithTimeRangeTypeTwoWeeks_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-30")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": 200,
		"2024-03-21": 500,
		"2024-03-22": 600,
		"2024-03-23": 700,
		"2024-03-28": 800,
		"2024-03-29": 900,
	}

	s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23", "03/24", "03/25", "03/26", "03/27", "03/28", "03/29", "03/30"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 700, 0, 0, 0, 0, 800, 900, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeTwoWeeks, domain.TransactionTypeExpense, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithTimeRangeTypeOneMonth_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-01": 100,
		"2024-03-02": 200,
		"2024-03-05": 500,
		"2024-03-06": 600,
		"2024-03-14": 700,
		"2024-03-19": 800,
		"2024-03-20": 900,
	}

	s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/01", "03/02", "03/03", "03/04", "03/05", "03/06", "03/07", "03/08", "03/09", "03/10", "03/11", "03/12", "03/13", "03/14", "03/15", "03/16", "03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23", "03/24", "03/25", "03/26", "03/27", "03/28", "03/29", "03/30", "03/31"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 0, 0, 0, 0, 0, 0, 0, 700, 0, 0, 0, 0, 800, 900, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneMonth, domain.TransactionTypeExpense, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_GetChartDataFail_ReturnError(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	s.mockTransaction.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
		Return(domain.DateToChartData{}, errors.New("error")).Once()

	// prepare expected result
	expResult := domain.ChartData{}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeekDay, domain.TransactionTypeExpense, domain.User{ID: 1})
	s.Require().Equal(errors.New("error"), err, desc)
	s.Require().Equal(expResult, result, desc)
}

func (s *TransactionSuite) TestGetPieChartData() {
	tests := []struct {
		desc            string
		setupFun        func() domain.ChartDateRange
		transactionType domain.TransactionType
		user            domain.User
		expResult       domain.ChartData
		expErr          error
	}{
		{
			desc: "when no error, return chart data",
			setupFun: func() domain.ChartDateRange {
				start, err := time.Parse(time.DateOnly, "2024-03-17")
				s.Require().NoError(err)
				end, err := time.Parse(time.DateOnly, "2024-03-23")
				s.Require().NoError(err)

				chartDataRange := domain.ChartDateRange{
					Start: start,
					End:   end,
				}

				chartData := domain.ChartData{
					Labels:   []string{"label1", "label2"},
					Datasets: []float64{100, 200},
				}

				s.mockTransaction.On("GetPieChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
					Return(chartData, nil).Once()

				return chartDataRange
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
			setupFun: func() domain.ChartDateRange {
				start, err := time.Parse(time.DateOnly, "2024-03-17")
				s.Require().NoError(err)
				end, err := time.Parse(time.DateOnly, "2024-03-23")
				s.Require().NoError(err)

				chartDataRange := domain.ChartDateRange{
					Start: start,
					End:   end,
				}

				s.mockTransaction.On("GetPieChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
					Return(domain.ChartData{}, errors.New("error")).Once()

				return chartDataRange
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
			dateRange := t.setupFun()

			result, err := s.transactionUC.GetPieChartData(mockCtx, dateRange, t.transactionType, t.user)
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
