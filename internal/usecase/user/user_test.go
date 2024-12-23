package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
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
	uc           *UC
	mockUserRepo *mocks.UserRepo
	mockRedis    *mocks.RedisService
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	logger.Register()
}

func (s *UserSuite) SetupTest() {
	s.mockUserRepo = mocks.NewUserRepo(s.T())
	s.mockRedis = mocks.NewRedisService(s.T())
	s.uc = New(s.mockUserRepo, s.mockRedis)
}

func (s *UserSuite) TearDownTest() {
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserSuite) TestSignup() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when email not exists, signup successfully": singup_EmailNotExists_SignupSuccessfully,
		"when email exists, return error":            singup_EmailExists_ReturnError,
		"when set cache fail, dont return error":     singup_SetCacheFail_DontReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func singup_EmailNotExists_SignupSuccessfully(s *UserSuite, desc string) {
	// prepare mock data
	mockUser := domain.User{
		ID:    1,
		Name:  "username",
		Email: "email.com",
	}

	// prepare mock service
	s.mockUserRepo.On("FindByEmail", "email.com").Return(domain.User{}, domain.ErrEmailNotFound).Once()
	s.mockUserRepo.On("Create", "username", "email.com", mock.Anything).Return(nil).Once()
	s.mockUserRepo.On("FindByEmail", "email.com").Return(mockUser, nil).Once()
	s.mockRedis.On("Set", mockCTX, mock.Anything, mockUser.Email, 7*24*time.Hour).Return(nil).Once()

	token, err := s.uc.Signup(mockCTX, mockUser)
	s.Require().NoError(err, desc)
	s.Require().NotEmpty(token.Access, desc)
	s.Require().NotEmpty(token.Refresh, desc)
}

func singup_EmailExists_ReturnError(s *UserSuite, desc string) {
	mockUser := domain.User{
		ID:    1,
		Name:  "username",
		Email: "email.com",
	}
	s.mockUserRepo.On("FindByEmail", "email.com").Return(mockUser, nil).Once()

	token, err := s.uc.Signup(mockCTX, mockUser)
	s.Require().ErrorIs(err, domain.ErrEmailAlreadyExists, desc)
	s.Require().Empty(token, desc)
}

func singup_SetCacheFail_DontReturnError(s *UserSuite, desc string) {
	// prepare mock data
	mockUser := domain.User{
		ID:    1,
		Name:  "username",
		Email: "email.com",
	}

	// prepare mock service
	s.mockUserRepo.On("FindByEmail", "email.com").Return(domain.User{}, domain.ErrEmailNotFound).Once()
	s.mockUserRepo.On("Create", "username", "email.com", mock.Anything).Return(nil).Once()
	s.mockUserRepo.On("FindByEmail", "email.com").Return(mockUser, nil).Once()
	s.mockRedis.On("Set", mockCTX, mock.Anything, mockUser.Email, 7*24*time.Hour).Return(errors.New("set fail")).Once()

	token, err := s.uc.Signup(mockCTX, mockUser)
	s.Require().NoError(err, desc)
	s.Require().NotEmpty(token.Access, desc)
	s.Require().NotEmpty(token.Refresh, desc)
}

func (s *UserSuite) TestLogin() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when no error, return successfully":     login_NoError_ReturnSuccessfully,
		"when email not exists, return error":    login_EmailNotExists_ReturnError,
		"when password not match, return error":  login_PasswordNotMatch_ReturnError,
		"when redis set fail, dont return error": login_RedisSetFail_DontReturnError,
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

	mockUser := domain.User{
		ID:            1,
		Name:          "username",
		Email:         "email.com",
		Password:      "password",
		Password_hash: hashedPassword,
	}
	s.mockUserRepo.On("FindByEmail", "email.com").Return(mockUser, nil).Once()
	s.mockRedis.On("Set", mockCTX, mock.Anything, mockUser.Email, 7*24*time.Hour).Return(nil).Once()

	token, err := s.uc.Login(mockCTX, mockUser)
	s.Require().NoError(err, desc)
	s.Require().NotEmpty(token.Access, desc)
	s.Require().NotEmpty(token.Refresh, desc)
}

