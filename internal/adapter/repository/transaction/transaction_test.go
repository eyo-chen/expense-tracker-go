package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/codeutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/dockerutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX     = context.Background()
	mockLoc, _  = time.LoadLocation("")
	mockTimeNow = time.Unix(1629446406, 0).Truncate(24 * time.Hour).In(mockLoc)
)

type TransactionSuite struct {
	suite.Suite
	dk      *dockerutil.Container
	db      *sql.DB
	migrate *migrate.Migrate
	repo    *Repo
	f       *TransactionFactory
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) SetupSuite() {
	s.dk = dockerutil.RunDocker(dockerutil.ImageMySQL)
	db, migrate := testutil.ConnToDB(s.dk.Port)
	logger.Register()

	s.db = db
	s.migrate = migrate
	s.repo = New(s.db)
	s.f = NewTransactionFactory(db)
}

func (s *TransactionSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	s.dk.PurgeDocker()
}

func (s *TransactionSuite) SetupTest() {
	s.repo = New(s.db)
	s.f = NewTransactionFactory(s.db)
}

func (s *TransactionSuite) TearDownTest() {
	tx, err := s.db.Begin()
	if err != nil {
		s.Require().NoError(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.Require().NoError(err)
		}
	}()

	if _, err := tx.Exec("DELETE FROM transactions"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM icons"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM sub_categories"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM main_categories"); err != nil {
		s.Require().NoError(err)
	}

	if _, err := tx.Exec("DELETE FROM users"); err != nil {
		s.Require().NoError(err)
	}

	s.Require().NoError(tx.Commit())
	s.f.Reset()
}

