package icon

import (
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

type IconSuite struct {
	suite.Suite

	hlr        *Hlr
	mockIconUC *mocks.IconUC
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupSuite() {
	logger.Register()
}

func (s *IconSuite) SetupTest() {
	s.mockIconUC = new(mocks.IconUC)
	s.hlr = New(s.mockIconUC)
}

func (s *IconSuite) TearDownTest() {
	s.mockIconUC.AssertExpectations(s.T())
}

func (s *IconSuite) TestList() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when no error, return successfully": list_NoError_ReturnSuccessfully,
		"when list failed, return error":     list_ListFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func list_NoError_ReturnSuccessfully(s *IconSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.List))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/icon", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare mock data
	mockIcons := []domain.DefaultIcon{
		{
			ID:  1,
			URL: "http://example.com/icon/1",
		},
		{
			ID:  2,
			URL: "http://example.com/icon/2",
		},
	}
	// mock service
	s.mockIconUC.On("List").Return(mockIcons, nil)

	// prepare expected response
	expResp := map[string]interface{}{
		"icons": []interface{}{
			map[string]interface{}{
				"id":  float64(1),
				"url": "http://example.com/icon/1",
			},
			map[string]interface{}{
				"id":  float64(2),
				"url": "http://example.com/icon/2",
			},
		},
	}

	// action
	s.hlr.List(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func list_ListFail_ReturnError(s *IconSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.List))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/icon", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	mockErr := errors.New("list failed")

	// mock service
	s.mockIconUC.On("List").Return(nil, mockErr)

	// prepare expected error
	expErr := map[string]interface{}{
		"error": mockErr.Error(),
	}

	// action
	s.hlr.List(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expErr, responseBody, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}

func (s *IconSuite) TestListByUserID() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when no error, return successfully": listByUserID_NoError_ReturnSuccessfully,
		"when list failed, return error":     listByUserID_ListFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func listByUserID_NoError_ReturnSuccessfully(s *IconSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.ListByUserID))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/user-icon", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	mockUser := domain.User{ID: 1}
	req = ctxutil.SetUser(req, &mockUser)

	// prepare mock data
	mockIcons := []domain.Icon{
		{
			ID:   1,
			Type: domain.IconTypeCustom,
			URL:  "http://example.com/icon/1",
		},
		{
			ID:   1,
			Type: domain.IconTypeDefault,
			URL:  "http://example.com/icon/2",
		},
	}

	// mock service
	s.mockIconUC.On("ListByUserID", req.Context(), mockUser.ID).Return(mockIcons, nil)

	// prepare expected response
	expResp := map[string]interface{}{
		"icons": []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"type": "custom",
				"url":  "http://example.com/icon/1",
			},
			map[string]interface{}{
				"id":   float64(1),
				"type": "default",
				"url":  "http://example.com/icon/2",
			},
		},
	}

	// action
	s.hlr.ListByUserID(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func listByUserID_ListFail_ReturnError(s *IconSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.ListByUserID))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/icon", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	mockUser := domain.User{ID: 1}
	req = ctxutil.SetUser(req, &mockUser)

	mockErr := errors.New("list failed")

	// mock service
	s.mockIconUC.On("ListByUserID", req.Context(), mockUser.ID).Return(nil, mockErr)

	// prepare expected error
	expErr := map[string]interface{}{
		"error": mockErr.Error(),
	}

	// action
	s.hlr.ListByUserID(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expErr, responseBody, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}
