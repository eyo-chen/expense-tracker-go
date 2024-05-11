package user

import (
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	userUC   *UserUC
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
	s.mockUser.On("FindByEmail", "email.com").Return(nil, nil).Once()
	s.mockUser.On("Create", "username", "email.com", mock.Anything).Return(nil).Once()
	s.mockUser.On("FindByEmail", "email.com").
		Return(&domain.User{
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
	s.Require().NoError(err)
	s.Require().NotEmpty(token)
}

func singup_EmailExists_ReturnError(s *UserSuite, desc string) {
	userByEmail := &domain.User{
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
	s.Require().Equal(domain.ErrDataAlreadyExists, err)
	s.Require().Empty(token)
}
