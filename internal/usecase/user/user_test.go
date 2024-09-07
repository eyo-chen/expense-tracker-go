package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/auth"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX = context.Background()
)

type UserSuite struct {
	suite.Suite
	userUC    interfaces.UserUC
	mockUser  *mocks.UserModel
	mockRedis *mocks.RedisService
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	logger.Register()
}

func (s *UserSuite) SetupTest() {
	s.mockUser = mocks.NewUserModel(s.T())
	s.mockRedis = mocks.NewRedisService(s.T())
	s.userUC = New(s.mockUser, s.mockRedis)
}

func (s *UserSuite) TearDownTest() {
	s.mockUser.AssertExpectations(s.T())
}

func (s *UserSuite) TestSignup() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when email not exists, signup successfully": singup_EmailNotExists_SignupSuccessfully,
		"when email exists, return error":            singup_EmailExists_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func singup_EmailNotExists_SignupSuccessfully(s *UserSuite, desc string) {
	s.mockUser.On("FindByEmail", "email.com").Return(domain.User{}, domain.ErrEmailNotFound).Once()
	s.mockUser.On("Create", "username", "email.com", mock.Anything).Return(nil).Once()
	s.mockUser.On("FindByEmail", "email.com").
		Return(domain.User{
			ID:       1,
			Name:     "username",
			Email:    "email.com",
			Password: "password",
		}, nil).Once()

	input := domain.User{
		Name:     "username",
		Email:    "email.com",
		Password: "password",
	}
	token, err := s.userUC.Signup(input)
	s.Require().NoError(err, desc)
	s.Require().NotEmpty(token, desc)
}

func singup_EmailExists_ReturnError(s *UserSuite, desc string) {
	userByEmail := domain.User{
		ID:    1,
		Name:  "username",
		Email: "email.com",
	}
	s.mockUser.On("FindByEmail", "email.com").Return(userByEmail, nil).Once()

	input := domain.User{
		Name:     "username",
		Email:    "email.com",
		Password: "password",
	}
	token, err := s.userUC.Signup(input)
	s.Require().Equal(domain.ErrEmailAlreadyExists, err, desc)
	s.Require().Empty(token, desc)
}

func (s *UserSuite) TestLogin() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when no error, return successfully":    login_NoError_ReturnSuccessfully,
		"when email not exists, return error":   login_EmailNotExists_ReturnError,
		"when password not match, return error": login_PasswordNotMatch_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func login_NoError_ReturnSuccessfully(s *UserSuite, desc string) {
	hashedPassword, err := auth.GenerateHashPassword("password")
	s.Require().NoError(err)

	userByEmail := domain.User{
		ID:            1,
		Name:          "username",
		Email:         "email.com",
		Password:      "password",
		Password_hash: hashedPassword,
	}
	s.mockUser.On("FindByEmail", "email.com").Return(userByEmail, nil).Once()

	input := domain.User{
		Email:    "email.com",
		Password: "password",
	}
	token, err := s.userUC.Login(input)
	s.Require().NoError(err, desc)
	s.Require().NotEmpty(token, desc)
}

func login_EmailNotExists_ReturnError(s *UserSuite, desc string) {
	s.mockUser.On("FindByEmail", "email.com").Return(domain.User{}, domain.ErrEmailNotFound).Once()

	input := domain.User{
		Email:    "email.com",
		Password: "password",
	}
	token, err := s.userUC.Login(input)
	s.Require().Equal(domain.ErrAuthentication, err, desc)
	s.Require().Empty(token, desc)
}

func login_PasswordNotMatch_ReturnError(s *UserSuite, desc string) {
	hashedPassword, err := auth.GenerateHashPassword("password")
	s.Require().NoError(err)

	userByEmail := domain.User{
		ID:            1,
		Name:          "username",
		Email:         "email.com",
		Password:      "password",
		Password_hash: hashedPassword,
	}
	s.mockUser.On("FindByEmail", "email.com").Return(userByEmail, nil).Once()

	input := domain.User{
		Email:    "email.com",
		Password: "password2", // wrong password
	}
	token, err := s.userUC.Login(input)
	s.Require().Equal(domain.ErrAuthentication, err, desc)
	s.Require().Empty(token, desc)
}

