package subcateg_test

import (
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/usecase"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type SubCategSuite struct {
	suite.Suite
	db            *sql.DB
	f             *model.Factory
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
	s.subCategModel = subcateg.NewSubCategModel(db)
	s.f = model.NewFactory(db)
}

func (s *SubCategSuite) TearDownSuite() {
	s.db.Close()
	dockerutil.PurgeDocker()
}

func (s *SubCategSuite) SetupTest() {
	s.subCategModel = subcateg.NewSubCategModel(s.db)
	s.f = model.NewFactory(s.db)
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
	user, err := s.f.NewUser()
	s.Require().NoError(err, desc)
	mainCateg, err := s.f.NewMainCateg(user)
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
	var result subcateg.SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user.ID, mainCateg.ID, subCateg.Name).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(subCateg.Name, result.Name, desc)
	s.Require().Equal(subCateg.MainCategID, result.MainCategID, desc)
}

func create_DuplicateNameUserMainCateg_ReturnError(s *SubCategSuite, desc string) {
	// prepare data
	user, err := s.f.NewUser()
	s.Require().NoError(err, desc)
	mainCateg, err := s.f.NewMainCateg(user)
	s.Require().NoError(err, desc)

	// prepare existing data
	subCateg, err := s.f.NewSubCateg(user, mainCateg)
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

func (s *SubCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when no duplicate data, update successfully":                  updatesub_NoDuplicateData_UpdateSuccessfully,
		"when there are multiple main categories, update successfully": updatesub_WithMultipleMainCateg_UpdateSuccessfully,
		"when update to duplicate name, return error":                  updatesub_DuplicateName_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func updatesub_NoDuplicateData_UpdateSuccessfully(s *SubCategSuite, desc string) {
	// prepare data
	user, err := s.f.NewUser()
	s.Require().NoError(err, desc)
	mainCateg, err := s.f.NewMainCateg(user)
	s.Require().NoError(err, desc)
	subCateg, err := s.f.NewSubCateg(user, mainCateg)
	s.Require().NoError(err, desc)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		ID:          subCateg.ID,
		Name:        "updated test",
		MainCategID: mainCateg.ID,
	}

	// action
	err = s.subCategModel.Update(inputSubCateg)
	s.Require().NoError(err, desc)

	// check
	var result subcateg.SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user.ID, mainCateg.ID, inputSubCateg.Name).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputSubCateg.Name, result.Name, desc)
	s.Require().Equal(inputSubCateg.MainCategID, result.MainCategID, desc)
}

func updatesub_WithMultipleMainCateg_UpdateSuccessfully(s *SubCategSuite, desc string) {
	// prepare existing data
	// user
	user1, err := s.f.NewUser()
	s.Require().NoError(err, desc)

	// main category1
	overwrite := map[string]interface{}{"Name": "test1"}
	mainCateg, err := s.f.NewMainCateg(user1, overwrite)
	s.Require().NoError(err, desc)

	// main category2
	overwrite = map[string]interface{}{"Name": "test2"}
	_, err = s.f.NewMainCateg(user1, overwrite)
	s.Require().NoError(err, desc)

	// sub category
	subCateg, err := s.f.NewSubCateg(user1, mainCateg)
	s.Require().NoError(err, desc)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		ID:          subCateg.ID,
		Name:        "updated test",
		MainCategID: mainCateg.ID,
	}

	// action
	err = s.subCategModel.Update(inputSubCateg)
	s.Require().NoError(err, desc)

	// check
	var result subcateg.SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user1.ID, mainCateg.ID, inputSubCateg.Name).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputSubCateg.Name, result.Name, desc)
	s.Require().Equal(inputSubCateg.MainCategID, result.MainCategID, desc)
}

func updatesub_DuplicateName_ReturnError(s *SubCategSuite, desc string) {
	// prepare existing data
	user, err := s.f.NewUser()
	s.Require().NoError(err, desc)
	mainCateg, err := s.f.NewMainCateg(user)
	s.Require().NoError(err, desc)

	overwrite := map[string]interface{}{"Name": "test1"}
	subCateg1, err := s.f.NewSubCateg(user, mainCateg, overwrite)
	s.Require().NoError(err, desc)

	overwrite = map[string]interface{}{"Name": "test2"}
	subCateg2, err := s.f.NewSubCateg(user, mainCateg, overwrite)
	s.Require().NoError(err, desc)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		ID:          subCateg1.ID,
		Name:        subCateg2.Name,
		MainCategID: mainCateg.ID,
	}

	// action and check
	err = s.subCategModel.Update(inputSubCateg)
	s.Require().Equal(domain.ErrUniqueNameUserMainCateg, err, desc)
}
