package transaction_test

import (
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
	"github.com/stretchr/testify/mock"
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

	s.mockTransactionUC.On("Delete", mock.Anything, int64(1), user).Return(nil)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.Delete))
	defer srv.Close()
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	defer req.Body.Close()
	res := httptest.NewRecorder()
	defer res.Result().Body.Close()

	// set context value on request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = ctxutil.SetUser(req, &user)

	s.transactionHlr.Delete(res, req)

	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func delete_IncorrectID_ReturnBadReq(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.Delete))
	defer srv.Close()
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	defer req.Body.Close()
	res := httptest.NewRecorder()
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
	defer srv.Close()
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	defer req.Body.Close()
	res := httptest.NewRecorder()
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

	s.mockTransactionUC.On("Delete", mock.Anything, int64(1), user).Return(domain.ErrTransactionDataNotFound)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.transactionHlr.Delete))
	defer srv.Close()
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/transaction/id", nil)
	defer req.Body.Close()
	res := httptest.NewRecorder()
	defer res.Result().Body.Close()

	// set context value on request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = ctxutil.SetUser(req, &user)

	s.transactionHlr.Delete(res, req)

	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}
