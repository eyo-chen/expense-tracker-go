package maincateg

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
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

type MainCategSuite struct {
	suite.Suite
	db            *sql.DB
	dk            *dockerutil.Container
	migrate       *migrate.Migrate
	f             *factory
	mainCategRepo *Repo
	userRepo      *user.Repo
	iconRepo      *icon.Repo
}

func TestMainCategSuite(t *testing.T) {
	suite.Run(t, new(MainCategSuite))
}

func (s *MainCategSuite) SetupSuite() {
	s.dk = dockerutil.RunDocker(dockerutil.ImageMySQL)
	db, migrate := testutil.ConnToDB(s.dk.Port)
	logger.Register()
	s.db = db
	s.migrate = migrate
	s.mainCategRepo = New(db)
	s.userRepo = user.New(db)
	s.iconRepo = icon.New(db)

	s.f = newFactory(db)
}

func (s *MainCategSuite) TearDownSuite() {
	if err := s.db.Close(); err != nil {
		logger.Error("Unable to close mysql database", "error", err)
	}
	s.migrate.Close()
	s.dk.PurgeDocker()
}

func (s *MainCategSuite) SetupTest() {
	s.mainCategRepo = New(s.db)
	s.userRepo = user.New(s.db)
	s.iconRepo = icon.New(s.db)
	s.f = newFactory(s.db)
}

func (s *MainCategSuite) TearDownTest() {
	tx, err := s.db.Begin()
	s.Require().NoError(err)
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.Require().NoError(err)
		}
	}()

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
	users, err := s.f.InsertUsers(mockCTX, 1)
	s.Require().NoError(err, desc)

	categ := domain.MainCateg{
		Name:     "test",
		Type:     domain.TransactionTypeExpense,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	err = s.mainCategRepo.Create(mockCTX, categ, users[0].ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_type, icon_data
							 FROM main_categories
							 WHERE user_id = ?
							 AND name = ?
							 AND type = ?
							 `
	var result MainCateg
	err = s.db.QueryRow(checkStmt, users[0].ID, "test", domain.TransactionTypeExpense.ToModelValue()).Scan(&result.ID, &result.Name, &result.Type, &result.IconType, &result.IconData)
	s.Require().NoError(err, desc)
	s.Require().Equal(categ.Name, result.Name, desc)
	s.Require().Equal(categ.Type.ToModelValue(), result.Type, desc)
	s.Require().Equal(categ.IconType.ToModelValue(), result.IconType, desc)
	s.Require().Equal(categ.IconData, result.IconData, desc)
}

func create_DuplicateName_ReturnError(s *MainCategSuite, desc string) {
	createdMainCateg, user, err := s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	categ := domain.MainCateg{
		Name:     createdMainCateg.Name,
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "url",
	}
	err = s.mainCategRepo.Create(mockCTX, categ, user.ID)
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
	mainCategList, users, err := s.f.InsertMainCategListWithAss(mockCTX, 2, 1, 2, "expense", "income")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:       mainCategList[1].ID,
			Name:     mainCategList[1].Name,
			Type:     domain.TransactionTypeIncome,
			IconType: domain.CvtToIconType(mainCategList[1].IconType),
			IconData: mainCategList[1].IconData,
		},
	}

	categs, err := s.mainCategRepo.GetAll(mockCTX, users[0].ID, domain.TransactionTypeIncome)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, categs, desc)
}

func getAll_ExpenseType_ReturnOnlyExpenseTypeData(s *MainCategSuite, desc string) {
	mainCategList, users, err := s.f.InsertMainCategListWithAss(mockCTX, 2, 1, 2, "expense", "income")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:       mainCategList[0].ID,
			Name:     mainCategList[0].Name,
			Type:     domain.TransactionTypeExpense,
			IconType: domain.CvtToIconType(mainCategList[0].IconType),
			IconData: mainCategList[0].IconData,
		},
	}

	categs, err := s.mainCategRepo.GetAll(mockCTX, users[0].ID, domain.TransactionTypeExpense)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, categs, desc)
}

func getAll_UnSpecifiedType_ReturnAllData(s *MainCategSuite, desc string) {
	mainCategList, users, err := s.f.InsertMainCategListWithAss(mockCTX, 2, 1, 2, "expense", "income")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:       mainCategList[0].ID,
			Name:     mainCategList[0].Name,
			Type:     domain.TransactionTypeExpense,
			IconType: domain.CvtToIconType(mainCategList[0].IconType),
			IconData: mainCategList[0].IconData,
		},
		{
			ID:       mainCategList[1].ID,
			Name:     mainCategList[1].Name,
			Type:     domain.TransactionTypeIncome,
			IconType: domain.CvtToIconType(mainCategList[1].IconType),
			IconData: mainCategList[1].IconData,
		},
	}

	categs, err := s.mainCategRepo.GetAll(mockCTX, users[0].ID, domain.TransactionTypeUnSpecified)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, categs, desc)
}

func getAll_MultipleUsers_ReturnCorrectData(s *MainCategSuite, desc string) {
	mainCategList, users, err := s.f.InsertMainCategListWithAss(mockCTX, 3, 2, 3, "expense", "income", "expense")
	s.Require().NoError(err, desc)

	expResult := []domain.MainCateg{
		{
			ID:       mainCategList[1].ID,
			Name:     mainCategList[1].Name,
			Type:     domain.TransactionTypeIncome,
			IconType: domain.CvtToIconType(mainCategList[1].IconType),
			IconData: mainCategList[1].IconData,
		},
		{
			ID:       mainCategList[2].ID,
			Name:     mainCategList[2].Name,
			Type:     domain.TransactionTypeExpense,
			IconType: domain.CvtToIconType(mainCategList[2].IconType),
			IconData: mainCategList[2].IconData,
		},
	}

	categs, err := s.mainCategRepo.GetAll(mockCTX, users[1].ID, domain.TransactionTypeUnSpecified)
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
	categ, _, err := s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	inputCateg := domain.MainCateg{
		ID:       categ.ID,
		Name:     "test2",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "new-url",
	}
	err = s.mainCategRepo.Update(mockCTX, inputCateg)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_type, icon_data
							 FROM main_categories
							 WHERE id = ?
							 `
	var result MainCateg
	err = s.db.QueryRow(checkStmt, categ.ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconType, &result.IconData)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputCateg.Name, result.Name, desc)
	s.Require().Equal(inputCateg.Type.ToModelValue(), result.Type, desc)
	s.Require().Equal(inputCateg.IconType.ToModelValue(), result.IconType, desc)
	s.Require().Equal(inputCateg.IconData, result.IconData, desc)
}

