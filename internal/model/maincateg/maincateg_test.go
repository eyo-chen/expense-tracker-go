package maincateg_test

import (
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

type MainCategSuite struct {
	suite.Suite
	db             *sql.DB
	migrate        *migrate.Migrate
	f              *maincateg.MainCategFactory
	mainCategModel interfaces.MainCategModel
	userModel      interfaces.UserModel
	iconModel      interfaces.IconModel
}

func TestMainCategSuite(t *testing.T) {
	suite.Run(t, new(MainCategSuite))
}

func (s *MainCategSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db, migrate := testutil.ConnToDB(port)
	logger.Register()
	s.db = db
	s.migrate = migrate
	s.mainCategModel = maincateg.NewMainCategModel(db)
	s.userModel = user.NewUserModel(db)
	s.iconModel = icon.NewIconModel(db)

	s.f = maincateg.NewMainCategFactory(db)
}

func (s *MainCategSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *MainCategSuite) SetupTest() {
	s.mainCategModel = maincateg.NewMainCategModel(s.db)
	s.userModel = user.NewUserModel(s.db)
	s.iconModel = icon.NewIconModel(s.db)
	s.f = maincateg.NewMainCategFactory(s.db)
}

func (s *MainCategSuite) TearDownTest() {
	tx, err := s.db.Begin()
	s.Require().NoError(err)
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM main_categories")
	s.Require().NoError(err)

	_, err = tx.Exec("DELETE FROM users")
	s.Require().NoError(err)

	_, err = tx.Exec("DELETE FROM icons")
	s.Require().NoError(err)

	s.Require().NoError(tx.Commit())

	s.f.Reset()
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
	users, icons, err := s.f.InsertUserAndIcon(1, 1)
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: "test",
		Type: domain.Expense,
		Icon: domain.Icon{
			ID: icons[0].ID,
		},
	}
	err = s.mainCategModel.Create(categ, users[0].ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE user_id = ?
							 AND name = ?
							 AND type = ?
							 `
	var result maincateg.MainCateg
	err = s.db.QueryRow(checkStmt, users[0].ID, "test", domain.Expense.ToModelValue()).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(categ.Name, result.Name, desc)
	s.Require().Equal(categ.Type.ToModelValue(), result.Type, desc)
	s.Require().Equal(icons[0].ID, result.IconID, desc)
}

func create_DuplicateName_ReturnError(s *MainCategSuite, desc string) {
	createdMainCateg, user, _, err := s.f.InsertMainCategWithAss(maincateg.MainCateg{})
	s.Require().NoError(err, desc)

	icon, err := s.f.Icon.Build().Insert()
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: createdMainCateg.Name,
		Type: domain.Income,
		Icon: domain.Icon{
			ID: icon.ID,
		},
	}
	err = s.mainCategModel.Create(categ, user.ID)
	s.Require().EqualError(err, domain.ErrUniqueNameUserType.Error(), desc)
}

func create_DuplicateIcon_ReturnError(s *MainCategSuite, desc string) {
	createdMainCateg, user, icon, err := s.f.InsertMainCategWithAss(maincateg.MainCateg{})
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: createdMainCateg.Name + "1", // different name
		Type: domain.Expense,
		Icon: domain.Icon{
			ID: icon.ID,
		},
	}
	err = s.mainCategModel.Create(categ, user.ID)
	s.Require().EqualError(err, domain.ErrUniqueIconUser.Error(), desc)
}

func (s *MainCategSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when specify income type, return only income type data":   getAll_IncomeType_ReturnOnlyIncomeTypeData,
		"when specify expense type, return only expense type data": getAll_ExpenseType_ReturnOnlyExpenseTypeData,
		"when specify unspecified type, return all data":           getAll_UnSpecifiedType_ReturnAllData,
		"when multiple users, return correct data":                 getAll_MultipleUsers_ReturnCorrectData,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAll_IncomeType_ReturnOnlyIncomeTypeData(s *MainCategSuite, desc string) {
	mainCategList, users, icons, err := s.f.InsertMainCategListWithAss(2, 1, 2, "expense", "income")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:   mainCategList[1].ID,
			Name: mainCategList[1].Name,
			Type: domain.Income,
			Icon: domain.Icon{
				ID:  icons[1].ID,
				URL: icons[1].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(users[0].ID, domain.Income)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, categs, desc)
}

func getAll_ExpenseType_ReturnOnlyExpenseTypeData(s *MainCategSuite, desc string) {
	mainCategList, users, icons, err := s.f.InsertMainCategListWithAss(2, 1, 2, "expense", "income")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:   mainCategList[0].ID,
			Name: mainCategList[0].Name,
			Type: domain.Expense,
			Icon: domain.Icon{
				ID:  icons[0].ID,
				URL: icons[0].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(users[0].ID, domain.Expense)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, categs, desc)
}

func getAll_UnSpecifiedType_ReturnAllData(s *MainCategSuite, desc string) {
	mainCategList, users, icons, err := s.f.InsertMainCategListWithAss(2, 1, 2, "expense", "income")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:   mainCategList[0].ID,
			Name: mainCategList[0].Name,
			Type: domain.Expense,
			Icon: domain.Icon{
				ID:  icons[0].ID,
				URL: icons[0].URL,
			},
		},
		{
			ID:   mainCategList[1].ID,
			Name: mainCategList[1].Name,
			Type: domain.Income,
			Icon: domain.Icon{
				ID:  icons[1].ID,
				URL: icons[1].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(users[0].ID, domain.UnSpecified)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, categs, desc)
}

func getAll_MultipleUsers_ReturnCorrectData(s *MainCategSuite, desc string) {
	mainCategList, users, icons, err := s.f.InsertMainCategListWithAss(3, 2, 3, "expense", "income", "expense")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:   mainCategList[1].ID,
			Name: mainCategList[1].Name,
			Type: domain.Income,
			Icon: domain.Icon{
				ID:  icons[1].ID,
				URL: icons[1].URL,
			},
		},
		{
			ID:   mainCategList[2].ID,
			Name: mainCategList[2].Name,
			Type: domain.Expense,
			Icon: domain.Icon{
				ID:  icons[2].ID,
				URL: icons[2].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(users[1].ID, domain.UnSpecified)
	s.Require().NoError(err)
	s.Require().Equal(expResult, categs)
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
	categ, _, _, err := s.f.InsertMainCategWithAss(maincateg.MainCateg{})
	s.Require().NoError(err, desc)

	inputCateg := &domain.MainCateg{
		ID:   categ.ID,
		Name: "test2",
		Type: domain.Income,
		Icon: domain.Icon{
			ID: categ.IconID,
		},
	}
	err = s.mainCategModel.Update(inputCateg)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE id = ?
							 `
	var result maincateg.MainCateg
	err = s.db.QueryRow(checkStmt, categ.ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputCateg.Name, result.Name, desc)
	s.Require().Equal(inputCateg.Type.ToModelValue(), result.Type, desc)
}

func update_WithMultipleUser_UpdateSuccessfully(s *MainCategSuite, desc string) {
	categs, _, _, err := s.f.InsertMainCategListWithAss(2, 2, 1)
	s.Require().NoError(err, desc)

	inputCateg := &domain.MainCateg{
		ID:   categs[0].ID,
		Name: "update name",
		Type: domain.Income,
		Icon: domain.Icon{
			ID: categs[0].IconID,
		},
	}
	err = s.mainCategModel.Update(inputCateg)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE id = ?
							 `
	// check if the data is updated
	var result maincateg.MainCateg
	err = s.db.QueryRow(checkStmt, categs[0].ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputCateg.Name, result.Name, desc)
	s.Require().Equal(inputCateg.Type.ToModelValue(), result.Type, desc)

	// check if the other data is not updated
	var result2 maincateg.MainCateg
	err = s.db.QueryRow(checkStmt, categs[1].ID).Scan(&result2.ID, &result2.Name, &result2.Type, &result2.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(categs[1].Name, result2.Name, desc)
	s.Require().Equal(categs[1].Type, result2.Type, desc)
}

func update_DuplicateName_ReturnError(s *MainCategSuite, desc string) {
	categs, _, _, err := s.f.InsertMainCategListWithAss(2, 1, 2)
	s.Require().NoError(err, desc)

	domainMainCateg := &domain.MainCateg{
		ID:   categs[0].ID,
		Name: categs[1].Name, // update categ1 with categ2 name
		Type: domain.CvtToTransactionType(categs[0].Type),
		Icon: domain.Icon{
			ID: categs[0].IconID,
		},
	}
	err = s.mainCategModel.Update(domainMainCateg)
	s.Require().EqualError(err, domain.ErrUniqueNameUserType.Error(), desc)
}

func update_DuplicateIcon_ReturnError(s *MainCategSuite, desc string) {
	categs, _, _, err := s.f.InsertMainCategListWithAss(2, 1, 2)
	s.Require().NoError(err, desc)

	domainMainCateg := &domain.MainCateg{
		ID:   categs[0].ID,
		Name: categs[0].Name + "2", // make sure the name is different
		Type: domain.CvtToTransactionType(categs[0].Type),
		Icon: domain.Icon{
			ID: categs[1].IconID, // update categ1 with categ2 icon
		},
	}
	err = s.mainCategModel.Update(domainMainCateg)
	s.Require().EqualError(err, domain.ErrUniqueIconUser.Error(), desc)
}
