package transaction_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		"when no error, return data":                          getBarChartData_NoError_ReturnData,
		"when no start date, return bad request":              getBarChartData_NoStartDate_ReturnBadReq,
		"when no end date, return bad request":                getBarChartData_NoEndDate_ReturnBadReq,
		"when start date before end date, return bad request": getBarChartData_StartDateBeforeEndDate_ReturnBadReq,
		"when no type, return bad request":                    getBarChartData_NoType_ReturnBadReq,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getBarChartData_NoError_ReturnData(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&end_date=2024-03-08&type=bar", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// mock service
	s.mockTransactionUC.On("GetBarChartData",
		req.Context(),
		domain.ChartTypeBar,
		domain.ChartDateRange{
			StartDate: "2024-03-01",
			EndDate:   "2024-03-08",
		},
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
	s.transactionHlr.GetBarChartData(res, req)

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
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
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?end_date=2024-03-08&type=bar", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"start_date": "Start date must be in YYYY-MM-DD format",
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
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&type=bar", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"start_date": "Start date must be before end date",
		"end_date":   "End date must be in YYYY-MM-DD format",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getBarChartData_StartDateBeforeEndDate_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.GetBarChartData))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-08&end_date=2024-03-01&type=bar", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"start_date": "Start date must be before end date",
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
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/transaction/chart?start_date=2024-03-01&end_date=2024-03-08", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// action
	s.transactionHlr.GetBarChartData(res, req)

	expResp := map[string]interface{}{
		"type": "Chart type must be bar, pie or line",
	}

	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}
