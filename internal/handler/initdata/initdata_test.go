package initdata

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

type InitDataSuite struct {
	suite.Suite
	hlr            *Hlr
	mockInitDataUC *mocks.InitDataUC
}

func TestInitDataSuite(t *testing.T) {
	suite.Run(t, new(InitDataSuite))
}

func (s *InitDataSuite) SetupSuite() {
	logger.Register()
}

func (s *InitDataSuite) SetupTest() {
	s.mockInitDataUC = new(mocks.InitDataUC)
	s.hlr = New(s.mockInitDataUC)
}

func (s *InitDataSuite) TearDownTest() {
	s.mockInitDataUC.AssertExpectations(s.T())
}

func (s *InitDataSuite) TestList() {
	for scenario, fn := range map[string]func(s *InitDataSuite, desc string){
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

func list_NoError_ReturnSuccessfully(s *InitDataSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.List))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/init-data", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// mock service
	s.mockInitDataUC.On("List").Return(
		domain.InitData{
			Expense: []domain.InitDataMainCateg{
				{
					Name: "Food",
					Icon: domain.Icon{
						ID:  1,
						URL: "http://example.com/icon/1",
					},
					SubCategs: []string{
						"breakfast", "brunch", "lunch", "dinner", "groceries", "drink", "snak",
					},
				},
			},
			Income: []domain.InitDataMainCateg{
				{
					Name: "salary",
					Icon: domain.Icon{
						ID:  12,
						URL: "http://example.com/icon/12",
					},
					SubCategs: []string{
						"salary", "bonus", "commission", "tips",
					},
				},
			},
		}, nil,
	)

	// prepare expected response
	expResp := map[string]interface{}{
		"init_data": map[string]interface{}{
			"expense": []interface{}{
				map[string]interface{}{
					"name": "Food",
					"icon": map[string]interface{}{
						"id": float64(1), "url": "http://example.com/icon/1",
					},
					"sub_categories": []interface{}{"breakfast", "brunch", "lunch", "dinner", "groceries", "drink", "snak"}}},
			"income": []interface{}{map[string]interface{}{
				"name": "salary",
				"icon": map[string]interface{}{
					"id": float64(12), "url": "http://example.com/icon/12",
				},
				"sub_categories": []interface{}{"salary", "bonus", "commission", "tips"},
			},
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

func list_ListFail_ReturnError(s *InitDataSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.List))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/init-data", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// mock service
	s.mockInitDataUC.On("List").Return(
		domain.InitData{}, errors.New("list failed"),
	)

	// action
	s.hlr.List(res, req)

	// assertion
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}

func (s *InitDataSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *InitDataSuite, desc string){
		"when no error, return successfully": create_NoError_ReturnSuccessfully,
		"when create failed, return error":   create_CreateFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_ReturnSuccessfully(s *InitDataSuite, desc string) {
	// prepare mock data
	mockData := domain.InitData{
		Expense: []domain.InitDataMainCateg{
			{
				Name:      "Food",
				Icon:      domain.Icon{ID: 1, URL: "http://example.com/icon/1"},
				SubCategs: []string{"breakfast", "brunch", "lunch", "dinner", "groceries", "drink", "snak"},
			},
		},
		Income: []domain.InitDataMainCateg{
			{
				Name:      "salary",
				Icon:      domain.Icon{ID: 12, URL: "http://example.com/icon/12"},
				SubCategs: []string{"salary", "bonus", "commission", "tips"},
			},
		},
	}
	body, err := json.Marshal(mockData)
	s.Require().NoError(err, desc)

	mockUser := domain.User{ID: 1}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/init-data", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// mock service
	s.mockInitDataUC.On("Create", req.Context(), mockData, int64(1)).Return(nil)

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Nil(responseBody, desc)
	s.Require().Equal(http.StatusCreated, res.Code, desc)
}

func create_CreateFail_ReturnError(s *InitDataSuite, desc string) {
	// prepare mock data
	mockData := domain.InitData{
		Expense: []domain.InitDataMainCateg{
			{
				Name:      "Food",
				Icon:      domain.Icon{ID: 1, URL: "http://example.com/icon/1"},
				SubCategs: []string{"breakfast", "brunch", "lunch", "dinner", "groceries", "drink", "snak"},
			},
		},
		Income: []domain.InitDataMainCateg{
			{
				Name:      "salary",
				Icon:      domain.Icon{ID: 12, URL: "http://example.com/icon/12"},
				SubCategs: []string{"salary", "bonus", "commission", "tips"},
			},
		},
	}
	body, err := json.Marshal(mockData)
	s.Require().NoError(err, desc)

	mockUser := domain.User{ID: 1}

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/init-data", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// mock service
	s.mockInitDataUC.On("Create", req.Context(), mockData, int64(1)).Return(errors.New("create failed"))

	// action
	s.hlr.Create(res, req)

	// expected response
	expResp := map[string]interface{}{"error": "create failed"}

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}