func (s *TransactionSuite) TestCreate() {
	user, main, sub, err := s.f.PrepareUserMainAndSubCateg(mockCTX)
	s.Require().NoError(err)

	t := domain.CreateTransactionInput{
		UserID:      user.ID,
		Type:        domain.CvtToTransactionType(main.Type),
		MainCategID: main.ID,
		SubCategID:  sub.ID,
		Price:       100,
		Note:        "test",
		Date:        mockTimeNow,
	}

	err = s.repo.Create(mockCTX, t)
	s.Require().NoError(err)

	var checkT Transaction
	stmt := "SELECT user_id, type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE user_id = ?"
	err = s.db.QueryRow(stmt, user.ID).Scan(&checkT.UserID, &checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Equal(t.UserID, checkT.UserID)
	s.Equal(t.Type.ToModelValue(), checkT.Type)
	s.Equal(t.MainCategID, checkT.MainCategID)
	s.Equal(t.SubCategID, checkT.SubCategID)
	s.Equal(t.Price, checkT.Price)
	s.Equal(t.Note, checkT.Note)

	s.TearDownTest()
}

func (s *TransactionSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return successfully":                                          getAll_NoError_ReturnSuccessfully,
		"when with multiple users, return successfully":                               getAll_WithMultipleUsers_ReturnSuccessfully,
		"when with many transactions, return all transactions":                        getAll_WithManyTransaction_ReturnSuccessfully,
		"when with search keyword, return data with keyword":                          getAll_WithSearchKeyword_ReturnDataWithKeyword,
		"when filter by start date, return data after start date":                     getAll_FilterByStartDate_ReturnDataAfterStartDate,
		"when filter by end date, return data before end date":                        getAll_FilterByEndDate_ReturnDataBeforeEndDate,
		"when filter by start and end date, return data between them":                 getAll_FilterByStartAndEndDate_ReturnDataBetweenStartAndEndDate,
		"when filter by min price, return data greater than min price":                getAll_FilterByMinPrice_ReturnDataGreaterThanMinPrice,
		"when filter by max price, return data less than max price":                   getAll_FilterByMaxPrice_ReturnDataLessThanMinPrice,
		"when filter by main category id, return data with main category id":          getAll_FilterByMainCategID_ReturnDataWithMainCategID,
		"when filter by sub category id, return data with sub category id":            getAll_FilterBySubCategID_ReturnDataWithSubCategID,
		"when filter by date, price and main category, return correct data":           getAll_FilterByDateAndPriceAndMainCateg_ReturnCorrectData,
		"when sort by date asc, return correct order":                                 getAll_SortByDateAsc_ReturnCorrectOrder,
		"when sort by date desc, return correct order":                                getAll_SortByDateDesc_ReturnCorrectOrder,
		"when sort by price asc, return correct order":                                getAll_SortByPriceAsc_ReturnCorrectOrder,
		"when sort by price desc, return correct order":                               getAll_SortByPriceDesc_ReturnCorrectOrder,
		"when sort by type asc, return correct order":                                 getAll_SortByTypeAsc_ReturnCorrectOrder,
		"when sort by type desc, return correct order":                                getAll_SortByTypeDesc_ReturnCorrectOrder,
		"when with next key cursor, return data after cursor key":                     getAll_WithNextKeyCursor_ReturnDataAfterCursorKey,
		"when with next key cursor and sort by date, return correct data":             getAll_WithNextKeyCursorAndSortByDate_ReturnCorrectData,
		"when with next key cursor and sort by price, return correct data":            getAll_WithNextKeyCursorAndSortByPrice_ReturnCorrectData,
		"when with next key cursor and sort by transaction type, return correct data": getAll_WithNextKeyCursorAndSortByTransType_ReturnCorrectData,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAll_NoError_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0)

	opt := domain.GetTransOpt{}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0)

	opt := domain.GetTransOpt{}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_WithSearchKeyword_ReturnDataWithKeyword(s *TransactionSuite, desc string) {
	ow1 := Transaction{Note: "mysql database"}
	ow2 := Transaction{Note: "postgresql database"}
	ow3 := Transaction{Note: "mongodb database"}
	ow4 := Transaction{Note: "go programming"}
	ow5 := Transaction{Note: "javascript programming"}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 5, ow1, ow2, ow3, ow4, ow5)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 1, 2)

	searchKeyword := "database"
	opt := domain.GetTransOpt{
		Search: domain.Search{
			Keyword: &searchKeyword,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_WithManyTransaction_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 3)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 1, 2)

	opt := domain.GetTransOpt{}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterByStartDate_ReturnDataAfterStartDate(s *TransactionSuite, desc string) {
	ow1 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 1, 2, 3)

	startDate := mockTimeNow.AddDate(0, 0, -2)

	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			StartDate: &startDate,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterByEndDate_ReturnDataBeforeEndDate(s *TransactionSuite, desc string) {
	ow1 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 1)

	endDate := mockTimeNow.AddDate(0, 0, -2)
	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			EndDate: &endDate,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterByMinPrice_ReturnDataGreaterThanMinPrice(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 1000}
	ow2 := Transaction{Price: 1500}
	ow3 := Transaction{Price: 2000}
	ow4 := Transaction{Price: 2500}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 1, 2, 3)

	minPrice := 1500.00
	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			MinPrice: &minPrice,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterByMaxPrice_ReturnDataLessThanMinPrice(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 1000}
	ow2 := Transaction{Price: 1500}
	ow3 := Transaction{Price: 2000}
	ow4 := Transaction{Price: 2500}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 1)

	maxPrice := 1500.00
	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			MaxPrice: &maxPrice,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterByStartAndEndDate_ReturnDataBetweenStartAndEndDate(s *TransactionSuite, desc string) {
	ow1 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 1, 2)

	startDate := mockTimeNow.AddDate(0, 0, -2)
	endDate := mockTimeNow.AddDate(0, 0, -1)
	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			StartDate: &startDate,
			EndDate:   &endDate,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterByMainCategID_ReturnDataWithMainCategID(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 2)

	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			MainCategIDs: []int64{mainList[0].ID, mainList[2].ID},
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterBySubCategID_ReturnDataWithSubCategID(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 1, 3)

	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			SubCategIDs: []int64{subList[1].ID, subList[3].ID},
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_FilterByDateAndPriceAndMainCateg_ReturnCorrectData(s *TransactionSuite, desc string) {
	ow1 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3), Price: 1000}
	ow2 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2), Price: 2000}
	ow3 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1), Price: 3000}
	ow4 := Transaction{Date: mockTimeNow.AddDate(0, 0, 0), Price: 4000}
	ow5 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3), Price: 1000}
	ow6 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2), Price: 2000}
	ow7 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1), Price: 3000}
	ow8 := Transaction{Date: mockTimeNow.AddDate(0, 0, 0), Price: 4000}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 8, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 1, 2)

	startDate := mockTimeNow.AddDate(0, 0, -2)
	endDate := mockTimeNow.AddDate(0, 0, 0)
	minPrice := 2000.00
	maxPrice := 3000.00
	opt := domain.GetTransOpt{
		Filter: domain.Filter{
			StartDate:    &startDate,
			EndDate:      &endDate,
			MinPrice:     &minPrice,
			MaxPrice:     &maxPrice,
			MainCategIDs: []int64{mainList[1].ID, mainList[2].ID},
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_SortByDateAsc_ReturnCorrectOrder(s *TransactionSuite, desc string) {
	ow1 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow3 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow4 := Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	ow5 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}

	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 5, ow1, ow2, ow3, ow4, ow5)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 2, 4, 1, 3)

	opt := domain.GetTransOpt{
		Sort: &domain.Sort{
			By:  domain.SortByTypeDate,
			Dir: domain.SortDirTypeAsc,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_SortByDateDesc_ReturnCorrectOrder(s *TransactionSuite, desc string) {
	ow1 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow3 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow4 := Transaction{Date: mockTimeNow.AddDate(0, 0, 0)}
	ow5 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}

	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 5, ow1, ow2, ow3, ow4, ow5)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 3, 1, 4, 2, 0)

	opt := domain.GetTransOpt{
		Sort: &domain.Sort{
			By:  domain.SortByTypeDate,
			Dir: domain.SortDirTypeDesc,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_SortByPriceAsc_ReturnCorrectOrder(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 2500}
	ow2 := Transaction{Price: 1500}
	ow3 := Transaction{Price: 500}
	ow4 := Transaction{Price: 2000}
	ow5 := Transaction{Price: 500}

	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 5, ow1, ow2, ow3, ow4, ow5)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 2, 4, 1, 3, 0)

	opt := domain.GetTransOpt{
		Sort: &domain.Sort{
			By:  domain.SortByTypePrice,
			Dir: domain.SortDirTypeAsc,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_SortByPriceDesc_ReturnCorrectOrder(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 2500}
	ow2 := Transaction{Price: 1500}
	ow3 := Transaction{Price: 500}
	ow4 := Transaction{Price: 2000}
	ow5 := Transaction{Price: 500}

	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 5, ow1, ow2, ow3, ow4, ow5)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 3, 1, 4, 2)

	opt := domain.GetTransOpt{
		Sort: &domain.Sort{
			By:  domain.SortByTypePrice,
			Dir: domain.SortDirTypeDesc,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_SortByTypeAsc_ReturnCorrectOrder(s *TransactionSuite, desc string) {
	ow1 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow2 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue()}
	ow3 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow4 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow5 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue()}

	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 5, ow1, ow2, ow3, ow4, ow5)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 0, 2, 3, 1, 4)

	opt := domain.GetTransOpt{
		Sort: &domain.Sort{
			By:  domain.SortByTypeTransType,
			Dir: domain.SortDirTypeAsc,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_SortByTypeDesc_ReturnCorrectOrder(s *TransactionSuite, desc string) {
	ow1 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow2 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue()}
	ow3 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow4 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow5 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue()}

	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 5, ow1, ow2, ow3, ow4, ow5)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 4, 1, 3, 2, 0)

	opt := domain.GetTransOpt{
		Sort: &domain.Sort{
			By:  domain.SortByTypeTransType,
			Dir: domain.SortDirTypeDesc,
		},
	}
	trans, decodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Empty(decodedNextKey, desc)
}

func getAll_WithNextKeyCursor_ReturnDataAfterCursorKey(s *TransactionSuite, desc string) {
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 8)
	s.Require().NoError(err, desc)

	// prepare encodedNextKey
	// Note that the order of the transactionList is based on the ID(default), and it's ascending
	encodedNextKey, err := codeutil.EncodeNextKeys(domain.DecodedNextKeys{{Field: "ID", Value: "1"}}, transactionList[4])
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 5, 6, 7)

	opt := domain.GetTransOpt{
		Cursor: domain.Cursor{
			NextKey: encodedNextKey,
			Size:    3,
		},
	}
	trans, deencodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	decodedNextKeyID := transactionList[4].ID
	s.Require().Equal(domain.DecodedNextKeys{{Field: "ID", Value: fmt.Sprint(decodedNextKeyID)}}, deencodedNextKey, desc)
}