func login_EmailNotExists_ReturnError(s *UserSuite, desc string) {
	s.mockUserRepo.On("FindByEmail", "email.com").Return(domain.User{}, domain.ErrEmailNotFound).Once()

	input := domain.User{
		Email:    "email.com",
		Password: "password",
	}
	token, err := s.uc.Login(mockCTX, input)
	s.Require().ErrorIs(err, domain.ErrAuthentication, desc)
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
	s.mockUserRepo.On("FindByEmail", "email.com").Return(userByEmail, nil).Once()

	input := domain.User{
		Email:    "email.com",
		Password: "password2", // wrong password
	}
	token, err := s.uc.Login(mockCTX, input)
	s.Require().ErrorIs(err, domain.ErrAuthentication, desc)
	s.Require().Empty(token, desc)
}

func login_RedisSetFail_DontReturnError(s *UserSuite, desc string) {
	hashedPassword, err := auth.GenerateHashPassword("password")
	s.Require().NoError(err)

	mockUser := domain.User{
		ID:            1,
		Name:          "username",
		Email:         "email.com",
		Password:      "password",
		Password_hash: hashedPassword,
	}
	s.mockUserRepo.On("FindByEmail", "email.com").Return(mockUser, nil).Once()
	s.mockRedis.On("Set", mockCTX, mock.Anything, mockUser.Email, 7*24*time.Hour).Return(errors.New("set fail")).Once()

	token, err := s.uc.Login(mockCTX, mockUser)
	s.Require().NoError(err, desc)
	s.Require().NotEmpty(token.Access, desc)
	s.Require().NotEmpty(token.Refresh, desc)
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

	s.mockUserRepo.On("FindByEmail", mockEmail).Return(mockUser, nil).Once()

	s.mockRedis.On("Set", mockCTX, hashToken(mockNewRefreshToken), mockEmail, 7*24*time.Hour).Return(nil).Once()

	expResp := domain.Token{
		Access:  mockAccessToken,
		Refresh: mockNewRefreshToken,
	}

	token, err := s.uc.Token(mockCTX, mockRefreshToken)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, token, desc)
}

func token_RedisGetDelFail_ReturnError(s *UserSuite, desc string) {
	mockRefreshToken := "refresh_token"
	mockErr := errors.New("redis get del fail")

	s.mockRedis.On("GetDel", mockCTX, hashToken(mockRefreshToken)).
		Return("", mockErr).Once()

	token, err := s.uc.Token(mockCTX, mockRefreshToken)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Empty(token, desc)
}

func token_UserNotFound_ReturnError(s *UserSuite, desc string) {
	mockRefreshToken := "refresh_token"
	mockEmail := "email.com"
	mockErr := errors.New("find by email fail")

	s.mockRedis.On("GetDel", mockCTX, hashToken(mockRefreshToken)).
		Return(mockEmail, nil).Once()

	s.mockUserRepo.On("FindByEmail", mockEmail).Return(domain.User{}, mockErr).Once()

	token, err := s.uc.Token(mockCTX, mockRefreshToken)
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

	s.mockUserRepo.On("FindByEmail", mockEmail).Return(mockUser, nil).Once()

	s.mockRedis.On("Set", mockCTX, hashToken(mockNewRefreshToken), mockEmail, 7*24*time.Hour).Return(mockErr).Once()

	expResp := domain.Token{
		Access:  mockAccessToken,
		Refresh: mockNewRefreshToken,
	}

	token, err := s.uc.Token(mockCTX, mockRefreshToken)
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
	s.mockUserRepo.On("GetInfo", int64(1)).Return(userByID, nil).Once()

	user, err := s.uc.GetInfo(1)
	s.Require().NoError(err, desc)
	s.Require().Equal(userByID, user, desc)
}

func getInfo_GetFail_ReturnError(s *UserSuite, desc string) {
	s.mockUserRepo.On("GetInfo", int64(1)).Return(domain.User{}, domain.ErrUserIDNotFound).Once()

	user, err := s.uc.GetInfo(1)
	s.Require().ErrorIs(err, domain.ErrUserIDNotFound, desc)
	s.Require().Empty(user, desc)
}
