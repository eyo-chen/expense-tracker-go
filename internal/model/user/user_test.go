package user

import (
	"context"
	"database/sql"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/dockerutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil/efactory/db/esql"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx = context.Background()
)

type UserSuite struct {
	suite.Suite
	db      *sql.DB
	migrate *migrate.Migrate
	f       *efactory.Factory[User]
	model   interfaces.UserModel
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db, migrate := testutil.ConnToDB(port)
	s.model = NewUserModel(db)
	logger.Register()
	s.db = db
	s.migrate = migrate
	s.f = efactory.New(User{}).SetConfig(efactory.Config[User]{
		DB: &esql.Config{
			DB: db,
		},
	})
}

func (s *UserSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *UserSuite) SetupTest() {
	s.model = NewUserModel(s.db)
}

func (s *UserSuite) TearDownTest() {
	if _, err := s.db.Exec("DELETE FROM users"); err != nil {
		s.Require().NoError(err)
	}

	s.f.Reset()
}

func (s *UserSuite) TestCreate() {
	user := domain.User{
		Name:          "username",
		Email:         "email.com",
		Password:      "password",
		Password_hash: "password_hash",
	}

	err := s.model.Create(user.Name, user.Email, user.Password_hash)
	s.Require().NoError(err)

	// check if user is created
	checkedUser := User{}
	err = s.db.QueryRow("SELECT name, email, password_hash FROM users WHERE email = ?", user.Email).Scan(&checkedUser.Name, &checkedUser.Email, &checkedUser.Password_hash)
	s.Require().NoError(err)
	s.Require().Equal(user.Name, checkedUser.Name)
	s.Require().Equal(user.Email, checkedUser.Email)
	s.Require().Equal(user.Password_hash, checkedUser.Password_hash)
}

func (s *UserSuite) TestFindByEmail() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when user is found, return successfully": findByEmail_FoundUser_ReturnSuccessfully,
		"when user is not found, return error":    findByEmail_NotFound_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func findByEmail_FoundUser_ReturnSuccessfully(s *UserSuite, desc string) {
	users, err := s.f.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	expResult := domain.User{
		ID:            users[0].ID,
		Name:          users[0].Name,
		Email:         users[0].Email,
		Password_hash: users[0].Password_hash,
	}

	user, err := s.model.FindByEmail(users[0].Email)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, user, desc)
}

func findByEmail_NotFound_ReturnError(s *UserSuite, desc string) {
	_, err := s.f.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	_, err = s.model.FindByEmail("notfound")
	s.Require().Error(err, desc)
	s.Require().Equal(domain.ErrEmailNotFound, err, desc)
}

func (s *UserSuite) TestGetInfo() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when user is found, return successfully": getInfo_FoundUser_ReturnSuccessfully,
		"when user is not found, return error":    getInfo_NotFound_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getInfo_FoundUser_ReturnSuccessfully(s *UserSuite, desc string) {
	users, err := s.f.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	expResult := domain.User{
		ID:                users[0].ID,
		Name:              users[0].Name,
		Email:             users[0].Email,
		IsSetInitCategory: users[0].IsSetInitCategory,
	}

	user, err := s.model.GetInfo(users[0].ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, user, desc)
}

func getInfo_NotFound_ReturnError(s *UserSuite, desc string) {
	_, err := s.f.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	user, err := s.model.GetInfo(999)
	s.Require().Empty(user, desc)
	s.Require().EqualError(err, domain.ErrUserIDNotFound.Error(), desc)
}

func (s *UserSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *UserSuite, desc string){
		"when is_set_init_category is set, update successfully": update_IsSetInitCategory_UpdateSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_IsSetInitCategory_UpdateSuccessfully(s *UserSuite, desc string) {
	// prepare mock data
	users, err := s.f.BuildList(2).SetZero(0, "IsSetInitCategory").SetZero(1, "IsSetInitCategory").Insert()
	s.Require().NoError(err, desc)

	// prepare update option
	t := true
	opt := domain.UpdateUserOpt{IsSetInitCategory: &t}

	// action
	err = s.model.Update(mockCtx, users[0].ID, opt)
	s.Require().NoError(err, desc)

	// check if user is updated
	checkStmt := "SELECT name, email, is_set_init_category FROM users WHERE id = ?"
	var checkedUser User
	err = s.db.QueryRow(checkStmt, users[0].ID).Scan(&checkedUser.Name, &checkedUser.Email, &checkedUser.IsSetInitCategory)
	s.Require().NoError(err, desc)
	s.Require().True(checkedUser.IsSetInitCategory, desc)
	s.Require().Equal(users[0].Name, checkedUser.Name, desc)
	s.Require().Equal(users[0].Email, checkedUser.Email, desc)

	// check if other user is not updated
	var checkedUser2 User
	err = s.db.QueryRow(checkStmt, users[1].ID).Scan(&checkedUser2.Name, &checkedUser2.Email, &checkedUser2.IsSetInitCategory)
	s.Require().NoError(err, desc)
	s.Require().False(checkedUser2.IsSetInitCategory, desc)
	s.Require().Equal(users[1].Name, checkedUser2.Name, desc)
	s.Require().Equal(users[1].Email, checkedUser2.Email, desc)
}