func getAll_WithNextKeyCursorAndSortByDate_ReturnCorrectData(s *TransactionSuite, desc string) {
	ow1 := Transaction{Date: mockTimeNow.AddDate(0, 0, -5)}
	ow2 := Transaction{Date: mockTimeNow.AddDate(0, 0, -4)}
	ow3 := Transaction{Date: mockTimeNow.AddDate(0, 0, -4)}
	ow4 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow5 := Transaction{Date: mockTimeNow.AddDate(0, 0, -3)}
	ow6 := Transaction{Date: mockTimeNow.AddDate(0, 0, -2)}
	ow7 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	ow8 := Transaction{Date: mockTimeNow.AddDate(0, 0, -1)}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 8, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8)
	s.Require().NoError(err, desc)

	// prepare encodedNextKey
	encodedNextKey, err := codeutil.EncodeNextKeys(domain.DecodedNextKeys{{Field: "Date"}, {Field: "ID"}}, transactionList[3])
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 2, 1, 0)
	expDecodedNextKey := domain.DecodedNextKeys{
		{Field: "Date", Value: transactionList[3].Date.Format("2006-01-02")},
		{Field: "ID", Value: fmt.Sprint(transactionList[3].ID)},
	}

	opt := domain.GetTransOpt{
		Cursor: domain.Cursor{
			NextKey: encodedNextKey,
			Size:    3,
		},
		Sort: &domain.Sort{
			By:  domain.SortByTypeDate,
			Dir: domain.SortDirTypeDesc,
		},
	}
	trans, deencodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Equal(expDecodedNextKey, deencodedNextKey, desc)
}

func getAll_WithNextKeyCursorAndSortByPrice_ReturnCorrectData(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 300}
	ow2 := Transaction{Price: 300}
	ow3 := Transaction{Price: 500}
	ow4 := Transaction{Price: 1000}
	ow5 := Transaction{Price: 1000}
	ow6 := Transaction{Price: 2000}
	ow7 := Transaction{Price: 2500}
	ow8 := Transaction{Price: 3000}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 8, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8)
	s.Require().NoError(err, desc)

	// prepare encodedNextKey
	encodedNextKey, err := codeutil.EncodeNextKeys(domain.DecodedNextKeys{{Field: "Price"}, {Field: "ID"}}, transactionList[4])
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 3, 2, 1)
	expDecodedNextKey := domain.DecodedNextKeys{
		{Field: "Price", Value: fmt.Sprintf("%f", transactionList[4].Price)},
		{Field: "ID", Value: fmt.Sprint(transactionList[4].ID)},
	}

	opt := domain.GetTransOpt{
		Cursor: domain.Cursor{
			NextKey: encodedNextKey,
			Size:    3,
		},
		Sort: &domain.Sort{
			By:  domain.SortByTypePrice,
			Dir: domain.SortDirTypeDesc,
		},
	}
	trans, deencodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Equal(expDecodedNextKey, deencodedNextKey, desc)
}

func getAll_WithNextKeyCursorAndSortByTransType_ReturnCorrectData(s *TransactionSuite, desc string) {
	ow1 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow2 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue()}
	ow3 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow4 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow5 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue()}
	ow6 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow7 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue()}
	ow8 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue()}
	transactionList, user, mainList, subList, err := s.f.InsertTransactionsWithOneUser(mockCTX, 8, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8)
	s.Require().NoError(err, desc)

	// prepare encodedNextKey
	encodedNextKey, err := codeutil.EncodeNextKeys(domain.DecodedNextKeys{{Field: "Type"}, {Field: "ID"}}, transactionList[4])
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow3, ow4)
	s.Require().NoError(err, desc)

	expResult := GetAll_GenExpResult(transactionList, user, mainList, subList, 1, 6, 5)
	expDecodedNextKey := domain.DecodedNextKeys{
		{Field: "Type", Value: transactionList[4].Type},
		{Field: "ID", Value: fmt.Sprint(transactionList[4].ID)},
	}

	opt := domain.GetTransOpt{
		Cursor: domain.Cursor{
			NextKey: encodedNextKey,
			Size:    3,
		},
		Sort: &domain.Sort{
			By:  domain.SortByTypeTransType,
			Dir: domain.SortDirTypeDesc,
		},
	}
	trans, deencodedNextKey, err := s.repo.GetAll(mockCTX, opt, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
	s.Require().Equal(expDecodedNextKey, deencodedNextKey, desc)
}

func (s *TransactionSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, update successfully":       update_WithOneData_UpdateSuccessfully,
		"when with multiple data, update successfully":  update_WithMultipleData_UpdateSuccessfully,
		"when with multiple users, update successfully": update_WithMultipleUsers_UpdateSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_WithOneData_UpdateSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	t := domain.UpdateTransactionInput{
		ID:          transactions[0].ID,
		Type:        domain.CvtToTransactionType(transactions[0].Type),
		MainCategID: transactions[0].MainCategID,
		SubCategID:  transactions[0].SubCategID,
		Price:       999,
		Note:        "update note",
		Date:        mockTimeNow.AddDate(0, 0, -1),
	}

	err = s.repo.Update(mockCTX, t)
	s.Require().NoError(err, desc)

	var checkT Transaction
	stmt := "SELECT type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Require().Equal(t.Type.ToModelValue(), checkT.Type)
	s.Require().Equal(t.MainCategID, checkT.MainCategID)
	s.Require().Equal(t.SubCategID, checkT.SubCategID)
	s.Require().Equal(t.Price, checkT.Price)
	s.Require().Equal(t.Note, checkT.Note)
	s.Require().Equal(t.Date, checkT.Date)
}

