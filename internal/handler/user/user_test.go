package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	hlr        *Hlr
	mockUserUC *mocks.UserUC
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	logger.Register()
}

func (s *UserSuite) SetupTest() {
	s.mockUserUC = mocks.NewUserUC(s.T())
	s.hlr = NewUserHandler(s.mockUserUC)
}

func (s *UserSuite) TearDownTest() {
	s.mockUserUC.AssertExpectations(s.T())
}

func (s *UserSuite) TestSignup() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when no error, create successfully":  signup_NoError_CreateSuccessfully,
		"when empty name, return error":       signup_EmptyName_ReturnError,
		"when invalid email, return error":    signup_InvalidEmail_ReturnError,
		"when invalid password, return error": signup_InvalidPasswordReturnError,
		"when email exists, return error":     signup_EmailExists_CreateSuccessfully,
		"when signup fail, return error":      signup_SignupFail_CreateSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func signup_NoError_CreateSuccessfully(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Name:     "username",
		Email:    "email@gmail.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/signup", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	s.mockUserUC.On("Signup", user).Return("token", nil).Once()

	// prepare expected response
	expResp := map[string]interface{}{
		"token": "token",
	}

	// action
	s.hlr.Signup(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusCreated, res.Code, desc)
}

func signup_EmptyName_ReturnError(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Name:     "",
		Email:    "email@gmail.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/signup", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare expected response
	expResp := map[string]interface{}{
		"name": "Name can't be empty",
	}

	// action
	s.hlr.Signup(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func signup_InvalidEmail_ReturnError(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Name:     "username",
		Email:    "email.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/signup", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare expected response
	expResp := map[string]interface{}{
		"email": "Invalid email address",
	}

	// action
	s.hlr.Signup(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func signup_InvalidPasswordReturnError(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Name:     "username",
		Email:    "email@gmail.com",
		Password: "p",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/signup", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare expected response
	expResp := map[string]interface{}{
		"password": "Password must be at least 8 characters long",
	}

	// action
	s.hlr.Signup(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func signup_EmailExists_CreateSuccessfully(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Name:     "username",
		Email:    "email@gmail.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/signup", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	s.mockUserUC.On("Signup", user).Return("", domain.ErrEmailAlreadyExists).Once()

	// prepare expected response
	expResp := map[string]interface{}{
		"error": domain.ErrEmailAlreadyExists.Error(),
	}

	// action
	s.hlr.Signup(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func signup_SignupFail_CreateSuccessfully(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Name:     "username",
		Email:    "email@gmail.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/signup", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	s.mockUserUC.On("Signup", user).Return("", errors.New("error")).Once()

	// prepare expected response
	expResp := map[string]interface{}{
		"error": "error",
	}

	// action
	s.hlr.Signup(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}

func (s *UserSuite) TestLogin() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when no error, login successfully":   login_NoError_LoginSuccessfully,
		"when invalid email, return error":    login_InvalidEmail_ReturnError,
		"when invalid password, return error": login_InvalidPassword_ReturnError,
		"when auth error, return error":       login_AuthError_ReturnError,
		"when login fail, return error":       login_LoginFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func login_NoError_LoginSuccessfully(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Email:    "email@gmail.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/login", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	s.mockUserUC.On("Login", user).Return("token", nil).Once()

	// prepare expected response
	expResp := map[string]interface{}{
		"token": "token",
	}

	// action
	s.hlr.Login(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func login_InvalidEmail_ReturnError(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Email:    "email.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/login", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare expected response
	expResp := map[string]interface{}{
		"email": "Invalid email address",
	}

	// action
	s.hlr.Login(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func login_InvalidPassword_ReturnError(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Email:    "email@gmail.com",
		Password: "p",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/login", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare expected response
	expResp := map[string]interface{}{
		"password": "Password must be at least 8 characters long",
	}

	// action
	s.hlr.Login(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func login_AuthError_ReturnError(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Email:    "email@gmail.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/login", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	s.mockUserUC.On("Login", user).Return("", domain.ErrAuthentication).Once()

	// prepare expected response
	expResp := map[string]interface{}{
		"error": domain.ErrAuthentication.Error(),
	}

	// action
	s.hlr.Login(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusUnauthorized, res.Code, desc)
}

func login_LoginFail_ReturnError(s *UserSuite, desc string) {
	// prepare data
	user := domain.User{
		Email:    "email@gmail.com",
		Password: "password",
	}
	body, err := json.Marshal(user)
	s.Require().NoError(err, desc)

	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.Signup))
	req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/user/login", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	s.mockUserUC.On("Login", user).Return("", errors.New("error")).Once()

	// prepare expected response
	expResp := map[string]interface{}{
		"error": "error",
	}

	// action
	s.hlr.Login(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}

func (s *UserSuite) TestGetInfo() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when no error, get info successfully":          getInfo_NoError_GetInfoSuccessfully,
		"when user not found, return bad request error": getInfo_UserNotFound_ReturnBadRequest,
		"when get info fail, return error":              getInfo_GetInfoFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getInfo_NoError_GetInfoSuccessfully(s *UserSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetInfo))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/user/info", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	user := domain.User{
		ID:                1,
		Name:              "username",
		Email:             "aaa@gmail.com",
		IsSetInitCategory: true,
	}

	s.mockUserUC.On("GetInfo", user.ID).Return(user, nil).Once()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// prepare expected response
	expResp := map[string]interface{}{
		"id":                   float64(1),
		"name":                 "username",
		"email":                "aaa@gmail.com",
		"is_set_init_category": true,
	}

	// action
	s.hlr.GetInfo(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusOK, res.Code, desc)
}

func getInfo_UserNotFound_ReturnBadRequest(s *UserSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetInfo))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/user/info", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	user := domain.User{
		ID:    1,
		Name:  "username",
		Email: "aaa@gmail.com",
	}

	s.mockUserUC.On("GetInfo", user.ID).Return(domain.User{}, domain.ErrUserIDNotFound).Once()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// prepare expected response
	expResp := map[string]interface{}{
		"error": domain.ErrUserIDNotFound.Error(),
	}

	// action
	s.hlr.GetInfo(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}

func getInfo_GetInfoFail_ReturnError(s *UserSuite, desc string) {
	// prepare request, and response recorder
	srv := httptest.NewServer(http.HandlerFunc(s.hlr.GetInfo))
	req := httptest.NewRequest(http.MethodGet, srv.URL+"/v1/user/info", nil)
	res := httptest.NewRecorder()
	defer srv.Close()
	defer req.Body.Close()
	defer res.Result().Body.Close()

	// prepare service
	user := domain.User{}

	s.mockUserUC.On("GetInfo", user.ID).Return(domain.User{}, errors.New("error")).Once()

	// set context value on request
	req = ctxutil.SetUser(req, &user)

	// prepare expected response
	expResp := map[string]interface{}{
		"error": "error",
	}

	// action
	s.hlr.GetInfo(res, req)

	// assertion
	var responseBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusInternalServerError, res.Code, desc)
}