func update_WithMultipleUser_UpdateSuccessfully(s *MainCategSuite, desc string) {
	categs, _, err := s.f.InsertMainCategListWithAss(mockCTX, 2, 2, 1)
	s.Require().NoError(err, desc)

	inputCateg := domain.MainCateg{
		ID:       categs[0].ID,
		Name:     "update name",
		Type:     domain.TransactionTypeIncome,
		IconType: domain.IconTypeDefault,
		IconData: "new-url",
	}
	err = s.mainCategRepo.Update(mockCTX, inputCateg)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_type, icon_data
							 FROM main_categories
							 WHERE id = ?
							 `
	// check if the data is updated
	var result MainCateg
	err = s.db.QueryRow(checkStmt, categs[0].ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconType, &result.IconData)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputCateg.Name, result.Name, desc)
	s.Require().Equal(inputCateg.Type.ToModelValue(), result.Type, desc)
	s.Require().Equal(inputCateg.IconType.ToModelValue(), result.IconType, desc)
	s.Require().Equal(inputCateg.IconData, result.IconData, desc)

	// check if the other data is not updated
	var result2 MainCateg
	err = s.db.QueryRow(checkStmt, categs[1].ID).Scan(&result2.ID, &result2.Name, &result2.Type, &result2.IconType, &result2.IconData)
	s.Require().NoError(err, desc)
	s.Require().Equal(categs[1].Name, result2.Name, desc)
	s.Require().Equal(categs[1].Type, result2.Type, desc)
	s.Require().Equal(categs[1].IconType, result2.IconType, desc)
	s.Require().Equal(categs[1].IconData, result2.IconData, desc)
}

func update_DuplicateName_ReturnError(s *MainCategSuite, desc string) {
	categs, _, err := s.f.InsertMainCategListWithAss(mockCTX, 2, 1, 2)
	s.Require().NoError(err, desc)

	domainMainCateg := domain.MainCateg{
		ID:       categs[0].ID,
		Name:     categs[1].Name, // update categ1 with categ2 name
		Type:     domain.CvtToTransactionType(categs[0].Type),
		IconType: domain.IconTypeDefault,
		IconData: "new-url",
	}
	err = s.mainCategRepo.Update(mockCTX, domainMainCateg)
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
	categ, _, err := s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	err = s.mainCategRepo.Delete(categ.ID)
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
	categ, user, err := s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	// prepare more user data
	_, _, err = s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)
	_, _, err = s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	// prepare expected result
	expResult := domain.MainCateg{
		ID:       categ.ID,
		Name:     categ.Name,
		Type:     domain.CvtToTransactionType(categ.Type),
		IconType: domain.CvtToIconType(categ.IconType),
		IconData: categ.IconData,
	}

	result, err := s.mainCategRepo.GetByID(categ.ID, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, *result, desc)
}

func getByID_FindNoData_ReturnSuccessfully(s *MainCategSuite, desc string) {
	// prepare existing data
	_, user, err := s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	// prepare more user data
	_, _, err = s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)
	_, _, err = s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	result, err := s.mainCategRepo.GetByID(0, user.ID)
	s.Require().ErrorIs(err, domain.ErrMainCategNotFound, desc)
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
	users, err := s.f.InsertUsers(mockCTX, 1)
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name:     "test1",
			Type:     domain.TransactionTypeExpense,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
	}

	err = s.mainCategRepo.BatchCreate(mockCTX, categs, users[0].ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_type, icon_data
							 FROM main_categories
							 WHERE user_id = ?
							 `
	var result MainCateg
	err = s.db.QueryRow(checkStmt, users[0].ID).Scan(&result.ID, &result.Name, &result.Type, &result.IconType, &result.IconData)
	s.Require().NoError(err, desc)
	s.Require().Equal(categs[0].Name, result.Name, desc)
	s.Require().Equal(categs[0].Type.ToModelValue(), result.Type, desc)
	s.Require().Equal(categs[0].IconType.ToModelValue(), result.IconType, desc)
	s.Require().Equal(categs[0].IconData, result.IconData, desc)
}

