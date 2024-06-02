package maincateg

import (
	"context"
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/interfaces"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx = context.Background()
)

type MainCategSuite struct {
	suite.Suite
	db             *sql.DB
	migrate        *migrate.Migrate
	f              *factory
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
	s.mainCategModel = NewMainCategModel(db)
	s.userModel = user.NewUserModel(db)
	s.iconModel = icon.NewIconModel(db)

	s.f = newFactory(db)
}

func (s *MainCategSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *MainCategSuite) SetupTest() {
	s.mainCategModel = NewMainCategModel(s.db)
	s.userModel = user.NewUserModel(s.db)
	s.iconModel = icon.NewIconModel(s.db)
	s.f = newFactory(s.db)
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
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoDuplicate_CreateSuccessfully(s *MainCategSuite, desc string) {
	users, icons, err := s.f.InsertUsersAndIcons(1, 1)
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: "test",
		Type: domain.TransactionTypeExpense,
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
	var result MainCateg
	err = s.db.QueryRow(checkStmt, users[0].ID, "test", domain.TransactionTypeExpense.ToModelValue()).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(categ.Name, result.Name, desc)
	s.Require().Equal(categ.Type.ToModelValue(), result.Type, desc)
	s.Require().Equal(icons[0].ID, result.IconID, desc)
}

func create_DuplicateName_ReturnError(s *MainCategSuite, desc string) {
	createdMainCateg, user, _, err := s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	icon, err := s.f.Icon.Build().Insert()
	s.Require().NoError(err, desc)

	categ := &domain.MainCateg{
		Name: createdMainCateg.Name,
		Type: domain.TransactionTypeIncome,
		Icon: domain.Icon{
			ID: icon.ID,
		},
	}
	err = s.mainCategModel.Create(categ, user.ID)
	s.Require().EqualError(err, domain.ErrUniqueNameUserType.Error(), desc)
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
			Type: domain.TransactionTypeIncome,
			Icon: domain.Icon{
				ID:  icons[1].ID,
				URL: icons[1].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(mockCtx, users[0].ID, domain.TransactionTypeIncome)
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
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{
				ID:  icons[0].ID,
				URL: icons[0].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(mockCtx, users[0].ID, domain.TransactionTypeExpense)
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
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{
				ID:  icons[0].ID,
				URL: icons[0].URL,
			},
		},
		{
			ID:   mainCategList[1].ID,
			Name: mainCategList[1].Name,
			Type: domain.TransactionTypeIncome,
			Icon: domain.Icon{
				ID:  icons[1].ID,
				URL: icons[1].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(mockCtx, users[0].ID, domain.TransactionTypeUnSpecified)
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
			Type: domain.TransactionTypeIncome,
			Icon: domain.Icon{
				ID:  icons[1].ID,
				URL: icons[1].URL,
			},
		},
		{
			ID:   mainCategList[2].ID,
			Name: mainCategList[2].Name,
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{
				ID:  icons[2].ID,
				URL: icons[2].URL,
			},
		},
	}

	categs, err := s.mainCategModel.GetAll(mockCtx, users[1].ID, domain.TransactionTypeUnSpecified)
	s.Require().NoError(err)
	s.Require().Equal(expResult, categs)
}

func (s *MainCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no duplicate data, update successfully":  update_NoDuplicate_UpdateSuccessfully,
		"when with multiple user, update successfully": update_WithMultipleUser_UpdateSuccessfully,
		"when duplicate name, return error":            update_DuplicateName_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoDuplicate_UpdateSuccessfully(s *MainCategSuite, desc string) {
	categ, _, _, err := s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	inputCateg := &domain.MainCateg{
		ID:   categ.ID,
		Name: "test2",
		Type: domain.TransactionTypeIncome,
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
	var result MainCateg
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
		Type: domain.TransactionTypeIncome,
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
	var result MainCateg
	err = s.db.QueryRow(checkStmt, categs[0].ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputCateg.Name, result.Name, desc)
	s.Require().Equal(inputCateg.Type.ToModelValue(), result.Type, desc)

	// check if the other data is not updated
	var result2 MainCateg
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

func (s *MainCategSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, delete successfully": delete_NoError_DeleteSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func delete_NoError_DeleteSuccessfully(s *MainCategSuite, desc string) {
	categ, _, _, err := s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	err = s.mainCategModel.Delete(categ.ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id
							 FROM main_categories
							 WHERE id = ?
							 `
	err = s.db.QueryRow(checkStmt, categ.ID).Scan(&categ.ID)
	s.Require().EqualError(err, sql.ErrNoRows.Error(), desc)
}

func (s *MainCategSuite) TestGetByID() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when has data, return successfully":     getByID_NoError_ReturnSuccessfully,
		"when find no data, return successfully": getByID_FindNoData_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByID_NoError_ReturnSuccessfully(s *MainCategSuite, desc string) {
	// prepare existing data
	categ, user, _, err := s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	// prepare more user data
	_, _, _, err = s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)
	_, _, _, err = s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	// prepare expected result
	expResult := domain.MainCateg{
		ID:   categ.ID,
		Name: categ.Name,
		Type: domain.CvtToTransactionType(categ.Type),
	}

	result, err := s.mainCategModel.GetByID(categ.ID, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, *result, desc)
}

func getByID_FindNoData_ReturnSuccessfully(s *MainCategSuite, desc string) {
	// prepare existing data
	_, user, _, err := s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	// prepare more user data
	_, _, _, err = s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)
	_, _, _, err = s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	result, err := s.mainCategModel.GetByID(0, user.ID)
	s.Require().Equal(domain.ErrMainCategNotFound, err, desc)
	s.Require().Nil(result, desc)
}

func (s *MainCategSuite) TestCreateBatch() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when insert one data, insert successfully":            createBatch_InsertOneData_InsertSuccessfully,
		"when insert many data, insert successfully":           createBatch_InsertManyData_InsertSuccessfully,
		"when insert duplicate name data, return error":        createBatch_InsertDuplicateNameData_ReturnError,
		"when already exist data with same name, return error": createBatch_AlreadyExistData_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func createBatch_InsertOneData_InsertSuccessfully(s *MainCategSuite, desc string) {
	users, icons, err := s.f.InsertUsersAndIcons(1, 1)
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name: "test1",
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{
				ID: icons[0].ID,
			},
		},
	}

	err = s.mainCategModel.BatchCreate(mockCtx, categs, users[0].ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE user_id = ?
							 `
	var result MainCateg
	err = s.db.QueryRow(checkStmt, users[0].ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
	s.Require().NoError(err, desc)
	s.Require().Equal(categs[0].Name, result.Name, desc)
	s.Require().Equal(categs[0].Type.ToModelValue(), result.Type, desc)
	s.Require().Equal(icons[0].ID, result.IconID, desc)
}

func createBatch_InsertManyData_InsertSuccessfully(s *MainCategSuite, desc string) {
	users, icons, err := s.f.InsertUsersAndIcons(1, 3)
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name: "test1",
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{ID: icons[0].ID},
		},
		{
			Name: "test2",
			Type: domain.TransactionTypeIncome,
			Icon: domain.Icon{ID: icons[1].ID},
		},
		{
			Name: "test3",
			Type: domain.TransactionTypeIncome,
			Icon: domain.Icon{ID: icons[2].ID},
		},
	}

	err = s.mainCategModel.BatchCreate(mockCtx, categs, users[0].ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_id
							 FROM main_categories
							 WHERE user_id = ?
							 ORDER BY id
							 `
	rows, err := s.db.Query(checkStmt, users[0].ID)
	s.Require().NoError(err, desc)
	defer rows.Close()

	var result MainCateg
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&result.ID, &result.Name, &result.Type, &result.IconID)
		s.Require().NoError(err, desc)
		s.Require().Equal(categs[i].Name, result.Name, desc)
		s.Require().Equal(categs[i].Type.ToModelValue(), result.Type, desc)
		s.Require().Equal(icons[i].ID, result.IconID, desc)
	}
}

func createBatch_InsertDuplicateNameData_ReturnError(s *MainCategSuite, desc string) {
	users, icons, err := s.f.InsertUsersAndIcons(1, 2)
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name: "test1",
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{ID: icons[0].ID},
		},
		{
			Name: "test1",
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{ID: icons[1].ID},
		},
	}

	err = s.mainCategModel.BatchCreate(mockCtx, categs, users[0].ID)
	s.Require().EqualError(err, domain.ErrUniqueNameUserType.Error(), desc)

	// check if the data is not inserted
	stms := `SELECT COUNT(*)
						FROM main_categories
						WHERE user_id = ?
						`
	var count int
	err = s.db.QueryRow(stms, users[0].ID).Scan(&count)
	s.Require().NoError(err, desc)
	s.Require().Equal(0, count, desc)
}

func createBatch_AlreadyExistData_ReturnError(s *MainCategSuite, desc string) {
	maincateg, user, _, err := s.f.InsertMainCategWithAss(MainCateg{})
	s.Require().NoError(err, desc)

	icons, err := s.f.Icon.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name: "test1",
			Type: domain.TransactionTypeExpense,
			Icon: domain.Icon{ID: icons[0].ID},
		},
		{
			Name: maincateg.Name,
			Type: domain.CvtToTransactionType(maincateg.Type),
			Icon: domain.Icon{ID: icons[1].ID},
		},
	}

	err = s.mainCategModel.BatchCreate(mockCtx, categs, user.ID)
	s.Require().EqualError(err, domain.ErrUniqueNameUserType.Error(), desc)

	// check if the data is not inserted
	stms := `SELECT COUNT(*)
						FROM main_categories
						WHERE user_id = ?
						`
	var count int
	err = s.db.QueryRow(stms, user.ID).Scan(&count)
	s.Require().NoError(err, desc)
}