func update_WithMultipleData_UpdateSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 2)
	s.Require().NoError(err, desc)

	t := domain.UpdateTransactionInput{
		ID:          transactions[0].ID,
		Type:        domain.CvtToTransactionType(transactions[0].Type),
		MainCategID: transactions[0].MainCategID,
		SubCategID:  transactions[0].SubCategID,
		Price:       999,
		Note:        "update note",
		Date:        mockTimeNow.AddDate(0, 0, -1),
	}

	err = s.repo.Update(mockCTX, t)
	s.Require().NoError(err, desc)

	var checkT Transaction
	stmt := "SELECT type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Require().Equal(t.Type.ToModelValue(), checkT.Type)
	s.Require().Equal(t.Price, checkT.Price)
	s.Require().Equal(t.Note, checkT.Note)
	s.Require().Equal(t.Date, checkT.Date)

	// check if other data is not updated
	var checkT2 Transaction
	err = s.db.QueryRow(stmt, transactions[1].ID).Scan(&checkT2.Type, &checkT2.MainCategID, &checkT2.SubCategID, &checkT2.Price, &checkT2.Note, &checkT2.Date)
	s.Require().NoError(err)
	s.Require().Equal(transactions[1].Type, checkT2.Type)
	s.Require().Equal(transactions[1].Price, checkT2.Price)
	s.Require().Equal(transactions[1].Note, checkT2.Note)
	s.Require().Equal(transactions[1].Date, checkT2.Date)
}

func update_WithMultipleUsers_UpdateSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	// prepare more users
	transactions2, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	t := domain.UpdateTransactionInput{
		ID:          transactions[0].ID,
		Type:        domain.CvtToTransactionType(transactions[0].Type),
		MainCategID: transactions[0].MainCategID,
		SubCategID:  transactions[0].SubCategID,
		Price:       999,
		Note:        "update note",
		Date:        mockTimeNow,
	}

	err = s.repo.Update(mockCTX, t)
	s.Require().NoError(err, desc)

	var checkT Transaction
	stmt := "SELECT type, main_category_id, sub_category_id, price, note, date FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.Type, &checkT.MainCategID, &checkT.SubCategID, &checkT.Price, &checkT.Note, &checkT.Date)
	s.Require().NoError(err)
	s.Require().Equal(t.Type.ToModelValue(), checkT.Type)
	s.Require().Equal(t.Price, checkT.Price)
	s.Require().Equal(t.Note, checkT.Note)
	s.Require().Equal(t.Date, checkT.Date)

	// check if other user's data is not updated
	var checkT2 Transaction
	err = s.db.QueryRow(stmt, transactions2[0].ID).Scan(&checkT2.Type, &checkT2.MainCategID, &checkT2.SubCategID, &checkT2.Price, &checkT2.Note, &checkT2.Date)
	s.Require().NoError(err)
	s.Require().Equal(transactions2[0].Type, checkT2.Type)
	s.Require().Equal(transactions2[0].Price, checkT2.Price)
	s.Require().Equal(transactions2[0].Note, checkT2.Note)
	s.Require().Equal(transactions2[0].Date, checkT2.Date)
}

func (s *TransactionSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, delete successfully":       delete_WithOneData_DeleteSuccessfully,
		"when with multiple data, delete successfully":  delete_WithMultipleData_DeleteSuccessfully,
		"when with multiple users, delete successfully": delete_WithMultipleUsers_DeleteSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func delete_WithOneData_DeleteSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	err = s.repo.Delete(mockCTX, transactions[0].ID)
	s.Require().NoError(err, desc)

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	s.Require().NoError(err, desc)
	s.Require().Equal(0, count, desc)

	// check if data exists
	var checkT Transaction
	stmt := "SELECT id FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.ID)
	s.Require().ErrorIs(err, sql.ErrNoRows, desc)
}

func delete_WithMultipleData_DeleteSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 3)
	s.Require().NoError(err, desc)

	err = s.repo.Delete(mockCTX, transactions[0].ID)
	s.Require().NoError(err, desc)

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	s.Require().NoError(err, desc)
	s.Require().Equal(2, count, desc)

	// check if data exists
	var checkT Transaction
	stmt := "SELECT id FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.ID)
	s.Require().ErrorIs(err, sql.ErrNoRows, desc)
}

func delete_WithMultipleUsers_DeleteSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 3)
	s.Require().NoError(err, desc)

	// prepare more users
	transactions2, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	err = s.repo.Delete(mockCTX, transactions[0].ID)
	s.Require().NoError(err, desc)

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	s.Require().NoError(err, desc)
	s.Require().Equal(3, count, desc)

	// check if data exists
	var checkT Transaction
	stmt := "SELECT id FROM transactions WHERE id = ?"
	err = s.db.QueryRow(stmt, transactions[0].ID).Scan(&checkT.ID)
	s.Require().ErrorIs(err, sql.ErrNoRows, desc)

	// check if other user's data still exists
	var countUser2 int
	stmt = "SELECT COUNT(*) FROM transactions WHERE user_id = ?"
	err = s.db.QueryRow(stmt, transactions2[0].UserID).Scan(&countUser2)
	s.Require().NoError(err, desc)
	s.Require().Equal(1, countUser2, desc)
}

func (s *TransactionSuite) TestGetAccInfo() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return successfully":                                  getAccInfo_NoError_ReturnSuccessfully,
		"when with multiple users, return data only with one user":            getAccInfo_WithMultipleUsers_ReturnDataOnlyWithOneUser,
		"when with many transaction, return correct calculation":              getAccInfo_WithManyTransaction_ReturnCorrectCalculation,
		"when query start date, return accumulated data after start date":     getAccInfo_QueryStartDate_ReturnDataAfterStartDate,
		"when query end date, return accumulated data before end date":        getAccInfo_QueryEndDate_ReturnDataBeforeEndDate,
		"when query start and end date, return accumulated data between them": getAccInfo_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAccInfo_NoError_ReturnSuccessfully(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue()}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 999,
		TotalIncome:  0,
		TotalBalance: -999,
	}

	query := domain.GetAccInfoQuery{}
	accInfo, err := s.repo.GetAccInfo(mockCTX, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_WithMultipleUsers_ReturnDataOnlyWithOneUser(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue()}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow1)
	s.Require().NoError(err, desc)

	ow2 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue()}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow2)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 999,
		TotalIncome:  0,
		TotalBalance: -999,
	}

	query := domain.GetAccInfoQuery{}
	accInfo, err := s.repo.GetAccInfo(mockCTX, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_WithManyTransaction_ReturnCorrectCalculation(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue()}
	ow2 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue()}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue()}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 3, ow1, ow2, ow3)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 1999,
		TotalIncome:  1000,
		TotalBalance: -999,
	}

	query := domain.GetAccInfoQuery{}
	accInfo, err := s.repo.GetAccInfo(mockCTX, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryStartDate_ReturnDataAfterStartDate(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 0,
		TotalIncome:  3000,
		TotalBalance: 3000,
	}

	startDate := mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly)
	query := domain.GetAccInfoQuery{
		StartDate: &startDate,
	}
	accInfo, err := s.repo.GetAccInfo(mockCTX, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryEndDate_ReturnDataBeforeEndDate(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 1999,
		TotalIncome:  1000,
		TotalBalance: -999,
	}

	endDate := mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly)
	query := domain.GetAccInfoQuery{
		EndDate: &endDate,
	}
	accInfo, err := s.repo.GetAccInfo(mockCTX, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func getAccInfo_QueryStartAndEndDate_ReturnDataBetweenStartAndEndDate(s *TransactionSuite, desc string) {
	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -3)}
	ow2 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -2)}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, -1)}
	ow4 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: mockTimeNow.AddDate(0, 0, 0)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 4, ow1, ow2, ow3, ow4)
	s.Require().NoError(err, desc)

	// prepare two more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := domain.AccInfo{
		TotalExpense: 1000,
		TotalIncome:  1000,
		TotalBalance: 0,
	}

	startDate := mockTimeNow.AddDate(0, 0, -2).Format(time.DateOnly)
	endDate := mockTimeNow.AddDate(0, 0, -1).Format(time.DateOnly)
	query := domain.GetAccInfoQuery{
		StartDate: &startDate,
		EndDate:   &endDate,
	}
	accInfo, err := s.repo.GetAccInfo(mockCTX, query, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, accInfo, desc)
}

