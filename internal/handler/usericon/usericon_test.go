package usericon

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/ctxutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type UserIconSuite struct {
	suite.Suite
	hlr            *Hlr
	mockUserIconUC *mocks.UserIconUC
}

func TestUserIconSuite(t *testing.T) {
	suite.Run(t, new(UserIconSuite))
}

func (s *UserIconSuite) SetupSuite() {
	logger.Register()
}

func (s *UserIconSuite) SetupTest() {
	s.mockUserIconUC = new(mocks.UserIconUC)
	s.hlr = New(s.mockUserIconUC)
}

func (s *UserIconSuite) TearDownTest() {
	s.mockUserIconUC.AssertExpectations(s.T())
}

func (s *UserIconSuite) TestGetPutObjectURL() {
	for scenario, fn := range map[string]func(s *UserIconSuite, desc string){
		"when no error, return successfully":                getPutObjectURL_NoError_ReturnSuccessfully,
		"when empty file name, return bad request":          getPutObjectURL_EmptyFileName_ReturnBadRequest,
		"when get put object url fail, return server error": getPutObjectURL_GetFail_ReturnServerError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getPutObjectURL_NoError_ReturnSuccessfully(s *UserIconSuite, desc string) {
	// prepare mock data
	mockFileName := "test.png"
	mockInput := map[string]interface{}{
		"file_name": mockFileName,
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)

	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetPutObjectURL))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user-icon/url", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"url": "url",
	}

	// mock service
	s.mockUserIconUC.On("GetPutObjectURL", req.Context(), mockFileName, mockUser.ID).Return("url", nil)

	// action
	s.hlr.GetPutObjectURL(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func getPutObjectURL_EmptyFileName_ReturnBadRequest(s *UserIconSuite, desc string) {
	// prepare mock data
	mockFileName := ""
	mockInput := map[string]interface{}{
		"file_name": mockFileName,
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)

	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetPutObjectURL))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user-icon/url", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"file_name": "File name can't be empty",
	}

	// action
	s.hlr.GetPutObjectURL(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func getPutObjectURL_GetFail_ReturnServerError(s *UserIconSuite, desc string) {
	// prepare mock data
	mockFileName := "test.png"
	mockInput := map[string]interface{}{
		"file_name": mockFileName,
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)

	mockUser := domain.User{ID: 1}
	mockErr := errors.New("get put object url error")

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetPutObjectURL))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user-icon/url", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"error": mockErr.Error(),
	}

	// mock service
	s.mockUserIconUC.On("GetPutObjectURL", req.Context(), mockFileName, mockUser.ID).Return("", mockErr)

	// action
	s.hlr.GetPutObjectURL(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}
