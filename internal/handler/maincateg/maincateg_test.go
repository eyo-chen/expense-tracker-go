package maincateg

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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type MainCategSuite struct {
	suite.Suite
	hlr             *Hlr
	mockMainCategUC *mocks.MainCategUC
}

func TestMainCategSuite(t *testing.T) {
	suite.Run(t, new(MainCategSuite))
}

func (s *MainCategSuite) SetupSuite() {
	logger.Register()
}

func (s *MainCategSuite) SetupTest() {
	s.mockMainCategUC = new(mocks.MainCategUC)
	s.hlr = New(s.mockMainCategUC)
}

func (s *MainCategSuite) TearDownTest() {
	s.mockMainCategUC.AssertExpectations(s.T())
}

func (s *MainCategSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, return successfully":         create_NoError_CreateSuccessfully,
		"when invalid type, return bad request":      create_InvalidType_ReturnBadRequest,
		"when invalid icon type, return bad request": create_InvalidIconType_ReturnBadRequest,
		"when invalid icon data, return bad request": create_InvalidIconData_ReturnBadRequest,
		"when invalid name, return bad request":      create_InvalidName_ReturnBadRequest,
		"when bad request error, return bad request": create_BadRequestError_ReturnBadRequest,
		"when server error, return server error":     create_ServerError_ReturnServerError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_CreateSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockCateg := domain.MainCateg{
		Name:     "Food",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// mock service
	s.mockMainCategUC.On("Create", mockCateg, mockUser.ID).Return(nil)

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusCreated, res.Code, desc)
	s.Require().Empty(responseBody, desc)
}

func create_InvalidType_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "invalid",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"type": "Type must be income or expense",
	}

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func create_InvalidIconType_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "invalid",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"icon_type": "Icon type must be default or custom",
	}

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func create_InvalidIconData_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"icon_data": "Icon data can't be empty",
	}

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func create_InvalidName_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"name": "Name can't be empty",
	}

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func create_BadRequestError_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockCateg := domain.MainCateg{
		Name:     "Food",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"error": "icon not found",
	}

	// mock service
	s.mockMainCategUC.On("Create", mockCateg, mockUser.ID).Return(domain.ErrIconNotFound)

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func create_ServerError_ReturnServerError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockCateg := domain.MainCateg{
		Name:     "Food",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	mockUser := domain.User{ID: 1}
	mockErr := errors.New("server error")

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
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
	s.mockMainCategUC.On("Create", mockCateg, mockUser.ID).Return(mockErr)

	// action
	s.hlr.Create(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func (s *MainCategSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, return successfully":        getAll_NoError_GetAllSuccessfully,
		"when in correct type, return successfully": getAll_InCorrectType_GetAllSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAll_NoError_GetAllSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock data
	mockCategs := []domain.MainCateg{
		{
			ID:       1,
			Name:     "Food",
			Type:     domain.TransactionTypeIncome,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
		{
			ID:       2,
			Name:     "Transportation",
			Type:     domain.TransactionTypeExpense,
			IconType: domain.IconTypeCustom,
			IconData: "url",
		},
	}
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetAll))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/main-category?type=income", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"categories": []interface{}{
			map[string]interface{}{
				"id":        float64(1),
				"name":      "Food",
				"type":      "income",
				"icon_type": "default",
				"icon_data": "url",
			},
			map[string]interface{}{
				"id":        float64(2),
				"name":      "Transportation",
				"type":      "expense",
				"icon_type": "custom",
				"icon_data": "url",
			},
		},
	}

	// mock service
	s.mockMainCategUC.On("GetAll", req.Context(), mockUser.ID, domain.TransactionTypeIncome).
		Return(mockCategs, nil)

	// action
	s.hlr.GetAll(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func getAll_InCorrectType_GetAllSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock data
	mockCategs := []domain.MainCateg{}
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetAll))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/main-category?type=aaaa", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"categories": []interface{}{},
	}

	// mock service
	s.mockMainCategUC.On("GetAll", req.Context(), mockUser.ID, domain.TransactionTypeUnSpecified).
		Return(mockCategs, nil)

	// action
	s.hlr.GetAll(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func (s *MainCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, return successfully":         update_NoError_CreateSuccessfully,
		"when invalid id, return bad request":        update_InvalidID_CreateSuccessfully,
		"when invalid type, return bad request":      update_InvalidType_ReturnBadRequest,
		"when invalid icon type, return bad request": update_InvalidIconType_ReturnBadRequest,
		"when invalid icon data, return bad request": update_InvalidIconData_ReturnBadRequest,
		"when invalid name, return bad request":      update_InvalidName_ReturnBadRequest,
		"when bad request error, return bad request": update_BadRequestError_ReturnBadRequest,
		"when server error, return server error":     update_ServerError_ReturnServerError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoError_CreateSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockCateg := domain.MainCateg{
		ID:       1,
		Name:     "Food",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPut, srv.URL+"/v1/main-category/id", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// mock service
	s.mockMainCategUC.On("Update", mockCateg, mockUser.ID).Return(nil)

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
	s.Require().Empty(responseBody, desc)
}

func update_InvalidID_CreateSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPut, srv.URL+"/v1/main-category/id", bytes.NewBuffer([]byte("")))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "aaa"})

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().NotEmpty(responseBody, desc)
}