func (s *TransactionSuite) TestGetByIDAndUserID() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when only one data, return successfully":       getByIDAndUserID_OnlyOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":  getByIDAndUserID_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getByIDAndUserID_WithMultipleUsers_ReturnSuccessfully,
		"when id not found, return error":               getByIDAndUserID_IDNotFound_ReturnError,
		"when user id not found, return error":          getByIDAndUserID_UserIDNotFound_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByIDAndUserID_OnlyOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := domain.Transaction{
		ID:     transactions[0].ID,
		UserID: transactions[0].UserID,
		Type:   domain.CvtToTransactionType(transactions[0].Type),
		Price:  transactions[0].Price,
		Note:   transactions[0].Note,
		Date:   transactions[0].Date,
	}

	trans, err := s.repo.GetByIDAndUserID(mockCTX, transactions[0].ID, transactions[0].UserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, desc)
}

func getByIDAndUserID_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 3)
	s.Require().NoError(err, desc)

	expResult := domain.Transaction{
		ID:     transactions[0].ID,
		UserID: transactions[0].UserID,
		Type:   domain.CvtToTransactionType(transactions[0].Type),
		Price:  transactions[0].Price,
		Note:   transactions[0].Note,
		Date:   transactions[0].Date,
	}

	trans, err := s.repo.GetByIDAndUserID(mockCTX, transactions[0].ID, transactions[0].UserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, err)
}

func getByIDAndUserID_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 3)
	s.Require().NoError(err, desc)

	// prepare more users
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	expResult := domain.Transaction{
		ID:     transactions[0].ID,
		UserID: transactions[0].UserID,
		Type:   domain.CvtToTransactionType(transactions[0].Type),
		Price:  transactions[0].Price,
		Note:   transactions[0].Note,
		Date:   transactions[0].Date,
	}

	trans, err := s.repo.GetByIDAndUserID(mockCTX, transactions[0].ID, transactions[0].UserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, trans, err)
}

func getByIDAndUserID_IDNotFound_ReturnError(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	_, err = s.repo.GetByIDAndUserID(mockCTX, transactions[0].ID+1, transactions[0].UserID)
	s.Require().Error(err, desc)
}

func getByIDAndUserID_UserIDNotFound_ReturnError(s *TransactionSuite, desc string) {
	transactions, _, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1)
	s.Require().NoError(err, desc)

	_, err = s.repo.GetByIDAndUserID(mockCTX, transactions[0].ID, transactions[0].UserID+1)
	s.Require().Error(err, desc)
}

