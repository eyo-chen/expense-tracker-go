package subcateg

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

type SubCategSuite struct {
	suite.Suite
	subCategModel interfaces.SubCategModel
	db            *sql.DB
	migrate       *migrate.Migrate
	f             *factory
}

func TestSubCategSuite(t *testing.T) {
	suite.Run(t, new(SubCategSuite))
}

func (s *SubCategSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db, migrate := testutil.ConnToDB(port)
	logger.Register()
	s.db = db
	s.subCategModel = NewSubCategModel(db)
	s.migrate = migrate
	s.f = newFactory(db)
}

func (s *SubCategSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *SubCategSuite) SetupTest() {
	s.subCategModel = NewSubCategModel(s.db)
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
	// prepare existing data
	user, maincateg, err := s.f.InsertUserAndMaincateg()
	s.Require().NoError(err, desc)

	// prepare input data
	subCateg := &domain.SubCateg{
		Name:        "test",
		MainCategID: maincateg.ID,
	}

	// action
	err = s.subCategModel.Create(subCateg, user.ID)
	s.Require().NoError(err, desc)

	// check
	var result SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user.ID, maincateg.ID, subCateg.Name).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(subCateg.Name, result.Name, desc)
	s.Require().Equal(subCateg.MainCategID, result.MainCategID, desc)
}

func create_DuplicateNameUserMainCateg_ReturnError(s *SubCategSuite, desc string) {
	// prepare data
	subcategs, user, maincateg, err := s.f.InsertSubcategs(1)
	s.Require().NoError(err, desc)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		Name:        subcategs[0].Name,
		MainCategID: maincateg.ID,
	}

	// action and check
	err = s.subCategModel.Create(inputSubCateg, user.ID)
	s.Require().Equal(domain.ErrUniqueNameUserMainCateg, err, desc)
}

func (s *SubCategSuite) TestGetByMainCategID() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when find no data, return empty":                              getByMainCategID_FindNoData_ReturnEmpty,
		"when with one main category, return correct subcategories":    getByMainCategID_WithOneMainCateg_ReturnCorrectSubCategs,
		"when with many main categories, return correct subcategories": getByMainCategID_WithManyMainCategs_ReturnCorrectSubCategs,
		"when with many users, return correct subcategories":           getByMainCategID_WithManyUsers_ReturnCorrectSubCategs,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByMainCategID_FindNoData_ReturnEmpty(s *SubCategSuite, desc string) {
	// prepare data
	_, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{2})
	s.Require().NoError(err, desc)

	// expected result
	mainCateg := mainCategs[0]

	// action
	result, err := s.subCategModel.GetByMainCategID(user.ID, mainCateg.ID+9999)
	s.Require().NoError(err, desc)
	s.Require().Nil(result, desc)
}

