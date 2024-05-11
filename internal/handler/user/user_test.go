package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	userHlr    *UserHandler
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
	s.userHlr = NewUserHandler(s.mockUserUC)
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
	srv := httptest.NewServer(http.HandlerFunc(s.userHlr.Signup))
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
	s.userHlr.Signup(res, req)

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
	srv := httptest.NewServer(http.HandlerFunc(s.userHlr.Signup))
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
	s.userHlr.Signup(res, req)

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
	srv := httptest.NewServer(http.HandlerFunc(s.userHlr.Signup))
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
	s.userHlr.Signup(res, req)

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
	srv := httptest.NewServer(http.HandlerFunc(s.userHlr.Signup))
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
	s.userHlr.Signup(res, req)

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
	srv := httptest.NewServer(http.HandlerFunc(s.userHlr.Signup))
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
	s.userHlr.Signup(res, req)

	// assertion
	var responseBody map[string]interface{}
	err = json.Unmarshal(res.Body.Bytes(), &responseBody)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, responseBody, desc)
	s.Require().Equal(http.StatusBadRequest, res.Code, desc)
}