func createBatch_InsertManyData_InsertSuccessfully(s *MainCategSuite, desc string) {
	users, err := s.f.InsertUsers(mockCTX, 1)
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name:     "test1",
			Type:     domain.TransactionTypeExpense,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
		{
			Name:     "test2",
			Type:     domain.TransactionTypeIncome,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
		{
			Name:     "test3",
			Type:     domain.TransactionTypeIncome,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
	}

	err = s.mainCategRepo.BatchCreate(mockCTX, categs, users[0].ID)
	s.Require().NoError(err, desc)

	checkStmt := `SELECT id, name, type, icon_type, icon_data
							 FROM main_categories
							 WHERE user_id = ?
							 ORDER BY id
							 `
	rows, err := s.db.Query(checkStmt, users[0].ID)
	s.Require().NoError(err, desc)
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Unable to close rows", "package", packageName, "err", err)
		}
	}()

	var result MainCateg
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&result.ID, &result.Name, &result.Type, &result.IconType, &result.IconData)
		s.Require().NoError(err, desc)
		s.Require().Equal(categs[i].Name, result.Name, desc)
		s.Require().Equal(categs[i].Type.ToModelValue(), result.Type, desc)
		s.Require().Equal(categs[i].IconType.ToModelValue(), result.IconType, desc)
		s.Require().Equal(categs[i].IconData, result.IconData, desc)
	}
}

func createBatch_InsertDuplicateNameData_ReturnError(s *MainCategSuite, desc string) {
	users, err := s.f.InsertUsers(mockCTX, 1)
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name:     "test1",
			Type:     domain.TransactionTypeExpense,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
		{
			Name:     "test1",
			Type:     domain.TransactionTypeExpense,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
	}

	err = s.mainCategRepo.BatchCreate(mockCTX, categs, users[0].ID)
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
	maincateg, user, err := s.f.InsertMainCategWithAss(mockCTX, MainCateg{})
	s.Require().NoError(err, desc)

	categs := []domain.MainCateg{
		{
			Name:     "test1",
			Type:     domain.TransactionTypeExpense,
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
		{
			Name:     maincateg.Name,
			Type:     domain.CvtToTransactionType(maincateg.Type),
			IconType: domain.IconTypeDefault,
			IconData: "url",
		},
	}

	err = s.mainCategRepo.BatchCreate(mockCTX, categs, user.ID)
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
