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

type MainCategSuite struct {
	suite.Suite
	db             *sql.DB
	f              *factory
	mainCategModel usecase.MainCategModel
	userModel      usecase.UserModel
	iconModel      usecase.IconModel
}

func TestMainCategSuite(t *testing.T) {
	suite.Run(t, new(MainCategSuite))
}

func (s *MainCategSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db := testutil.ConnToDB(port)
	logger.Register()
	s.db = db
	s.mainCategModel = newMainCategModel(db)
	s.userModel = newUserModel(db)
	s.iconModel = newIconModel(db)
	s.f = newFactory(db)
}

func (s *MainCategSuite) TearDownSuite() {
	s.db.Close()
	dockerutil.PurgeDocker()
}

func (s *MainCategSuite) SetupTest() {
	s.mainCategModel = newMainCategModel(s.db)
	s.userModel = newUserModel(s.db)
	s.iconModel = newIconModel(s.db)
	s.f = newFactory(s.db)
}

func (s *MainCategSuite) TearDownTest() {
	tx, err := s.db.Begin()
	if err != nil {
		s.Require().NoError(err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM main_categories"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM users"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM icons"); err != nil {
		s.Require().NoError(err)
	}

	s.Require().NoError(tx.Commit())
}

func (s *MainCategSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no duplicate data, create successfully": create_NoDuplicate_CreateSuccessfully,
		"when duplicate name, return error":           create_DuplicateName_ReturnError,
		"when duplicate icon, return error":           create_DuplicateIcon_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoDuplicate_CreateSuccessfully(s *MainCategSuite, desc string) {
	user, err := s.f.newUser()
	s.Require().NoError(err, desc)

	icon, err := s.f.newIcon()
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: "test",
		Type: domain.Expense,
		Icon: &domain.Icon{
			ID: icon.ID,
		},
	}
	err = s.mainCategModel.Create(categ, user.ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE user_id = ?
							 AND name = ?
							 AND type = ?
							 `
	var result MainCateg
	err = s.db.QueryRow(checkStmt, user.ID, "test", domain.Expense.ModelValue()).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(categ.Name, result.Name, desc)
	s.Require().Equal(categ.Type.ModelValue(), result.Type, desc)
	s.Require().Equal(icon.ID, result.IconID, desc)
}

func create_DuplicateName_ReturnError(s *MainCategSuite, desc string) {
	user, err := s.f.newUser()
	s.Require().NoError(err, desc)

	_, err = s.f.newIcon()
	s.Require().NoError(err, desc)

	icon1, err := s.f.newIcon()
	s.Require().NoError(err, desc)

	overwrite := map[string]any{
		"Type": domain.Expense.ModelValue(),
	}
	createdMainCateg, err := s.f.newMainCateg(user, overwrite)
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: createdMainCateg.Name,
		Type: domain.Expense,
		Icon: &domain.Icon{
			ID: icon1.ID,
		},
	}
	err = s.mainCategModel.Create(categ, user.ID)
	s.Require().EqualError(err, domain.ErrUniqueNameUserType.Error(), desc)
}

func create_DuplicateIcon_ReturnError(s *MainCategSuite, desc string) {
	user, err := s.f.newUser()
	s.Require().NoError(err, desc)

	icon, err := s.f.newIcon()
	s.Require().NoError(err, desc)

	overwrite := map[string]any{
		"IconID": icon.ID,
	}
	createdMainCateg, err := s.f.newMainCateg(user, overwrite)
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: createdMainCateg.Name + "1", // different name
		Type: domain.Expense,
		Icon: &domain.Icon{
			ID: createdMainCateg.IconID,
		},
	}
	err = s.mainCategModel.Create(categ, user.ID)
	s.Require().EqualError(err, domain.ErrUniqueIconUser.Error(), desc)
}

func (s *MainCategSuite) TestGetAll() {
	overwrite := map[string]any{
		"Email": "test1@gmail.com",
	}
	user1, err := s.f.newUser(overwrite)
	s.Require().NoError(err)

	categ1, err := s.f.newMainCateg(user1)
	s.Require().NoError(err)
	_, err = s.f.newMainCateg(nil)
	s.Require().NoError(err)

	categs, err := s.mainCategModel.GetAll(user1.ID)
	s.Require().NoError(err)

	s.Require().Equal(1, len(categs))
	s.Require().Equal(categ1.Name, categs[0].Name)
	s.Require().Equal(categ1.Type, categs[0].Type.ModelValue())
}

func (s *MainCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no duplicate data, update successfully":  update_NoDuplicate_UpdateSuccessfully,
		"when with multiple user, update successfully": update_WithMultipleUser_UpdateSuccessfully,
		"when duplicate name, return error":            update_DuplicateName_ReturnError,
		"when duplicate icon, return error":            update_DuplicateIcon_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoDuplicate_UpdateSuccessfully(s *MainCategSuite, desc string) {
	// prepare existing data
	mainCateg, err := s.f.newMainCateg(nil)
	s.Require().NoError(err, desc)

	// prepare updating data with different name and type
	domainMainCateg := &domain.MainCateg{
		ID:   mainCateg.ID,
		Name: "test2",
		Type: domain.Income,
		Icon: &domain.Icon{
			ID: mainCateg.IconID,
		},
	}
	err = s.mainCategModel.Update(domainMainCateg)
	s.Require().NoError(err, desc)

	// check if the data is updated
	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE id = ?
							 `
	var result MainCateg
	err = s.db.QueryRow(checkStmt, mainCateg.ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(domainMainCateg.Name, result.Name, desc)
	s.Require().Equal(domainMainCateg.Type.ModelValue(), result.Type, desc)
}

func update_WithMultipleUser_UpdateSuccessfully(s *MainCategSuite, desc string) {
	// prepare two users
	overwrite := map[string]any{"Email": "test@gmail.com"}
	user, err := s.f.newUser(overwrite)
	s.Require().NoError(err, desc)
	overwrite = map[string]any{"Email": "test1@gmail.com"}
	user1, err := s.f.newUser(overwrite)
	s.Require().NoError(err, desc)

	// prepare two existing datas for each user
	createdMainCateg, err := s.f.newMainCateg(user)
	s.Require().NoError(err, desc)
	createdMainCateg1, err := s.f.newMainCateg(user1)
	s.Require().NoError(err, desc)

	// prepare updating data with different name and type
	domainMainCateg := &domain.MainCateg{
		ID:   createdMainCateg.ID,
		Name: "update name",
		Type: domain.Income,
		Icon: &domain.Icon{
			ID: createdMainCateg.IconID,
		},
	}
	err = s.mainCategModel.Update(domainMainCateg)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE id = ?
							 `
	// check if the data is updated
	var result MainCateg
	err = s.db.QueryRow(checkStmt, createdMainCateg.ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(domainMainCateg.Name, result.Name, desc)
	s.Require().Equal(domainMainCateg.Type.ModelValue(), result.Type, desc)

	// check if the data of other user is not updated
	var result2 MainCateg
	err = s.db.QueryRow(checkStmt, createdMainCateg1.ID).Scan(&result2.ID, &result2.Name, &result2.Type, &result2.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(createdMainCateg1.Name, result2.Name, desc)
	s.Require().Equal(createdMainCateg1.Type, result2.Type, desc)
}

func update_DuplicateName_ReturnError(s *MainCategSuite, desc string) {
	// prepare user
	user, err := s.f.newUser()
	s.Require().NoError(err, desc)

	// prepare existing data
	overwrite := map[string]any{
		"Type": domain.Expense.ModelValue(),
	}
	createdMainCateg, err := s.f.newMainCateg(user, overwrite)
	s.Require().NoError(err, desc)

	// prepare existing data for with different name
	overwrite = map[string]any{
		"Name": "test1",
		"Type": domain.Expense.ModelValue(),
	}
	createdMainCateg1, err := s.f.newMainCateg(user, overwrite)
	s.Require().NoError(err, desc)

	// prepare updating data with duplicate name
	domainMainCateg := &domain.MainCateg{
		ID:   createdMainCateg.ID,
		Name: createdMainCateg1.Name,
		Type: domain.Expense,
		Icon: &domain.Icon{
			ID: createdMainCateg.IconID,
		},
	}
	err = s.mainCategModel.Update(domainMainCateg)
	s.Require().EqualError(err, domain.ErrUniqueNameUserType.Error(), desc)
}

func update_DuplicateIcon_ReturnError(s *MainCategSuite, desc string) {
	// prepare user
	user, err := s.f.newUser()
	s.Require().NoError(err, desc)

	// prepare existing data
	icon, err := s.f.newIcon()
	s.Require().NoError(err, desc)
	overwrite := map[string]any{
		"IconID": icon.ID,
	}
	createdMainCateg, err := s.f.newMainCateg(user, overwrite)
	s.Require().NoError(err, desc)

	// prepare existing data for with different icon and name
	icon1, err := s.f.newIcon()
	s.Require().NoError(err, desc)
	overwrite = map[string]any{
		"Name":   "test1",
		"IconID": icon1.ID,
	}
	createdMainCateg1, err := s.f.newMainCateg(user, overwrite)
	s.Require().NoError(err, desc)

	// prepare updating data with duplicate icon
	domainMainCateg := &domain.MainCateg{
		ID:   createdMainCateg.ID,
		Name: createdMainCateg.Name + "2",
		Type: domain.Expense,
		Icon: &domain.Icon{
			ID: createdMainCateg1.IconID,
		},
	}

	err = s.mainCategModel.Update(domainMainCateg)
	s.Require().EqualError(err, domain.ErrUniqueIconUser.Error(), desc)
}
