package user

import (
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/auth"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	userUC   interfaces.UserUC
	mockUser *mocks.UserModel
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	logger.Register()
}

func (s *UserSuite) SetupTest() {
	s.mockUser = mocks.NewUserModel(s.T())
	s.userUC = NewUserUC(s.mockUser)
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
