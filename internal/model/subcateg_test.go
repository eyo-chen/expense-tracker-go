package model

import (
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type SubCategSuite struct {
	suite.Suite
	db            *sql.DB
	f             *factory
	subCategModel usecase.SubCategModel
}

func TestSubCategSuite(t *testing.T) {
	suite.Run(t, new(SubCategSuite))
}

func (s *SubCategSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db := testutil.ConnToDB(port)
	logger.Register()
	s.db = db
	s.subCategModel = newSubCategModel(db)
	s.f = newFactory(db)
}

func (s *SubCategSuite) TearDownSuite() {
	s.db.Close()
	dockerutil.PurgeDocker()
}

func (s *SubCategSuite) SetupTest() {
	s.subCategModel = newSubCategModel(s.db)
	s.f = newFactory(s.db)
}

func (s *SubCategSuite) TearDownTest() {
	tx, err := s.db.Begin()
	if err != nil {
		s.Require().NoError(err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM sub_categories"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM main_categories"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM icons"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM users"); err != nil {
		s.Require().NoError(err)
	}

	s.Require().NoError(tx.Commit())
}

func (s *SubCategSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when no duplicate data, create successfully": create_NoDuplicateData_CreateSuccessfully,
		"when there is duplicate name, return error":  create_DuplicateNameUserMainCateg_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoDuplicateData_CreateSuccessfully(s *SubCategSuite, desc string) {
	// prepare data
	user, err := s.f.newUser()
	s.Require().NoError(err, desc)
	mainCateg, err := s.f.newMainCateg(user)
	s.Require().NoError(err, desc)

	// prepare input data
	subCateg := &domain.SubCateg{
		Name:        "test",
		MainCategID: mainCateg.ID,
	}

	// action
	err = s.subCategModel.Create(subCateg, user.ID)
	s.Require().NoError(err, desc)

	// check
	var result domain.SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user.ID, mainCateg.ID, subCateg.Name).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(subCateg.Name, result.Name, desc)
	s.Require().Equal(subCateg.MainCategID, result.MainCategID, desc)
}

func create_DuplicateNameUserMainCateg_ReturnError(s *SubCategSuite, desc string) {
	// prepare data
	user, err := s.f.newUser()
	s.Require().NoError(err, desc)
	mainCateg, err := s.f.newMainCateg(user)
	s.Require().NoError(err, desc)

	// prepare existing data
	subCateg, err := s.f.newSubCateg(user, mainCateg)
	s.Require().NoError(err, desc)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		Name:        subCateg.Name,
		MainCategID: mainCateg.ID,
	}

	// action and check
	err = s.subCategModel.Create(inputSubCateg, user.ID)
	s.Require().Equal(domain.ErrUniqueNameUserMainCateg, err, desc)
}