func (s *TransactionSuite) TestGetDailyBarChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, return successfully":                       getDailyBarChartData_WithOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":                  getDailyBarChartData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully":                 getDailyBarChartData_WithMultipleUsers_ReturnSuccessfully,
		"when no main category ids, do not filter by main category ids": getDailyBarChartData_NoMainCategIDs_DoNotFilterByMainCategIDs,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getDailyBarChartData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)

	ow1 := Transaction{
		Price: 999,
		Type:  domain.TransactionTypeExpense.ToModelValue(),
		Date:  start,
	}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow1)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03-17": 999,
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetDailyBarChartData(mockCTX, dataRange, transactionType, nil, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getDailyBarChartData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-21")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[0].ID}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[1].ID}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1), MainCategID: mainCategList[2].ID}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 1), MainCategID: mainCategList[3].ID}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 2), MainCategID: mainCategList[3].ID}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 2), MainCategID: mainCategList[0].ID}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3), MainCategID: mainCategList[1].ID}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3), MainCategID: mainCategList[2].ID}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 4), MainCategID: mainCategList[3].ID}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4), MainCategID: mainCategList[2].ID}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5), MainCategID: mainCategList[0].ID}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5), MainCategID: mainCategList[1].ID}
	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 10, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03-17": 1,
		"2024-03-18": 1000,
		"2024-03-20": 1001,
		"2024-03-21": 2000,
	}

	// only get data with mainCategList[1] and mainCategList[2]
	mainCategIDs := []int64{mainCategList[1].ID, mainCategList[2].ID}
	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetDailyBarChartData(mockCTX, dataRange, transactionType, mainCategIDs, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getDailyBarChartData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-21")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[0].ID}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[1].ID}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1), MainCategID: mainCategList[2].ID}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 1), MainCategID: mainCategList[3].ID}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 2), MainCategID: mainCategList[3].ID}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 2), MainCategID: mainCategList[0].ID}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3), MainCategID: mainCategList[1].ID}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3), MainCategID: mainCategList[2].ID}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 4), MainCategID: mainCategList[3].ID}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4), MainCategID: mainCategList[2].ID}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5), MainCategID: mainCategList[0].ID}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5), MainCategID: mainCategList[1].ID}
	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 10, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	// prepare more users
	ow13 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow14 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 2)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow13, ow14)
	s.Require().NoError(err, desc)

	ow15 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3)}
	ow16 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow15, ow16)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03-17": 1,
		"2024-03-18": 1000,
		"2024-03-20": 1001,
		"2024-03-21": 2000,
	}

	// only get data with mainCategList[1] and mainCategList[2]
	mainCategIDs := []int64{mainCategList[1].ID, mainCategList[2].ID}
	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetDailyBarChartData(mockCTX, dataRange, transactionType, mainCategIDs, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getDailyBarChartData_NoMainCategIDs_DoNotFilterByMainCategIDs(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-21")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[0].ID}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[1].ID}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1), MainCategID: mainCategList[2].ID}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 1), MainCategID: mainCategList[3].ID}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 2), MainCategID: mainCategList[3].ID}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 2), MainCategID: mainCategList[0].ID}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3), MainCategID: mainCategList[1].ID}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3), MainCategID: mainCategList[2].ID}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 4), MainCategID: mainCategList[3].ID}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4), MainCategID: mainCategList[2].ID}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5), MainCategID: mainCategList[0].ID}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5), MainCategID: mainCategList[1].ID}
	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 10, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03-17": 1000,
		"2024-03-18": 1000,
		"2024-03-19": 999,
		"2024-03-20": 1001,
		"2024-03-21": 2000,
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetDailyBarChartData(mockCTX, dataRange, transactionType, nil, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func (s *TransactionSuite) TestGetMonthlyBarChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, return successfully":                      getMonthlyBarChartData_WithOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":                 getMonthlyBarChartData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully":                getMonthlyBarChartData_WithMultipleUsers_ReturnSuccessfully,
		"when no main category ids, do no filter by main category ids": getMonthlyBarChartData_NoMainCategIDs_DoNotFilterByMainCategIDs,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getMonthlyBarChartData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-20")
	s.Require().NoError(err, desc)

	ow1 := Transaction{
		Price: 999,
		Type:  domain.TransactionTypeExpense.ToModelValue(),
		Date:  start,
	}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow1)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03": 999,
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetMonthlyBarChartData(mockCTX, dataRange, transactionType, nil, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getMonthlyBarChartData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-08-30")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[0].ID}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[1].ID}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0), MainCategID: mainCategList[2].ID}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 1, 0), MainCategID: mainCategList[3].ID}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 3, 0), MainCategID: mainCategList[3].ID}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 3, 0), MainCategID: mainCategList[0].ID}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0), MainCategID: mainCategList[1].ID}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0), MainCategID: mainCategList[2].ID}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 5, 0), MainCategID: mainCategList[3].ID}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 5, 0), MainCategID: mainCategList[2].ID}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0), MainCategID: mainCategList[0].ID}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0), MainCategID: mainCategList[1].ID}
	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 12, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03": 1,
		"2024-04": 1000,
		"2024-07": 1001,
		"2024-08": 2000,
	}

	// only get data with mainCategList[1] and mainCategList[2]
	mainCategIDs := []int64{mainCategList[1].ID, mainCategList[2].ID}
	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetMonthlyBarChartData(mockCTX, dataRange, transactionType, mainCategIDs, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getMonthlyBarChartData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-08-30")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[0].ID}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[1].ID}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0), MainCategID: mainCategList[2].ID}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 1, 0), MainCategID: mainCategList[3].ID}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 3, 0), MainCategID: mainCategList[3].ID}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 3, 0), MainCategID: mainCategList[0].ID}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0), MainCategID: mainCategList[1].ID}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0), MainCategID: mainCategList[2].ID}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 5, 0), MainCategID: mainCategList[3].ID}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 5, 0), MainCategID: mainCategList[2].ID}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0), MainCategID: mainCategList[0].ID}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0), MainCategID: mainCategList[1].ID}
	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 12, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03": 1,
		"2024-04": 1000,
		"2024-07": 1001,
		"2024-08": 2000,
	}

	// prepare more users
	ow13 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow14 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow13, ow14)
	s.Require().NoError(err, desc)

	ow15 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	ow16 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow15, ow16)
	s.Require().NoError(err, desc)

	// only get data with mainCategList[1] and mainCategList[2]
	mainCategIDs := []int64{mainCategList[1].ID, mainCategList[2].ID}
	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetMonthlyBarChartData(mockCTX, dataRange, transactionType, mainCategIDs, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getMonthlyBarChartData_NoMainCategIDs_DoNotFilterByMainCategIDs(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-08-30")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[0].ID}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0), MainCategID: mainCategList[1].ID}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0), MainCategID: mainCategList[2].ID}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 1, 0), MainCategID: mainCategList[3].ID}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 3, 0), MainCategID: mainCategList[3].ID}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 3, 0), MainCategID: mainCategList[0].ID}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0), MainCategID: mainCategList[1].ID}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0), MainCategID: mainCategList[2].ID}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 5, 0), MainCategID: mainCategList[3].ID}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 5, 0), MainCategID: mainCategList[2].ID}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0), MainCategID: mainCategList[0].ID}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0), MainCategID: mainCategList[1].ID}
	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 12, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03": 1000,
		"2024-04": 1000,
		"2024-06": 999,
		"2024-07": 1001,
		"2024-08": 2000,
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetMonthlyBarChartData(mockCTX, dataRange, transactionType, nil, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func (s *TransactionSuite) TestGetPieChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when pie chart with one data, return successfully":      getPieChartData_WithOneData_ReturnSuccessfully,
		"when pie chart with multiple data, return successfully": getPieChartData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully":          getPieChartData_WithMultipleUsers_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getPieChartData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-19")
	s.Require().NoError(err, desc)

	ow1 := Transaction{
		Price: 999,
		Type:  domain.TransactionTypeExpense.ToModelValue(),
		Date:  start.AddDate(0, 0, 1),
	}
	_, user, mainCategs, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow1)
	s.Require().NoError(err, desc)

	expResult := domain.ChartData{
		Labels:   []string{mainCategs[0].Name},
		Datasets: []float64{999},
	}

	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	transactionType := domain.TransactionTypeExpense

	chartData, err := s.repo.GetPieChartData(mockCTX, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getPieChartData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-21")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: start}
	ow2 := Transaction{Price: 1, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: start}
	ow3 := Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: start}
	ow4 := Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: start}
	ow5 := Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: start}
	ow6 := Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: start}

	// income typ data
	ow7 := Transaction{Price: 999, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: start}
	ow8 := Transaction{Price: 1, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: start}

	// data out of date range
	ow9 := Transaction{Price: 1000, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: start.AddDate(0, 0, 10)}
	ow10 := Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: start.AddDate(0, 0, 10)}

	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 10, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	expResult := domain.ChartData{
		Labels:   []string{mainCategList[0].Name, mainCategList[1].Name, mainCategList[2].Name},
		Datasets: []float64{1000, 2000, 4000},
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetPieChartData(mockCTX, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getPieChartData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-21")
	s.Require().NoError(err, desc)

	mainCategOW1 := maincateg.MainCateg{Name: "food", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW2 := maincateg.MainCateg{Name: "clothes", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW3 := maincateg.MainCateg{Name: "transportation", Type: domain.TransactionTypeExpense.ToModelValue()}
	mainCategOW4 := maincateg.MainCateg{Name: "salary", Type: domain.TransactionTypeIncome.ToModelValue()} // income type
	mainCategList, user, err := s.f.InsertMainCategList(mockCTX, 5, mainCategOW1, mainCategOW2, mainCategOW3, mainCategOW4)
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: start}
	ow2 := Transaction{Price: 1, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: start}
	ow3 := Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: start}
	ow4 := Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: start}
	ow5 := Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: start}
	ow6 := Transaction{Price: 2000, Type: mainCategList[2].Type, MainCategID: mainCategList[2].ID, Date: start}

	// income typ data
	ow7 := Transaction{Price: 999, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: start}
	ow8 := Transaction{Price: 1, Type: mainCategList[3].Type, MainCategID: mainCategList[3].ID, Date: start}

	// data out of date range
	ow9 := Transaction{Price: 1000, Type: mainCategList[0].Type, MainCategID: mainCategList[0].ID, Date: start.AddDate(0, 0, 10)}
	ow10 := Transaction{Price: 1000, Type: mainCategList[1].Type, MainCategID: mainCategList[1].ID, Date: start.AddDate(0, 0, 10)}

	_, _, err = s.f.InsertTransactionWithGivenUser(mockCTX, 10, user, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	// prepare more user
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 2)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow11, ow12)
	s.Require().NoError(err, desc)

	ow13 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3)}
	ow14 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow13, ow14)
	s.Require().NoError(err, desc)

	expResult := domain.ChartData{
		Labels:   []string{mainCategList[0].Name, mainCategList[1].Name, mainCategList[2].Name},
		Datasets: []float64{1000, 2000, 4000},
	}

	transactionType := domain.TransactionTypeExpense
	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetPieChartData(mockCTX, dataRange, transactionType, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func (s *TransactionSuite) TestGetDailyLineChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with two data, return successfully":       getDailyLineChartData_WithTwoData_ReturnSuccessFully,
		"when with multiple data, return successfully":  getDailyLineChartData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getDailyLineChartData_WithMultipleUsers_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getDailyLineChartData_WithTwoData_ReturnSuccessFully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-21")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow2 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03-18": 1,
	}

	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	chartData, err := s.repo.GetDailyLineChartData(mockCTX, dataRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getDailyLineChartData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-06")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 3)}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3)}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4)}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4)}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 5)}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5)}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 6)}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 6)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 12, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03-01": -1000,
		"2024-03-02": -1000,
		"2024-03-04": 1,
		"2024-03-05": -1000,
		"2024-03-06": -2000,
	}

	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetDailyLineChartData(mockCTX, dataRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getDailyLineChartData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-03-06")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 3)}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 3)}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4)}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 4)}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 0, 5)}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 5)}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 6)}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 6)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 12, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	// prepare more users
	ow13 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow14 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow13, ow14)
	s.Require().NoError(err, desc)

	ow15 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	ow16 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 1)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow15, ow16)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03-01": -1000,
		"2024-03-02": -1000,
		"2024-03-04": 1,
		"2024-03-05": -1000,
		"2024-03-06": -2000,
	}

	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetDailyLineChartData(mockCTX, dataRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func (s *TransactionSuite) TestGetMonthlyLineChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with two data, return successfully":       getMonthlyLineChartData_WithTwoData_ReturnSuccessFully,
		"when with multiple data, return successfully":  getMonthlyLineChartData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getMonthlyLineChartData_WithMultipleUsers_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getMonthlyLineChartData_WithTwoData_ReturnSuccessFully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-06-21")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	ow2 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-04": 1,
	}

	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	chartData, err := s.repo.GetMonthlyLineChartData(mockCTX, dataRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getMonthlyLineChartData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-08-01")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 3, 0)}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 3, 0)}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0)}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0)}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 5, 0)}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 5, 0)}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0)}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 12, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03": -1000,
		"2024-04": -1000,
		"2024-06": 1,
		"2024-07": -1000,
		"2024-08": -2000,
	}

	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetMonthlyLineChartData(mockCTX, dataRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func getMonthlyLineChartData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	end, err := time.Parse(time.DateOnly, "2024-08-01")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow2 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow3 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	ow4 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	ow5 := Transaction{Price: 2000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 3, 0)}
	ow6 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 3, 0)}
	ow7 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0)}
	ow8 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 4, 0)}
	ow9 := Transaction{Price: 1000, Type: domain.TransactionTypeIncome.ToModelValue(), Date: start.AddDate(0, 5, 0)}
	ow10 := Transaction{Price: 2000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 5, 0)}
	ow11 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0)}
	ow12 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 6, 0)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 12, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10, ow11, ow12)
	s.Require().NoError(err, desc)

	// prepare more users
	ow13 := Transaction{Price: 999, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	ow14 := Transaction{Price: 1, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 0, 0)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow13, ow14)
	s.Require().NoError(err, desc)

	ow15 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	ow16 := Transaction{Price: 1000, Type: domain.TransactionTypeExpense.ToModelValue(), Date: start.AddDate(0, 1, 0)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow15, ow16)
	s.Require().NoError(err, desc)

	expResult := domain.DateToChartData{
		"2024-03": -1000,
		"2024-04": -1000,
		"2024-06": 1,
		"2024-07": -1000,
		"2024-08": -2000,
	}

	dataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}
	chartData, err := s.repo.GetMonthlyLineChartData(mockCTX, dataRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, chartData, desc)
}

