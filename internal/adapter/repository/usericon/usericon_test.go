package usericon

import (
	"context"
	"database/sql"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/dockerutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX = context.Background()
)

type UserIconSuite struct {
	suite.Suite
	dk      *dockerutil.Container
	db      *sql.DB
	migrate *migrate.Migrate
	repo    *Repo
	factory *factory
}

func TestUserIconSuite(t *testing.T) {
	suite.Run(t, new(UserIconSuite))
}

func (s *UserIconSuite) SetupSuite() {
	s.dk = dockerutil.RunDocker(dockerutil.ImageMySQL)
	db, migrate := testutil.ConnToDB(s.dk.Port)
	s.repo = New(db)
	s.factory = newFactory(db)
	logger.Register()
	s.db = db
	s.migrate = migrate
}

func (s *UserIconSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	s.dk.PurgeDocker()
}

func (s *UserIconSuite) SetupTest() {
	s.repo = New(s.db)
}

func (s *UserIconSuite) TearDownTest() {
	if _, err := s.db.Exec("DELETE FROM user_icons"); err != nil {
		s.Require().NoError(err)
	}
}

func (s *UserIconSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *UserIconSuite, desc string){
		"when no error, return successfully": create_NoError_ReturnSuccessfully,
		"when user not found, return error":  create_UserNotFound_ReturnErr,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_ReturnSuccessfully(s *UserIconSuite, desc string) {
	// prepare mock data
	user, err := s.factory.InsertUser(mockCTX)
	s.Require().NoError(err)
	userIcon := domain.UserIcon{
		UserID:    user.ID,
		ObjectKey: "test",
	}

	// action
	err = s.repo.Create(mockCTX, userIcon)

	// assertion
	s.Require().NoError(err)

	// check data
	checkStmt := `SELECT user_id, object_key FROM user_icons WHERE user_id = ?`
	rows, err := s.db.Query(checkStmt, user.ID)
	s.Require().NoError(err)
	defer rows.Close()

	var userID int64
	var objectKey string
	for rows.Next() {
		if err := rows.Scan(&userID, &objectKey); err != nil {
			s.Require().NoError(err)
		}
	}

	s.Require().Equal(userIcon.UserID, userID)
	s.Require().Equal(userIcon.ObjectKey, objectKey)
}

func create_UserNotFound_ReturnErr(s *UserIconSuite, desc string) {
	// prepare mock data
	userIcon := domain.UserIcon{
		UserID:    111, // user not found
		ObjectKey: "test",
	}

	// action
	err := s.repo.Create(mockCTX, userIcon)

	// assertion
	s.Require().ErrorIs(err, domain.ErrUserNotFound)
}

func (s *UserIconSuite) TestGetByUserID() {
	for scenario, fn := range map[string]func(s *UserIconSuite, desc string){
		"when no error, return successfully": getByUserID_NoError_ReturnSuccessfully,
		"when no data, return empty":         getByUserID_NoData_ReturnEmpty,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByUserID_NoError_ReturnSuccessfully(s *UserIconSuite, desc string) {
	// prepare mock data
	mockUserIcons, user, err := s.factory.InsertManyWithOneUser(mockCTX, 2)
	s.Require().NoError(err)
	_, _, err = s.factory.InsertManyWithOneUser(mockCTX, 2) // prepare more mock data
	s.Require().NoError(err)

	// prepare expected data
	expResp := []domain.UserIcon{
		{ID: mockUserIcons[0].ID, UserID: user.ID, ObjectKey: mockUserIcons[0].ObjectKey},
		{ID: mockUserIcons[1].ID, UserID: user.ID, ObjectKey: mockUserIcons[1].ObjectKey},
	}

	// action
	userIcons, err := s.repo.GetByUserID(mockCTX, user.ID)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, userIcons, desc)
}

func getByUserID_NoData_ReturnEmpty(s *UserIconSuite, desc string) {
	// prepare mock data
	_, _, err := s.factory.InsertManyWithOneUser(mockCTX, 2)
	s.Require().NoError(err)

	// action
	userIcons, err := s.repo.GetByUserID(mockCTX, 9999)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Empty(userIcons, desc)
}
