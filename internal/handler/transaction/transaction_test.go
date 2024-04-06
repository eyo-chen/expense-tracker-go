package transaction_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/handler/transaction"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type TransactionSuite struct {
	suite.Suite
	transactionHlr    *transaction.TransactionHandler
	mockTransactionUC *mocks.TransactionUC
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) SetupSuite() {
	logger.Register()
}

func (s *TransactionSuite) SetupTest() {
	s.mockTransactionUC = mocks.NewTransactionUC(s.T())
	s.transactionHlr = transaction.NewTransactionHandler(s.mockTransactionUC)
}

func (s *TransactionSuite) TearDownTest() {
	s.mockTransactionUC.AssertExpectations(s.T())
}

func (s *TransactionSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, delete successfully":         delete_NoError_DeleteSuccessfully,
		"when id is incorrect, return bad request":   delete_IncorrectID_ReturnBadReq,
		"when id is less than 0, return bad request": delete_IDLessThanZero_ReturnBadReq,
		"when data not found, return bad request":    delete_DataNotFound_ReturnBadReq,
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

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.Delete))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = ctxutil.SetUser(req, &user)

	// mock service
	s.mockTransactionUC.On("Delete", req.Context(), int64(1), user).Return(nil)

	s.transactionHlr.Delete(res, req)

	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func delete_IncorrectID_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.Delete))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = mux.SetURLVars(req, map[string]string{"id": "a"})
	req = ctxutil.SetUser(req, &user)

	s.transactionHlr.Delete(res, req)

	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func delete_IDLessThanZero_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.Delete))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = mux.SetURLVars(req, map[string]string{"id": "-1"})
	req = ctxutil.SetUser(req, &user)

	s.transactionHlr.Delete(res, req)

	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func delete_DataNotFound_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.Delete))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = ctxutil.SetUser(req, &user)

	// mock service
	s.mockTransactionUC.On("Delete", req.Context(), int64(1), user).Return(domain.ErrTransactionDataNotFound)

	s.transactionHlr.Delete(res, req)

	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func (s *TransactionSuite) TestGetBarChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return data":                         getBarChartData_NoError_ReturnData,
		"when no main category ids, pass nil to service":     getBarChartData_NoMainCategoryIDs_PassNilToService,
		"when no start date, return bad request":             getBarChartData_NoStartDate_ReturnBadReq,
		"when no end date, return bad request":               getBarChartData_NoEndDate_ReturnBadReq,
		"when start date after end date, return bad request": getBarChartData_StartDateAfterEndDate_ReturnBadReq,
		"when no type, return bad request":                   getBarChartData_NoType_ReturnBadReq,
		"when no time range type, return bad request":        getBarChartData_NoTimeRangeType_ReturnBadReq,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getBarChartData_NoError_ReturnData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-08")
	s.Require().NoError(err, desc)
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&end_date=2024-03-08&type=expense&time_range=one_week_day&main_category_ids=1,2,3", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	dateRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	mainCategIDs := []int64{1, 2, 3}

	// mock service
	s.mockTransactionUC.On("GetBarChartData",
		req.Context(),
		dateRange,
		domain.TimeRangeTypeOneWeekDay,
		domain.TransactionTypeExpense,
		&mainCategIDs,
		user,
	).Return(domain.ChartData{
		Labels:   []string{"Mon", "Tue", "Wed"},
		Datasets: []float64{100, 200, 300},
	}, nil)

	// expected expected response
	expResp := map[string]interface{}{
		"chart_data": map[string]interface{}{
			"labels":   []interface{}{"Mon", "Tue", "Wed"},
			"datasets": []interface{}{100.0, 200.0, 300.0},
		},
	}

	// action
	s.transactionHlr.GetBarChartData(res, req)

	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func getBarChartData_NoMainCategoryIDs_PassNilToService(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-08")
	s.Require().NoError(err, desc)
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&end_date=2024-03-08&type=expense&time_range=one_week_day", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	dateRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	// mock service
	s.mockTransactionUC.On("GetBarChartData",
		req.Context(),
		dateRange,
		domain.TimeRangeTypeOneWeekDay,
		domain.TransactionTypeExpense,
		(*[]int64)(nil),
		user,
	).Return(domain.ChartData{
		Labels:   []string{"Mon", "Tue", "Wed"},
		Datasets: []float64{100, 200, 300},
	}, nil)

	// expected expected response
	expResp := map[string]interface{}{
		"chart_data": map[string]interface{}{
			"labels":   []interface{}{"Mon", "Tue", "Wed"},
			"datasets": []interface{}{100.0, 200.0, 300.0},
		},
	}

	// action
	s.transactionHlr.GetBarChartData(res, req)

	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func getBarChartData_NoStartDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?end_date=2024-03-08&type=expense&time_range=one_week_day", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"error": "start date must be in YYYY-MM-DD format",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getBarChartData_NoEndDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&type=expense&time_range=one_week_day", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"error": "end date must be in YYYY-MM-DD format",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getBarChartData_StartDateAfterEndDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-08&end_date=2024-03-01&type=expense&time_range=one_week_day", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"start_date": "start date must be before end date",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getBarChartData_NoType_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&end_date=2024-03-08&time_range=six_months", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"type": "transaction type must be income or expense",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getBarChartData_NoTimeRangeType_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&end_date=2024-03-08&type=income&time_range=xxx", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"time_range": "time range is invalid",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func (s *TransactionSuite) TestGetPieChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return data":                         getPieChartData_NoError_ReturnData,
		"when no start date, return bad request":             getPieChartData_NoStartDate_ReturnBadReq,
		"when no end date, return bad request":               getPieChartData_NoEndDate_ReturnBadReq,
		"when start date after end date, return bad request": getPieChartData_StartDateAfterEndDate_ReturnBadReq,
		"when no type, return bad request":                   getPieChartData_NoType_ReturnBadReq,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getPieChartData_NoError_ReturnData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-08")
	s.Require().NoError(err, desc)
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetPieChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/pie-chart?start_date=2024-03-01&end_date=2024-03-08&type=expense", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// mock service
	s.mockTransactionUC.On("GetPieChartData",
		req.Context(),
		domain.ChartDateRange{
			Start: start,
			End:   end,
		},
		domain.TransactionTypeExpense,
		user,
	).Return(domain.ChartData{
		Labels:   []string{"2024-03-01", "2024-03-02", "2024-03-03"},
		Datasets: []float64{100, 200, 300},
	}, nil)

	// expected expected response
	expResp := map[string]interface{}{
		"chart_data": map[string]interface{}{
			"labels":   []interface{}{"2024-03-01", "2024-03-02", "2024-03-03"},
			"datasets": []interface{}{100.0, 200.0, 300.0},
		},
	}

	// action
	s.transactionHlr.GetPieChartData(res, req)

	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func getPieChartData_NoStartDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetPieChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/pie-chart?end_date=2024-03-08&type=expense", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetPieChartData(res, req)

	expResp := map[string]interface{}{
		"error": "start date must be in YYYY-MM-DD format",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getPieChartData_NoEndDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetPieChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/pie-chart?start_date=2024-03-01&type=expense", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetPieChartData(res, req)

	expResp := map[string]interface{}{
		"error": "end date must be in YYYY-MM-DD format",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getPieChartData_StartDateAfterEndDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetPieChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/pie-chart?start_date=2024-03-08&end_date=2024-03-01&type=expense", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetPieChartData(res, req)

	expResp := map[string]interface{}{
		"start_date": "start date must be before end date",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getPieChartData_NoType_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetPieChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/pie-chart?start_date=2024-03-01&end_date=2024-03-08", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetPieChartData(res, req)

	expResp := map[string]interface{}{
		"type": "transaction type must be income or expense",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func (s *TransactionSuite) TestGetMonthlyData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return data":                          getMonthlyData_NoError_ReturnData,
		"when start dats is wrong format, return bad request": getMonthlyData_StartDateWrongFormat_ReturnBadReq,
		"when end date is wrong format, return bad request":   getMonthlyData_EndDateWrongFormat_ReturnBadReq,
		"when start date after end date, return bad request":  getMonthlyData_StartDateAfterEndDate_ReturnBadReq,
		"when get monthly data failed, return internal error": getMonthlyData_GetMonthylDataFail_ReturnInternalServerError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getMonthlyData_NoError_ReturnData(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-08")
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetMonthlyData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/monthly-data?start_date=2024-03-01&end_date=2024-03-08", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	montlyData := []domain.TransactionType{
		domain.TransactionTypeIncome,
		domain.TransactionTypeExpense,
		domain.TransactionTypeUnSpecified,
		domain.TransactionTypeBoth,
	}

	// mock service
	s.mockTransactionUC.On("GetMonthlyData",
		req.Context(), dateRange, user,
	).Return(montlyData, nil)

	expResp := map[string]interface{}{
		"monthly_data": []interface{}{
			"income",
			"expense",
			"no data",
			"both",
		},
	}

	// action
	s.transactionHlr.GetMonthlyData(res, req)

	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func getMonthlyData_StartDateWrongFormat_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetMonthlyData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/monthly-data?start_date=2024/03/01&end_date=2024-03-08", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetMonthlyData(res, req)

	expResp := map[string]interface{}{
		"error": "start date must be in YYYY-MM-DD format",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getMonthlyData_EndDateWrongFormat_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetMonthlyData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/monthly-data?start_date=2024-03-01&end_date=2024/03/08", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetMonthlyData(res, req)

	expResp := map[string]interface{}{
		"error": "end date must be in YYYY-MM-DD format",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getMonthlyData_StartDateAfterEndDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetMonthlyData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/monthly-data?start_date=2024-03-08&end_date=2024-03-01", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetMonthlyData(res, req)

	expResp := map[string]interface{}{
		"start_date": "start date must be before end date",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getMonthlyData_GetMonthylDataFail_ReturnInternalServerError(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-08")
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetMonthlyData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/monthly-data?start_date=2024-03-01&end_date=2024-03-08", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expResult := map[string]interface{}{
		"error": "error",
	}

	// mock service
	s.mockTransactionUC.On("GetMonthlyData",
		req.Context(), dateRange, user,
	).Return(nil, errors.New("error"))

	// action
	s.transactionHlr.GetMonthlyData(res, req)

	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, responseBody, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}