func (s *TransactionSuite) TestGetMonthlyData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, return successfully":       getMonthlyData_WithOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":  getMonthlyData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getMonthlyData_WithMultipleUsers_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getMonthlyData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 3)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow1)
	s.Require().NoError(err, desc)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expResult := domain.MonthDayToTransactionType{
		4: domain.TransactionTypeExpense, // startDate.AddDate(0, 0, 3)
	}

	monthlyData, err := s.repo.GetMonthlyData(mockCTX, dateRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyData, desc)
}

func getMonthlyData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow2 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow3 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow4 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow5 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow6 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow7 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow8 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow9 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	ow10 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 10, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expResult := domain.MonthDayToTransactionType{
		2: domain.TransactionTypeExpense, // startDate.AddDate(0, 0, 1)
		5: domain.TransactionTypeIncome,  // startDate.AddDate(0, 0, 4)
		6: domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 5)
		7: domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 6)
	}

	monthlyData, err := s.repo.GetMonthlyData(mockCTX, dateRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyData, desc)
}

func getMonthlyData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err, desc)
	endDate, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow2 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	ow3 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow4 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 4)}
	ow5 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow6 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 5)}
	ow7 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow8 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 6)}
	ow9 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	ow10 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 40)}
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 10, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8, ow9, ow10)
	s.Require().NoError(err, desc)

	// prepare more users
	ow11 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 1)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow11)
	s.Require().NoError(err, desc)
	ow12 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: startDate.AddDate(0, 0, 2)}
	_, _, _, _, err = s.f.InsertTransactionsWithOneUser(mockCTX, 1, ow12)
	s.Require().NoError(err, desc)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	expResult := domain.MonthDayToTransactionType{
		2: domain.TransactionTypeExpense, // startDate.AddDate(0, 0, 1)
		5: domain.TransactionTypeIncome,  // startDate.AddDate(0, 0, 4)
		6: domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 5)
		7: domain.TransactionTypeBoth,    // startDate.AddDate(0, 0, 6)
	}

	monthlyData, err := s.repo.GetMonthlyData(mockCTX, dateRange, user.ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyData, desc)
}

