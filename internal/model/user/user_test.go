package user_test

import (
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/internal/usecase"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	db      *sql.DB
	migrate *migrate.Migrate
	f       *model.Factory
	model   usecase.UserModel
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db, migrate := testutil.ConnToDB(port)
	s.model = user.NewUserModel(db)
	s.f = model.NewFactory(db)
	s.db = db
	s.migrate = migrate
}

func (s *UserSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *UserSuite) SetupTest() {
	s.model = user.NewUserModel(s.db)
	s.f = model.NewFactory(s.db)
}

func (s *UserSuite) TearDownTest() {
	if _, err := s.db.Exec("DELETE FROM users"); err != nil {
		s.Require().NoError(err)
	}
}

func (s *UserSuite) TestCreate() {
	tests := []struct {
		Desc         string
		Name         string
		Email        string
		PasswordHash string
		CheckFun     func() error
	}{
		{
			Desc:         "Create user successfully",
			Name:         "test",
			Email:        "test@gmail.com",
			PasswordHash: "test",
			CheckFun: func() error {
				stmt := `SELECT id, name, email, password_hash FROM users WHERE email = ? AND name = ?`
				var user user.User

				return s.db.QueryRow(stmt, "test@gmail.com", "test").Scan(&user.ID, &user.Name, &user.Email, &user.Password_hash)
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.Desc, func(t *testing.T) {
			err := s.model.Create(test.Name, test.Email, test.PasswordHash)
			s.Require().NoError(err, test.Desc)

			if test.CheckFun != nil {
				err = test.CheckFun()
				s.Require().NoError(err, test.Desc)
			}
		})
	}
}

func (s *UserSuite) TestFindByEmail() {
	tests := []struct {
		Desc     string
		Email    string
		SetupFun func() error
		Expected *domain.User
	}{
		{
			Desc:  "Find user successfully",
			Email: "test@gmail.com",
			SetupFun: func() error {
				overwrites := map[string]any{
					"Email": "test@gmail.com",
				}
				_, err := s.f.NewUser(overwrites)
				return err
			},
			Expected: &domain.User{
				Name:          "test",
				Email:         "test@gmail.com",
				Password_hash: "test",
			},
		},
		{
			Desc:  "User not found",
			Email: "test222@",
			SetupFun: func() error {
				overwrites := map[string]any{
					"Email": "test2222@gmail.com",
				}
				_, err := s.f.NewUser(overwrites)
				return err
			},
			Expected: nil,
		},
	}

	for _, test := range tests {
		s.T().Run(test.Desc, func(t *testing.T) {
			if test.SetupFun != nil {
				err := test.SetupFun()
				s.Require().NoError(err, test.Desc)
			}

			user, err := s.model.FindByEmail(test.Email)
			s.Require().NoError(err, test.Desc)

			if test.Expected == nil {
				s.Require().Nil(user, test.Desc)
				return
			}

			if test.Expected != nil && user == nil {
				s.Require().FailNow("Expected is not nil but user is nil", test.Desc)
			}

			s.Require().Equal(test.Expected.Name, user.Name, test.Desc)
			s.Require().Equal(test.Expected.Email, user.Email, test.Desc)
			s.Require().Equal(test.Expected.Password_hash, user.Password_hash, test.Desc)
		})
	}
}