func getByMainCategID_WithOneMainCateg_ReturnCorrectSubCategs(s *SubCategSuite, desc string) {
	// prepare data
	mainCategIDToSubCategs, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{2})
	s.Require().NoError(err, desc)

	// expected result
	mainCateg := mainCategs[0]
	expResult := []*domain.SubCateg{
		{
			ID:          mainCategIDToSubCategs[mainCateg.ID][0].ID,
			Name:        mainCategIDToSubCategs[mainCateg.ID][0].Name,
			MainCategID: mainCategIDToSubCategs[mainCateg.ID][0].MainCategID,
		},
		{
			ID:          mainCategIDToSubCategs[mainCateg.ID][1].ID,
			Name:        mainCategIDToSubCategs[mainCateg.ID][1].Name,
			MainCategID: mainCategIDToSubCategs[mainCateg.ID][1].MainCategID,
		},
	}

	// action
	result, err := s.subCategModel.GetByMainCategID(user.ID, mainCateg.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getByMainCategID_WithManyMainCategs_ReturnCorrectSubCategs(s *SubCategSuite, desc string) {
	// prepare data
	mainCategIDToSubCategs, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(3, []int{3, 2, 1})
	s.Require().NoError(err, desc)

	// expected result
	mainCateg := mainCategs[2] // choose the third main category
	expResult := []*domain.SubCateg{
		{
			ID:          mainCategIDToSubCategs[mainCateg.ID][0].ID,
			Name:        mainCategIDToSubCategs[mainCateg.ID][0].Name,
			MainCategID: mainCategIDToSubCategs[mainCateg.ID][0].MainCategID,
		},
	}

	// action
	result, err := s.subCategModel.GetByMainCategID(user.ID, mainCateg.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getByMainCategID_WithManyUsers_ReturnCorrectSubCategs(s *SubCategSuite, desc string) {
	// prepare data
	mainCategIDToSubCategs, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(3, []int{3, 2, 1})
	s.Require().NoError(err, desc)

	// prepare more data with different user
	_, _, _, err = s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{1})
	s.Require().NoError(err, desc)
	_, _, _, err = s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{1})
	s.Require().NoError(err, desc)

	// expected result
	mainCateg := mainCategs[2] // choose the third main category
	expResult := []*domain.SubCateg{
		{
			ID:          mainCategIDToSubCategs[mainCateg.ID][0].ID,
			Name:        mainCategIDToSubCategs[mainCateg.ID][0].Name,
			MainCategID: mainCategIDToSubCategs[mainCateg.ID][0].MainCategID,
		},
	}

	// action
	result, err := s.subCategModel.GetByMainCategID(user.ID, mainCateg.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func (s *SubCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when no duplicate data, update successfully":                  update_NoDuplicateData_UpdateSuccessfully,
		"when there are multiple main categories, update successfully": update_WithMultipleMainCateg_UpdateSuccessfully,
		"when update to duplicate name, return error":                  update_DuplicateName_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoDuplicateData_UpdateSuccessfully(s *SubCategSuite, desc string) {
	// prepare existing data
	subcategs, user, maincateg, err := s.f.InsertSubcategs(1)
	s.Require().NoError(err, desc)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		ID:          subcategs[0].ID,
		Name:        "updated test",
		MainCategID: maincateg.ID,
	}

	// action
	err = s.subCategModel.Update(inputSubCateg)
	s.Require().NoError(err, desc)

	// check
	var result SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user.ID, maincateg.ID, inputSubCateg.Name).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputSubCateg.Name, result.Name, desc)
	s.Require().Equal(inputSubCateg.MainCategID, result.MainCategID, desc)
}

func update_WithMultipleMainCateg_UpdateSuccessfully(s *SubCategSuite, desc string) {
	// prepare existing data
	subcategs, user, maincateg, err := s.f.InsertSubcategs(2)
	s.Require().NoError(err, desc)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		ID:          subcategs[0].ID,
		Name:        "updated test",
		MainCategID: maincateg.ID,
	}

	// action
	err = s.subCategModel.Update(inputSubCateg)
	s.Require().NoError(err, desc)

	// check
	var result SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user.ID, maincateg.ID, inputSubCateg.Name).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(inputSubCateg.Name, result.Name, desc)
	s.Require().Equal(inputSubCateg.MainCategID, result.MainCategID, desc)

	// check the other main category
	var result2 SubCateg
	checkStmt = `SELECT id, name, main_category_id FROM sub_categories WHERE user_id = ? AND main_category_id = ? AND name = ?`
	err = s.db.QueryRow(checkStmt, user.ID, subcategs[1].MainCategID, subcategs[1].Name).Scan(&result2.ID, &result2.Name, &result2.MainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(subcategs[1].Name, result2.Name, desc)
	s.Require().Equal(subcategs[1].MainCategID, result2.MainCategID, desc)
}

func update_DuplicateName_ReturnError(s *SubCategSuite, desc string) {
	// prepare existing data
	subcategs, _, maincateg, err := s.f.InsertSubcategs(2)
	s.Require().NoError(err, desc)

	fmt.Println("subcategs", subcategs)

	// prepare input data
	inputSubCateg := &domain.SubCateg{
		ID:          subcategs[0].ID,
		Name:        subcategs[1].Name, // update to duplicate name
		MainCategID: maincateg.ID,
	}

	// action and check
	err = s.subCategModel.Update(inputSubCateg)
	s.Require().Equal(domain.ErrUniqueNameUserMainCateg, err, desc)
}