func (s *TransactionSuite) TestGetMonthlyAggregatedData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when with one data, return successfully":       getMonthlyAggregatedData_WithOneData_ReturnSuccessfully,
		"when with multiple data, return successfully":  getMonthlyAggregatedData_WithMultipleData_ReturnSuccessfully,
		"when with multiple users, return successfully": getMonthlyAggregatedData_WithMultipleUsers_ReturnSuccessfully,
		"when with no data, return successfully":        getMonthlyAggregatedData_WithNoData_ReturnSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getMonthlyAggregatedData_WithOneData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	date, err := time.Parse(time.DateOnly, "2024-10-01")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date, Price: 1000}
	ow2 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 1, 0), Price: 1000} // out of date range
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 2, ow1, ow2)
	s.Require().NoError(err, desc)

	expResult := []domain.MonthlyAggregatedData{
		{
			UserID:       user.ID,
			TotalExpense: 1000,
			TotalIncome:  0,
		},
	}

	monthlyDataList, err := s.repo.GetMonthlyAggregatedData(mockCTX, date)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyDataList, desc)
}

func getMonthlyAggregatedData_WithMultipleData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	date, err := time.Parse(time.DateOnly, "2024-10-01")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date, Price: 100}
	ow2 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 1), Price: 200}
	ow3 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 2), Price: 300}
	ow4 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 3), Price: 400}
	ow5 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 4), Price: 500}
	ow6 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 5), Price: 600}
	ow7 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 1, 6), Price: 700}  // out of date range
	ow8 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, -1, 6), Price: 700} // out of date range
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 7, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8)
	s.Require().NoError(err, desc)

	expResult := []domain.MonthlyAggregatedData{
		{
			UserID:       user.ID,
			TotalExpense: 600,
			TotalIncome:  1500,
		},
	}

	monthlyDataList, err := s.repo.GetMonthlyAggregatedData(mockCTX, date)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyDataList, desc)
}

func getMonthlyAggregatedData_WithMultipleUsers_ReturnSuccessfully(s *TransactionSuite, desc string) {
	date, err := time.Parse(time.DateOnly, "2024-10-01")
	s.Require().NoError(err, desc)

	ow1 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date, Price: 100}
	ow2 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 1), Price: 200}
	ow3 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 2), Price: 300}
	ow4 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 3), Price: 400}
	ow5 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 4), Price: 500}
	ow6 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 5), Price: 600}
	ow7 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 1, 6), Price: 700}  // out of date range
	ow8 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, -1, 6), Price: 700} // out of date range
	_, user, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 8, ow1, ow2, ow3, ow4, ow5, ow6, ow7, ow8)
	s.Require().NoError(err, desc)

	ow9 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 1), Price: 100}
	ow10 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 0, 2), Price: 100}
	ow11 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 3), Price: 200}
	ow12 := Transaction{Type: domain.TransactionTypeIncome.ToModelValue(), Date: date.AddDate(0, 0, 4), Price: 300}
	ow13 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, 1, 5), Price: 200}
	ow14 := Transaction{Type: domain.TransactionTypeExpense.ToModelValue(), Date: date.AddDate(0, -1, 6), Price: 200}
	_, user2, _, _, err := s.f.InsertTransactionsWithOneUser(mockCTX, 6, ow9, ow10, ow11, ow12, ow13, ow14)
	s.Require().NoError(err, desc)

	expResult := []domain.MonthlyAggregatedData{
		{
			UserID:       user.ID,
			TotalExpense: 600,
			TotalIncome:  1500,
		},
		{
			UserID:       user2.ID,
			TotalExpense: 200,
			TotalIncome:  500,
		},
	}

	monthlyDataList, err := s.repo.GetMonthlyAggregatedData(mockCTX, date)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, monthlyDataList, desc)

}

func getMonthlyAggregatedData_WithNoData_ReturnSuccessfully(s *TransactionSuite, desc string) {
	data, err := time.Parse(time.DateOnly, "2024-10-01")
	s.Require().NoError(err, desc)

	monthlyDataList, err := s.repo.GetMonthlyAggregatedData(mockCTX, data)
	s.Require().NoError(err, desc)
	s.Require().Empty(monthlyDataList, desc)
}