func (s *UserSuite) TestToken() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when no error, return successfully":    token_NoError_ReturnSuccessfully,
		"when redis get del fail, return error": token_RedisGetDelFail_ReturnError,
		"when user not found, return error":     token_UserNotFound_ReturnError,
		"when set fail, dont return error":      token_SetFail_DontReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func token_NoError_ReturnSuccessfully(s *UserSuite, desc string) {
	mockRefreshToken := "refresh_token"
	mockEmail := "email.com"
	mockUser := domain.User{Name: "username", Email: mockEmail}
	mockAccessToken, err := genJWTToken(mockUser)
	s.Require().NoError(err)
	mockNewRefreshToken, err := genRefreshToken()
	s.Require().NoError(err)

	s.mockRedis.On("GetDel", mockCTX, hashToken(mockRefreshToken)).
		Return(mockEmail, nil).Once()

	s.mockUser.On("FindByEmail", mockEmail).Return(mockUser, nil).Once()

	s.mockRedis.On("Set", mockCTX, hashToken(mockNewRefreshToken), mockEmail, 7*24*time.Hour).Return(nil).Once()

	expResp := domain.Token{
		Access:  mockAccessToken,
		Refresh: mockNewRefreshToken,
	}

	token, err := s.userUC.Token(mockCTX, mockRefreshToken)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, token, desc)
}

func token_RedisGetDelFail_ReturnError(s *UserSuite, desc string) {
	mockRefreshToken := "refresh_token"
	mockErr := errors.New("redis get del fail")

	s.mockRedis.On("GetDel", mockCTX, hashToken(mockRefreshToken)).
		Return("", mockErr).Once()

	token, err := s.userUC.Token(mockCTX, mockRefreshToken)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Empty(token, desc)
}

func token_UserNotFound_ReturnError(s *UserSuite, desc string) {
	mockRefreshToken := "refresh_token"
	mockEmail := "email.com"
	mockErr := errors.New("find by email fail")

	s.mockRedis.On("GetDel", mockCTX, hashToken(mockRefreshToken)).
		Return(mockEmail, nil).Once()

	s.mockUser.On("FindByEmail", mockEmail).Return(domain.User{}, mockErr).Once()

	token, err := s.userUC.Token(mockCTX, mockRefreshToken)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Empty(token, desc)
}

func token_SetFail_DontReturnError(s *UserSuite, desc string) {
	mockRefreshToken := "refresh_token"
	mockEmail := "email.com"
	mockUser := domain.User{Name: "username", Email: mockEmail}
	mockErr := errors.New("set fail")
	mockAccessToken, err := genJWTToken(mockUser)
	s.Require().NoError(err)
	mockNewRefreshToken, err := genRefreshToken()
	s.Require().NoError(err)

	s.mockRedis.On("GetDel", mockCTX, hashToken(mockRefreshToken)).
		Return(mockEmail, nil).Once()

	s.mockUser.On("FindByEmail", mockEmail).Return(mockUser, nil).Once()

	s.mockRedis.On("Set", mockCTX, hashToken(mockNewRefreshToken), mockEmail, 7*24*time.Hour).Return(mockErr).Once()

	expResp := domain.Token{
		Access:  mockAccessToken,
		Refresh: mockNewRefreshToken,
	}

	token, err := s.userUC.Token(mockCTX, mockRefreshToken)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, token, desc)
}

func (s *UserSuite) TestGetInfo() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when no error, return successfully": getInfo_NoError_ReturnSuccessfully,
		"when get fail, return error":        getInfo_GetFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getInfo_NoError_ReturnSuccessfully(s *UserSuite, desc string) {
	userByID := domain.User{
		ID:    1,
		Name:  "username",
		Email: "email.com",
	}
	s.mockUser.On("GetInfo", int64(1)).Return(userByID, nil).Once()

	user, err := s.userUC.GetInfo(1)
	s.Require().NoError(err, desc)
	s.Require().Equal(userByID, user, desc)
}

func getInfo_GetFail_ReturnError(s *UserSuite, desc string) {
	s.mockUser.On("GetInfo", int64(1)).Return(domain.User{}, domain.ErrUserIDNotFound).Once()

	user, err := s.userUC.GetInfo(1)
	s.Require().Equal(domain.ErrUserIDNotFound, err, desc)
	s.Require().Empty(user, desc)
}
