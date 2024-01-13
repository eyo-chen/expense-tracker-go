package model

import (
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	db    *sql.DB
	model usecase.UserModel
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db := testutil.ConnToDB(port)
	s.model = newUserModel(db)
	s.db = db
}

func (s *UserSuite) TearDownSuite() {
	s.db.Close()
	dockerutil.PurgeDocker()
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
			Email:        "TestCreate@gmail.co",
			PasswordHash: "test",
			CheckFun: func() error {
				stmt := `SELECT id, name, email, password_hash FROM users WHERE email = ? AND name = ?`
				var user User

				return s.db.QueryRow(stmt, "TestCreate@gmail.com", "test").Scan(&user.ID, &user.Name, &user.Email, &user.Password_hash)
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.Desc, func(t *testing.T) {
			err := s.model.Create(test.Name, test.Email, test.PasswordHash)
			s.NoError(err, test.Desc)

			if test.CheckFun != nil {
				err = test.CheckFun()
				s.NoError(err, test.Desc)
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
			Email: "TestFindByEmail@gmail.com",
			SetupFun: func() error {
				stmt := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`
				_, err := s.db.Exec(stmt, "test", "TestFindByEmail@gmail.com", "test")

				return err
			},
			Expected: &domain.User{
				ID:            1,
				Name:          "test",
				Email:         "TestFindByEmail@gmail.com",
				Password_hash: "test",
			},
		},
		{
			Desc:     "User not found",
			Email:    "",
			SetupFun: nil,
			Expected: nil,
		},
	}

	for _, test := range tests {
		s.T().Run(test.Desc, func(t *testing.T) {
			if test.SetupFun != nil {
				err := test.SetupFun()
				s.NoError(err, test.Desc)
			}

			user, err := s.model.FindByEmail(test.Email)
			s.NoError(err, test.Desc)

			if test.Expected == nil {
				s.Nil(user, test.Desc)
				return
			}

			if test.Expected != nil && user == nil {
				s.FailNow("Expected is not nil but user is nil", test.Desc)
			}

			s.Equal(test.Expected.Name, user.Name, test.Desc)
			s.Equal(test.Expected.Email, user.Email, test.Desc)
			s.Equal(test.Expected.Password_hash, user.Password_hash, test.Desc)
		})
	}
}
