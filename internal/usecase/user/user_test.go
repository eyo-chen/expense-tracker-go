package user

import (
	"testing"

	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
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
}

func singup_EmailNotExists_SignupSuccessfully(s *UserSuite, desc string) {}

func singup_EmailExists_ReturnError(s *UserSuite, desc string) {}