func update_InvalidType_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "invalid",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category/id", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"type": "Type must be income or expense",
	}

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func update_InvalidIconType_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "invalid",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category/id", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"icon_type": "Icon type must be default or custom",
	}

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func update_InvalidIconData_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category/id", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"icon_data": "Icon data can't be empty",
	}

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func update_InvalidName_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category/id", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"name": "Name can't be empty",
	}

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func update_BadRequestError_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockCateg := domain.MainCateg{
		ID:       1,
		Name:     "Food",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	mockUser := domain.User{ID: 1}

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category/id", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"error": "icon not found",
	}

	// mock service
	s.mockMainCategUC.On("Update", mockCateg, mockUser.ID).Return(domain.ErrIconNotFound)

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func update_ServerError_ReturnServerError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockInput := map[string]interface{}{
		"name":      "Food",
		"type":      "income",
		"icon_type": "default",
		"icon_data": "url",
	}
	mockBody, err := json.Marshal(mockInput)
	s.Require().NoError(err, desc)
	mockCateg := domain.MainCateg{
		ID:       1,
		Name:     "Food",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	mockUser := domain.User{ID: 1}
	mockErr := errors.New("server error")

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Create))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/main-category", bytes.NewBuffer(mockBody))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// set context value on request
	req = ctxutil.SetUser(req, &mockUser)

	// prepare expResp
	expResp := map[string]interface{}{
		"error": mockErr.Error(),
	}

	// mock service
	s.mockMainCategUC.On("Update", mockCateg, mockUser.ID).Return(mockErr)

	// action
	s.hlr.Update(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}

func (s *MainCategSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, return successfully":         delete_NoError_DeleteSuccessfully,
		"when bad request error, return bad request": delete_BadRequestError_ReturnBadRequest,
		"when server error, return server error":     delete_ServerError_ReturnServerError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func delete_NoError_DeleteSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Delete))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/main-category/id", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// mock service
	s.mockMainCategUC.On("Delete", int64(1)).Return(nil)

	// action
	s.hlr.Delete(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
	s.Require().Empty(responseBody, desc)
}

func delete_BadRequestError_ReturnBadRequest(s *MainCategSuite, desc string) {
	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Delete))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/main-category/id", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "aaaa"})

	// action
	s.hlr.Delete(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
	s.Require().NotEmpty(responseBody, desc)
}

func delete_ServerError_ReturnServerError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockErr := errors.New("server error")

	// prepare mock request
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Delete))
	req := httptest.NewRequest(http.MethodDelete, srv.URL+"/v1/main-category/id", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// prepare expResp
	expResp := map[string]interface{}{
		"error": mockErr.Error(),
	}

	// mock service
	s.mockMainCategUC.On("Delete", int64(1)).Return(mockErr)

	// action
	s.hlr.Delete(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
	s.Require().Equal(expResp, responseBody, desc)
}
