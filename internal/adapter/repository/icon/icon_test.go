package icon

import (
	"context"
	"database/sql"
	"errors"
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

type IconSuite struct {
	suite.Suite
	dk      *dockerutil.Container
	db      *sql.DB
	migrate *migrate.Migrate
	repo    *Repo
	f       *factory
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupSuite() {
	s.dk = dockerutil.RunDocker(dockerutil.ImageMySQL)
	db, migrate := testutil.ConnToDB(s.dk.Port)
	logger.Register()
	s.repo = New(db)
	s.db = db
	s.migrate = migrate
	s.f = newFactory(db)
}

func (s *IconSuite) TearDownSuite() {
	if err := s.db.Close(); err != nil {
		logger.Error("Unable to close mysql database", "error", err)
	}
	s.migrate.Close()
	s.dk.PurgeDocker()
}

func (s *IconSuite) SetupTest() {
	s.repo = New(s.db)
	s.f = newFactory(s.db)
}

func (s *IconSuite) TearDownTest() {
	tx, err := s.db.Begin()
	s.Require().NoError(err)
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.Require().NoError(err)
		}
	}()

	_, err = tx.Exec("DELETE FROM icons")
	s.Require().NoError(err)

	s.Require().NoError(tx.Commit())
	s.f.Reset()
}

func (s *IconSuite) TestGetByID() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when has icon, return icon":   getByID_WithIcon_ReturnIcon,
		"when has no icon, return err": getByID_WithoutIcon_ReturnErr,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByID_WithIcon_ReturnIcon(s *IconSuite, desc string) {
	icons, err := s.f.InsertMany(mockCTX, 2)
	s.Require().NoError(err, desc)

	expRes := domain.DefaultIcon{
		ID:  icons[0].ID,
		URL: icons[0].URL,
	}

	res, err := s.repo.GetByID(mockCTX, icons[0].ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, res, desc)
}

func getByID_WithoutIcon_ReturnErr(s *IconSuite, desc string) {
	_, err := s.f.InsertMany(mockCTX, 2)
	s.Require().NoError(err, desc)

	res, err := s.repo.GetByID(mockCTX, 999)
	s.Require().Empty(res, desc)
	s.Require().ErrorIs(err, domain.ErrIconNotFound, desc)
}

func (s *IconSuite) TestList() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when has icons, return all":    list_WithIcons_ReturnAll,
		"when has no icons, return nil": list_WithoutIcons_ReturnNil,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func list_WithIcons_ReturnAll(s *IconSuite, desc string) {
	icons, err := s.f.InsertMany(mockCTX, 2)
	s.Require().NoError(err, desc)

	expRes := []domain.DefaultIcon{
		{
			ID:  icons[0].ID,
			URL: icons[0].URL,
		},
		{
			ID:  icons[1].ID,
			URL: icons[1].URL,
		},
	}

	res, err := s.repo.List()
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, res, desc)
}

func list_WithoutIcons_ReturnNil(s *IconSuite, desc string) {
	res, err := s.repo.List()
	s.Require().NoError(err, desc)
	s.Require().Nil(res, desc)
}

func (s *IconSuite) TestGetByIDs() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when has icon, return icon":   getByIDs_WithIcon_ReturnIcons,
		"when has no icon, return err": getByIDs_WithoutIcon_ReturnErr,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByIDs_WithIcon_ReturnIcons(s *IconSuite, desc string) {
	icons, err := s.f.InsertMany(mockCTX, 5)
	s.Require().NoError(err, desc)

	expRes := map[int64]domain.DefaultIcon{
		icons[0].ID: {
			ID:  icons[0].ID,
			URL: icons[0].URL,
		},
		icons[1].ID: {
			ID:  icons[1].ID,
			URL: icons[1].URL,
		},
	}

	ids := []int64{icons[0].ID, icons[1].ID, 999}
	res, err := s.repo.GetByIDs(ids)
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, res, desc)
}

func getByIDs_WithoutIcon_ReturnErr(s *IconSuite, desc string) {
	_, err := s.f.InsertMany(mockCTX, 5)
	s.Require().NoError(err, desc)

	ids := []int64{999}
	res, err := s.repo.GetByIDs(ids)
	s.Require().Nil(res, desc)
	s.Require().ErrorIs(err, domain.ErrIconNotFound, desc)
}