func (s *SubCategSuite) TestDelete() {
	mainCategIDToSubCategs, mainCategs, _, err := s.f.InsertSubcategsWithOneOrManyMainCateg(3, []int{3, 2, 1})
	s.Require().NoError(err, "test delete")

	// prepare more data with different user
	_, _, _, err = s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{1})
	s.Require().NoError(err, "test delete")
	_, _, _, err = s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{1})
	s.Require().NoError(err, "test delete")

	mainCateg := mainCategs[0] // choose the first main category

	// action
	err = s.subCategModel.Delete(mainCategIDToSubCategs[mainCateg.ID][0].ID)
	s.Require().NoError(err, "test delete")

	// check to see if the sub category is deleted
	var result SubCateg
	checkStmt := `SELECT id, name, main_category_id FROM sub_categories WHERE id = ?`
	err = s.db.QueryRow(checkStmt, mainCategIDToSubCategs[mainCateg.ID][0].ID).Scan(&result.ID, &result.Name, &result.MainCategID)
	s.Require().Equal(sql.ErrNoRows, err, "test delete")

	// check to see if the first main category still has the other sub categories
	checkStmt = `SELECT id, name, main_category_id FROM sub_categories WHERE main_category_id = ?`
	rows, err := s.db.Query(checkStmt, mainCateg.ID)
	s.Require().NoError(err, "test delete")
	defer rows.Close()
	var subCategs []SubCateg
	for rows.Next() {
		var subCateg SubCateg
		err := rows.Scan(&subCateg.ID, &subCateg.Name, &subCateg.MainCategID)
		s.Require().NoError(err, "test delete")
		subCategs = append(subCategs, subCateg)
	}
	s.Require().Len(subCategs, 2, "test delete")

	// check to see if the second main category still has the sub category
	checkStmt = `SELECT id, name, main_category_id FROM sub_categories WHERE main_category_id = ?`
	rows, err = s.db.Query(checkStmt, mainCategs[1].ID)
	s.Require().NoError(err, "test delete")
	defer rows.Close()
	var subCategs2 []SubCateg
	for rows.Next() {
		var subCateg SubCateg
		err := rows.Scan(&subCateg.ID, &subCateg.Name, &subCateg.MainCategID)
		s.Require().NoError(err, "test delete")
		subCategs2 = append(subCategs2, subCateg)
	}
	s.Require().Len(subCategs2, 2, "test delete")

	// check to see if the third main category still has the sub category
	checkStmt = `SELECT id, name, main_category_id FROM sub_categories WHERE main_category_id = ?`
	rows, err = s.db.Query(checkStmt, mainCategs[2].ID)
	s.Require().NoError(err, "test delete")
	defer rows.Close()
	var subCategs3 []SubCateg
	for rows.Next() {
		var subCateg SubCateg
		err := rows.Scan(&subCateg.ID, &subCateg.Name, &subCateg.MainCategID)
		s.Require().NoError(err, "test delete")
		subCategs3 = append(subCategs3, subCateg)
	}
	s.Require().Len(subCategs3, 1, "test delete")
}

func (s *SubCategSuite) TestGetByID() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when find no data, return error":           getByID_FindNoData_ReturnError,
		"when with one data, return correct data":   getByID_WithOneData_ReturnCorrectData,
		"when with many data, return correct data":  getByID_WithManyData_ReturnCorrectData,
		"when with many users, return correct data": getByID_WithManyUsers_ReturnCorrectData,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByID_FindNoData_ReturnError(s *SubCategSuite, desc string) {
	// prepare data
	_, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(3, []int{3, 2, 1})
	s.Require().NoError(err, desc)

	mainCateg := mainCategs[0]

	// action
	result, err := s.subCategModel.GetByID(mainCateg.ID+999, user.ID)
	s.Require().Equal(domain.ErrSubCategNotFound, err, desc)
	s.Require().Nil(result, desc)
}

func getByID_WithOneData_ReturnCorrectData(s *SubCategSuite, desc string) {
	// prepare data
	mainCategIDToSubCategs, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{1})
	s.Require().NoError(err, desc)

	mainCateg := mainCategs[0]
	subCateg := mainCategIDToSubCategs[mainCateg.ID][0]
	// prepare expected result
	expResult := &domain.SubCateg{
		ID:          subCateg.ID,
		Name:        subCateg.Name,
		MainCategID: subCateg.MainCategID,
	}

	// action
	result, err := s.subCategModel.GetByID(subCateg.ID, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getByID_WithManyData_ReturnCorrectData(s *SubCategSuite, desc string) {
	// prepare data
	mainCategIDToSubCategs, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(3, []int{3, 2, 1})
	s.Require().NoError(err, desc)

	// prepare expected result
	mainCateg := mainCategs[1]
	subCateg := mainCategIDToSubCategs[mainCateg.ID][1]
	expResult := &domain.SubCateg{
		ID:          subCateg.ID,
		Name:        subCateg.Name,
		MainCategID: subCateg.MainCategID,
	}

	// action
	result, err := s.subCategModel.GetByID(subCateg.ID, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getByID_WithManyUsers_ReturnCorrectData(s *SubCategSuite, desc string) {
	// prepare data
	mainCategIDToSubCategs, mainCategs, user, err := s.f.InsertSubcategsWithOneOrManyMainCateg(3, []int{3, 2, 1})
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, err = s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{1})
	s.Require().NoError(err, desc)
	_, _, _, err = s.f.InsertSubcategsWithOneOrManyMainCateg(1, []int{1})
	s.Require().NoError(err, desc)

	mainCateg := mainCategs[2]
	subCateg := mainCategIDToSubCategs[mainCateg.ID][0]
	// prepare expected result
	expResult := &domain.SubCateg{
		ID:          subCateg.ID,
		Name:        subCateg.Name,
		MainCategID: subCateg.MainCategID,
	}

	// action
	result, err := s.subCategModel.GetByID(subCateg.ID, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}
