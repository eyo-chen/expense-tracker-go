package transaction

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/codeutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx     = context.Background()
	mockLoc, _  = time.LoadLocation("")
	mockTimeNow = time.Unix(1629446406, 0).Truncate(24 * time.Hour).In(mockLoc)
)

type TransactionSuite struct {
	suite.Suite
	transactionUC       *TransactionUC
	mockTransactionRepo *mocks.TransactionRepo
	mockMainCategRepo   *mocks.MainCategRepo
	mockSubCategRepo    *mocks.SubCategRepo
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) SetupSuite() {
	logger.Register()
}

func (s *TransactionSuite) SetupTest() {
	s.mockTransactionRepo = mocks.NewTransactionRepo(s.T())
	s.mockMainCategRepo = mocks.NewMainCategRepo(s.T())
	s.mockSubCategRepo = mocks.NewSubCategRepo(s.T())
	s.transactionUC = NewTransactionUC(s.mockTransactionRepo, s.mockMainCategRepo, s.mockSubCategRepo)
}

func (s *TransactionSuite) TearDownTest() {
	s.mockTransactionRepo.AssertExpectations(s.T())
	s.mockMainCategRepo.AssertExpectations(s.T())
	s.mockSubCategRepo.AssertExpectations(s.T())
}

func (s *TransactionSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, create successfully":                                                      create_NoError_CreateSuccessfully,
		"when get main category fail, return error":                                               create_GetMainCategFail_ReturnError,
		"when type of main category not match transaction type, return error":                     create_TypeNotMatch_ReturnError,
		"when get sub category fail, return error":                                                create_GetSubCategFail_ReturnError,
		"when main category of sub category not match main category of transaction, return error": create_MainCategNotMatch_ReturnError,
		"when create fail, return error":                                                          create_CreateFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_CreateSuccessfully(s *TransactionSuite, desc string) {
	// prepare mock data
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}

	// prepare input
	transInput := domain.CreateTransactionInput{
		UserID:      1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	// prepare mock services
	s.mockMainCategRepo.Mock.On("GetByID", transInput.MainCategID, transInput.UserID).Return(&mainCateg, nil).Once()
	s.mockSubCategRepo.Mock.On("GetByID", transInput.SubCategID, transInput.UserID).Return(&subCateg, nil).Once()
	s.mockTransactionRepo.Mock.On("Create", mockCtx, transInput).Return(nil).Once()

	// action, assertion
	err := s.transactionUC.Create(mockCtx, transInput)
	s.Require().NoError(err, desc)
}

func create_GetMainCategFail_ReturnError(s *TransactionSuite, desc string) {
	// prepare input
	transInput := domain.CreateTransactionInput{
		UserID:      1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	// prepare mock services
	s.mockMainCategRepo.Mock.On("GetByID", transInput.MainCategID, transInput.UserID).Return(nil, errors.New("get main category fail")).Once()

	// action, assertion
	err := s.transactionUC.Create(mockCtx, transInput)
	s.Require().Equal(errors.New("get main category fail"), err, desc)
}

func create_TypeNotMatch_ReturnError(s *TransactionSuite, desc string) {
	// prepare mock data
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeIncome}

	// prepare input
	transInput := domain.CreateTransactionInput{
		UserID:      1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	// prepare mock services
	s.mockMainCategRepo.Mock.On("GetByID", transInput.MainCategID, transInput.UserID).Return(&mainCateg, nil).Once()

	// action, assertion
	err := s.transactionUC.Create(mockCtx, transInput)
	s.Require().EqualError(err, domain.ErrTypeNotConsistent.Error(), desc)
}

func create_GetSubCategFail_ReturnError(s *TransactionSuite, desc string) {
	// prepare mock data
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}

	// prepare input
	transInput := domain.CreateTransactionInput{
		UserID:      1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	// prepare mock services
	s.mockMainCategRepo.Mock.On("GetByID", transInput.MainCategID, transInput.UserID).Return(&mainCateg, nil).Once()
	s.mockSubCategRepo.Mock.On("GetByID", transInput.SubCategID, transInput.UserID).Return(nil, errors.New("get subcategory fail")).Once()

	// action, assertion
	err := s.transactionUC.Create(mockCtx, transInput)
	s.Require().EqualError(err, "get subcategory fail", desc)
}

func create_MainCategNotMatch_ReturnError(s *TransactionSuite, desc string) {
	// prepare mock data
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 2}

	// prepare input
	transInput := domain.CreateTransactionInput{
		UserID:      1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	// prepare mock services
	s.mockMainCategRepo.Mock.On("GetByID", transInput.MainCategID, transInput.UserID).Return(&mainCateg, nil).Once()
	s.mockSubCategRepo.Mock.On("GetByID", transInput.SubCategID, transInput.UserID).Return(&subCateg, nil).Once()

	// action, assertion
	err := s.transactionUC.Create(mockCtx, transInput)
	s.Require().EqualError(err, domain.ErrMainCategNotConsistent.Error(), desc)
}

func create_CreateFail_ReturnError(s *TransactionSuite, desc string) {
	// prepare mock data
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}

	// prepare input
	transInput := domain.CreateTransactionInput{
		UserID:      1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	// prepare mock services
	s.mockMainCategRepo.Mock.On("GetByID", transInput.MainCategID, transInput.UserID).Return(&mainCateg, nil).Once()
	s.mockSubCategRepo.Mock.On("GetByID", transInput.SubCategID, transInput.UserID).Return(&subCateg, nil).Once()
	s.mockTransactionRepo.Mock.On("Create", mockCtx, transInput).Return(errors.New("create fail")).Once()

	// action, assertion
	err := s.transactionUC.Create(mockCtx, transInput)
	s.Require().EqualError(err, "create fail", desc)
}

func (s *TransactionSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return transactions":                                        getAll_NoError_ReturnTransactions,
		"when get transactions fail, return error":                                  getAll_GetTransFail_ReturnError,
		"when it's the first page with size, return correct cursor":                 getAll_InitPageWithSize_ReturnCorrectCursor,
		"when it's the first page with size and sort, return correct cursor":        getAll_InitPageWithSizeAndSort_ReturnCorrectCursor,
		"when it's not the first page with decoded next key, return correct cursor": getAll_WithDecodedNextKey_ReturnCorrectCursor,
		"when size is empty value, return no cursor":                                getAll_SizeIsEmptyValue_ReturnNoCursor,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAll_NoError_ReturnTransactions(s *TransactionSuite, desc string) {
	mockDecodedNextKeys := domain.DecodedNextKeys{}
	mockOpt := domain.GetTransOpt{}
	mockUser := domain.User{ID: 1}
	mockTrans := []domain.Transaction{{ID: 1, UserID: 1}}

	s.mockTransactionRepo.On("GetAll", mockCtx, mockOpt, int64(1)).
		Return(mockTrans, mockDecodedNextKeys, nil).Once()

	result, cursor, err := s.transactionUC.GetAll(mockCtx, mockOpt, mockUser)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockTrans, result, desc)
	s.Require().Equal(domain.Cursor{}, cursor, desc)
}

func getAll_GetTransFail_ReturnError(s *TransactionSuite, desc string) {
	mockDecodedNextKeys := domain.DecodedNextKeys{}
	mockOpt := domain.GetTransOpt{}
	mockUser := domain.User{ID: 1}

	s.mockTransactionRepo.On("GetAll", mockCtx, mockOpt, int64(1)).
		Return(nil, mockDecodedNextKeys, errors.New("error")).Once()

	result, cursor, err := s.transactionUC.GetAll(mockCtx, mockOpt, mockUser)
	s.Require().Equal(errors.New("error"), err, desc)
	s.Require().Nil(result, desc)
	s.Require().Equal(domain.Cursor{}, cursor, desc)
}

func getAll_InitPageWithSize_ReturnCorrectCursor(s *TransactionSuite, desc string) {
	mockDecodedNextKeys := domain.DecodedNextKeys{}
	mockOpt := domain.GetTransOpt{Cursor: domain.Cursor{Size: 1}}
	mockUser := domain.User{ID: 1}
	mockTrans := []domain.Transaction{{ID: 1, UserID: 1}}

	s.mockTransactionRepo.On("GetAll", mockCtx, mockOpt, int64(1)).
		Return(mockTrans, mockDecodedNextKeys, nil).Once()

	result, cursor, err := s.transactionUC.GetAll(mockCtx, mockOpt, mockUser)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockTrans, result, desc)
	s.Require().Equal(1, cursor.Size, desc)

	// check decoded next key
	encodedNextKey, err := codeutil.DecodeNextKeys(cursor.NextKey, nil)
	s.Require().NoError(err, desc)
	s.Require().Equal(domain.DecodedNextKeys{{Field: "ID", Value: "1"}}, encodedNextKey, desc)
}

func getAll_InitPageWithSizeAndSort_ReturnCorrectCursor(s *TransactionSuite, desc string) {
	mockDecodedNextKeys := domain.DecodedNextKeys{}
	mockOpt := domain.GetTransOpt{
		Cursor: domain.Cursor{Size: 2},
		Sort:   &domain.Sort{By: domain.SortByTypePrice},
	}
	mockUser := domain.User{ID: 1}
	mockTrans := []domain.Transaction{{ID: 1, UserID: 1, Price: 100}, {ID: 2, UserID: 1, Price: 200.00}}

	expDecodedNextKeys := domain.DecodedNextKeys{
		{Field: "Price", Value: "200.000000"},
		{Field: "ID", Value: "2"},
	}

	s.mockTransactionRepo.On("GetAll", mockCtx, mockOpt, int64(1)).
		Return(mockTrans, mockDecodedNextKeys, nil).Once()

	result, cursor, err := s.transactionUC.GetAll(mockCtx, mockOpt, mockUser)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockTrans, result, desc)
	s.Require().Equal(2, cursor.Size, desc)

	// check decoded next key
	encodedNextKey, err := codeutil.DecodeNextKeys(cursor.NextKey, nil)
	s.Require().NoError(err, desc)
	s.Require().Equal(expDecodedNextKeys, encodedNextKey, desc)
}

func getAll_WithDecodedNextKey_ReturnCorrectCursor(s *TransactionSuite, desc string) {
	mockDecodedNextKeys := domain.DecodedNextKeys{{Field: "ID", Value: "1"}}
	mockOpt := domain.GetTransOpt{Cursor: domain.Cursor{Size: 1, NextKey: "eyJJRCI6IjEifQ=="}}
	mockUser := domain.User{ID: 1}
	mockTrans := []domain.Transaction{{ID: 2, UserID: 1}}

	s.mockTransactionRepo.On("GetAll", mockCtx, mockOpt, int64(1)).
		Return(mockTrans, mockDecodedNextKeys, nil).Once()

	result, cursor, err := s.transactionUC.GetAll(mockCtx, mockOpt, mockUser)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockTrans, result, desc)
	s.Require().Equal(1, cursor.Size, desc)

	// check encoded next key
	encodedNextKey, err := codeutil.DecodeNextKeys(cursor.NextKey, nil)
	s.Require().NoError(err, desc)
	s.Require().Equal(domain.DecodedNextKeys{{Field: "ID", Value: "2"}}, encodedNextKey, desc)
}

func getAll_SizeIsEmptyValue_ReturnNoCursor(s *TransactionSuite, desc string) {
	mockDecodedNextKeys := domain.DecodedNextKeys{}
	mockOpt := domain.GetTransOpt{Cursor: domain.Cursor{Size: 3}}
	mockUser := domain.User{ID: 1}
	mockTrans := []domain.Transaction{}

	s.mockTransactionRepo.On("GetAll", mockCtx, mockOpt, int64(1)).
		Return(mockTrans, mockDecodedNextKeys, nil).Once()

	result, cursor, err := s.transactionUC.GetAll(mockCtx, mockOpt, mockUser)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockTrans, result, desc)
	s.Require().Equal(domain.Cursor{}, cursor, desc)
}

func (s *TransactionSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, update successfully":                                                      update_NoError_UpdateSuccessfully,
		"when get main category fail, return error":                                               update_GetMainCategFail_ReturnError,
		"when type of main category not match transaction type, return error":                     update_TypeNotMatch_ReturnError,
		"when get sub category fail, return error":                                                update_GetSubCategFail_ReturnError,
		"when main category of sub category not match main category of transaction, return error": update_MainCategNotMatch_ReturnError,
		"when get transaction fail, return error":                                                 update_GetTransFail_UpdateSuccessfully,
		"when update fail, return error":                                                          update_UpdateFail_UpdateSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoError_UpdateSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCategRepo.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCategRepo.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	s.mockTransactionRepo.On("GetByIDAndUserID", mockCtx, trans.ID, user.ID).
		Return(domain.Transaction{}, nil).Once()

	s.mockTransactionRepo.On("Update", mockCtx, trans).
		Return(nil).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().NoError(err, desc)
}

func update_GetMainCategFail_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCategRepo.On("GetByID", trans.MainCategID, user.ID).
		Return(nil, errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func update_TypeNotMatch_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeIncome} // set type to income
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCategRepo.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(domain.ErrTypeNotConsistent, err, desc)
}

func update_GetSubCategFail_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCategRepo.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCategRepo.On("GetByID", trans.SubCategID, user.ID).
		Return(nil, errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func update_MainCategNotMatch_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 2}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCategRepo.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCategRepo.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(domain.ErrMainCategNotConsistent, err, desc)
}

func update_GetTransFail_UpdateSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCategRepo.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCategRepo.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	s.mockTransactionRepo.On("GetByIDAndUserID", mockCtx, trans.ID, user.ID).
		Return(domain.Transaction{}, errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func update_UpdateFail_UpdateSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{ID: 1}
	mainCateg := domain.MainCateg{ID: 1, Type: domain.TransactionTypeExpense}
	subCateg := domain.SubCateg{ID: 1, MainCategID: 1}
	trans := domain.UpdateTransactionInput{
		ID:          1,
		Type:        domain.TransactionTypeExpense,
		MainCategID: 1,
		SubCategID:  1,
		Price:       100,
		Date:        mockTimeNow,
		Note:        "note",
	}

	s.mockMainCategRepo.On("GetByID", trans.MainCategID, user.ID).
		Return(&mainCateg, nil).Once()

	s.mockSubCategRepo.On("GetByID", trans.SubCategID, user.ID).
		Return(&subCateg, nil).Once()

	s.mockTransactionRepo.On("GetByIDAndUserID", mockCtx, trans.ID, user.ID).
		Return(domain.Transaction{}, nil).Once()

	s.mockTransactionRepo.On("Update", mockCtx, trans).
		Return(errors.New("error")).Once()

	err := s.transactionUC.Update(mockCtx, trans, user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func (s *TransactionSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, delete successfully":       delete_NoError_DeleteSuccessfully,
		"when check permession fail, return error": delete_CheckPermessionFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func delete_NoError_DeleteSuccessfully(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	s.mockTransactionRepo.
		On("GetByIDAndUserID", mockCtx, int64(1), user.ID).
		Return(domain.Transaction{}, nil).Once()

	s.mockTransactionRepo.On("Delete", mockCtx, int64(1)).
		Return(nil).Once()

	err := s.transactionUC.Delete(mockCtx, int64(1), user)
	s.Require().NoError(err, desc)
}

func delete_CheckPermessionFail_ReturnError(s *TransactionSuite, desc string) {
	user := domain.User{
		ID: 1,
	}

	s.mockTransactionRepo.
		On("GetByIDAndUserID", mockCtx, int64(1), user.ID).
		Return(domain.Transaction{}, errors.New("error")).Once()

	err := s.transactionUC.Delete(mockCtx, int64(1), user)
	s.Require().Equal(errors.New("error"), err, desc)
}

func (s *TransactionSuite) TestGetAccInfo() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return acc info":       getAccInfo_NoError_ReturnAccInfo,
		"when get acc info fail, return error": getAccInfo_GetAccInfoFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAccInfo_NoError_ReturnAccInfo(s *TransactionSuite, desc string) {
	startDate := "2024-03-01"
	endDate := "2024-03-31"
	user := domain.User{ID: 1}
	query := domain.GetAccInfoQuery{StartDate: &startDate, EndDate: &endDate}
	accInfo := domain.AccInfo{
		TotalIncome:  100,
		TotalExpense: 200,
		TotalBalance: -100,
	}

	s.mockTransactionRepo.On("GetAccInfo", mockCtx, query, user.ID).
		Return(accInfo, nil).Once()

	result, err := s.transactionUC.GetAccInfo(mockCtx, query, user)
	s.Require().NoError(err, desc)
	s.Require().Equal(accInfo, result, desc)
}

func getAccInfo_GetAccInfoFail_ReturnError(s *TransactionSuite, desc string) {
	startDate := "2024-03-01"
	endDate := "2024-03-31"
	user := domain.User{ID: 1}
	query := domain.GetAccInfoQuery{StartDate: &startDate, EndDate: &endDate}

	s.mockTransactionRepo.On("GetAccInfo", mockCtx, query, user.ID).
		Return(domain.AccInfo{}, errors.New("get acc info fail")).Once()

	result, err := s.transactionUC.GetAccInfo(mockCtx, query, user)
	s.Require().EqualError(err, "get acc info fail", desc)
	s.Require().Equal(domain.AccInfo{}, result, desc)
}

func (s *TransactionSuite) TestGetBarChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when time range type is one week day, return week day data":         getBarChartData_WithOneWeekDay_ReturnWeekDayData,
		"when time range type is one week, return date data":                 getBarChartData_WithOneWeek_ReturnDateData,
		"when time range type is two weeks, return date data":                getBarChartData_WithTwoWeeks_ReturnDateData,
		"when time range type is one month, return date data":                getBarChartData_WithOneMonth_ReturnDateData,
		"when time ragne type is three months, return date accumulated data": getBarChartData_WithThreeMonths_ReturnDateData,
		"when time range type is six months, return month data":              getBarChartData_WithSixMonths_ReturnDateData,
		"when time range type is one year, return month data":                getBarChartData_WithOneYear_ReturnDateData,
		"when get chart data fail, return error":                             getBarChartData_GetChartDataFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getBarChartData_WithOneWeekDay_ReturnWeekDayData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": 200,
		"2024-03-21": 500,
		"2024-03-22": 600,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeekDay, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithOneWeek_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": 200,
		"2024-03-21": 500,
		"2024-03-22": 600,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeek, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithTwoWeeks_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-30")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": 200,
		"2024-03-21": 500,
		"2024-03-22": 600,
		"2024-03-23": 700,
		"2024-03-28": 800,
		"2024-03-29": 900,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23", "03/24", "03/25", "03/26", "03/27", "03/28", "03/29", "03/30"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 700, 0, 0, 0, 0, 800, 900, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeTwoWeeks, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithOneMonth_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-01": 100,
		"2024-03-02": 200,
		"2024-03-05": 500,
		"2024-03-06": 600,
		"2024-03-14": 700,
		"2024-03-19": 800,
		"2024-03-20": 900,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/01", "03/02", "03/03", "03/04", "03/05", "03/06", "03/07", "03/08", "03/09", "03/10", "03/11", "03/12", "03/13", "03/14", "03/15", "03/16", "03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23", "03/24", "03/25", "03/26", "03/27", "03/28", "03/29", "03/30", "03/31"},
		Datasets: []float64{100, 200, 0, 0, 500, 600, 0, 0, 0, 0, 0, 0, 0, 700, 0, 0, 0, 0, 800, 900, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneMonth, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithThreeMonths_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-05-31")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-01": 100,
		"2024-03-02": 200,
		"2024-03-05": 500,
		"2024-03-06": 600,
		"2024-03-14": 700,
		"2024-03-19": 800,
		"2024-03-20": 900,
		"2024-04-01": 1000,
		"2024-04-02": 1100,
		"2024-04-05": 1200,
		"2024-04-06": 1300,
		"2024-04-14": 1400,
		"2024-04-19": 1500,
		"2024-04-20": 1600,
		"2024-05-01": 1700,
		"2024-05-02": 1800,
		"2024-05-05": 1900,
		"2024-05-06": 2000,
		"2024-05-14": 2100,
		"2024-05-19": 2200,
		"2024-05-20": 2300,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/01", "03/04", "03/07", "03/10", "03/13", "03/16", "03/19", "03/22", "03/25", "03/28", "03/31", "04/03", "04/06", "04/09", "04/12", "04/15", "04/18", "04/21", "04/24", "04/27", "04/30", "05/03", "05/06", "05/09", "05/12", "05/15", "05/18", "05/21", "05/24", "05/27", "05/30"},
		Datasets: []float64{100, 200, 1100, 0, 0, 700, 800, 900, 0, 0, 0, 2100, 2500, 0, 0, 1400, 0, 3100, 0, 0, 0, 3500, 3900, 0, 0, 2100, 0, 4500, 0, 0, 0},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeThreeMonths, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithSixMonths_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-08-31")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03": 100,
		"2024-04": 500,
		"2024-06": 700,
		"2024-07": 900,
		"2024-08": 1100,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetMonthlyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"Mar", "Apr", "May", "Jun", "Jul", "Aug"},
		Datasets: []float64{100, 500, 0, 700, 900, 1100},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeSixMonths, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_WithOneYear_ReturnDateData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2025-02-28")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03": 100,
		"2024-04": 500,
		"2024-06": 700,
		"2024-07": 900,
		"2024-08": 1100,
		"2024-09": 1300,
		"2024-10": 1500,
		"2024-11": 1700,
		"2024-12": 1900,
		"2025-01": 2100,
		"2025-02": 2300,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetMonthlyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "Jan", "Feb"},
		Datasets: []float64{100, 500, 0, 700, 900, 1100, 1300, 1500, 1700, 1900, 2100, 2300},
	}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneYear, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getBarChartData_GetChartDataFail_ReturnError(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	mainCategIDs := []int64{1}

	s.mockTransactionRepo.On("GetDailyBarChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, mainCategIDs, int64(1)).
		Return(domain.DateToChartData{}, errors.New("error")).Once()

	// prepare expected result
	expResult := domain.ChartData{}

	result, err := s.transactionUC.GetBarChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeekDay, domain.TransactionTypeExpense, mainCategIDs, domain.User{ID: 1})
	s.Require().Equal(errors.New("error"), err, desc)
	s.Require().Equal(expResult, result, desc)
}

func (s *TransactionSuite) TestGetPieChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when no error, return chart data":       getPieChartData_NoError_ReturnChartData,
		"when get chart data fail, return error": getPieChartData_GetChartDataFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getPieChartData_NoError_ReturnChartData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	chartData := domain.ChartData{
		Labels:   []string{"label1", "label2"},
		Datasets: []float64{100, 200},
	}

	s.mockTransactionRepo.On("GetPieChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
		Return(chartData, nil).Once()

	result, err := s.transactionUC.GetPieChartData(mockCtx, chartDataRange, domain.TransactionTypeExpense, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(chartData, result, desc)
}

func getPieChartData_GetChartDataFail_ReturnError(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	s.mockTransactionRepo.On("GetPieChartData", mockCtx, chartDataRange, domain.TransactionTypeExpense, int64(1)).
		Return(domain.ChartData{}, errors.New("error")).Once()

	result, err := s.transactionUC.GetPieChartData(mockCtx, chartDataRange, domain.TransactionTypeExpense, domain.User{ID: 1})
	s.Require().EqualError(err, "error", desc)
	s.Require().Equal(domain.ChartData{}, result, desc)
}

func (s *TransactionSuite) TestGetLineChartData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when time range type is one week day, return week day data":         getLineChartData_WithOneWeekDay_ReturnWeekDayData,
		"when time range type is one week, return date data":                 getLineChartData_WithOneWeek_ReturnData,
		"when time range type is two weeks, return date data":                getLineChartData_WithTwoWeeks_ReturnData,
		"when time range type is one month, return date data":                getLineChartData_WithOneMonth_ReturnData,
		"when time ragne type is three months, return date accumulated data": getLineChartData_WithThreeMonths_ReturnData,
		"when time range type is six months, return month data":              getLineChartData_WithSixMonths_ReturnMonthData,
		"when time range type is one year, return month data":                getLineChartData_WithOneYear_ReturnMonthData,
		"when get chart data fail, return error":                             getLineChartData_GetChartDataFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getLineChartData_WithOneWeekDay_ReturnWeekDayData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": 200,
		"2024-03-21": -500,
		"2024-03-22": 600,
	}

	s.mockTransactionRepo.On("GetDailyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		Datasets: []float64{100, 200, 200, 200, -500, 600, 600},
	}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeekDay, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getLineChartData_WithOneWeek_ReturnData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": 100,
		"2024-03-18": -300,
		"2024-03-21": -500,
		"2024-03-22": 600,
	}

	s.mockTransactionRepo.On("GetDailyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23"},
		Datasets: []float64{100, -300, -300, -300, -500, 600, 600},
	}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeek, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getLineChartData_WithTwoWeeks_ReturnData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-30")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-17": -100,
		"2024-03-18": 200,
		"2024-03-21": 500,
		"2024-03-22": -600,
		"2024-03-23": 700,
		"2024-03-28": -800,
		"2024-03-29": 900,
	}

	s.mockTransactionRepo.On("GetDailyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23", "03/24", "03/25", "03/26", "03/27", "03/28", "03/29", "03/30"},
		Datasets: []float64{-100, 200, 200, 200, 500, -600, 700, 700, 700, 700, 700, -800, 900, 900},
	}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeek, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getLineChartData_WithOneMonth_ReturnData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-05")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-04-03")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-05": 500,
		"2024-03-06": -600,
		"2024-03-10": -100,
		"2024-03-14": 700,
		"2024-03-16": -500,
		"2024-03-19": -800,
		"2024-03-20": 900,
		"2024-03-25": 1000,
		"2024-03-31": -1400,
		"2024-04-02": 100,
	}

	s.mockTransactionRepo.On("GetDailyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/05", "03/06", "03/07", "03/08", "03/09", "03/10", "03/11", "03/12", "03/13", "03/14", "03/15", "03/16", "03/17", "03/18", "03/19", "03/20", "03/21", "03/22", "03/23", "03/24", "03/25", "03/26", "03/27", "03/28", "03/29", "03/30", "03/31", "04/01", "04/02", "04/03"},
		Datasets: []float64{500, -600, -600, -600, -600, -100, -100, -100, -100, 700, 700, -500, -500, -500, -800, 900, 900, 900, 900, 900, 1000, 1000, 1000, 1000, 1000, 1000, -1400, -1400, 100, 100},
	}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneMonth, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getLineChartData_WithThreeMonths_ReturnData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-05-31")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03-01": 100,
		"2024-03-02": 200,
		"2024-03-05": -500,
		"2024-03-06": -600,
		"2024-03-14": -700,
		"2024-03-19": 800,
		"2024-03-20": 900,
		"2024-04-01": 1000,
		"2024-04-02": -1100,
		"2024-04-05": -1200,
		"2024-04-06": 1300,
		"2024-04-14": -1400,
		"2024-04-19": 1500,
		"2024-04-20": 1600,
		"2024-05-01": -1700,
		"2024-05-02": -1800,
		"2024-05-05": 1900,
		"2024-05-06": 2000,
		"2024-05-14": -2100,
		"2024-05-19": 2200,
		"2024-05-20": -2300,
	}

	s.mockTransactionRepo.On("GetDailyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"03/01", "03/04", "03/07", "03/10", "03/13", "03/16", "03/19", "03/22", "03/25", "03/28", "03/31", "04/03", "04/06", "04/09", "04/12", "04/15", "04/18", "04/21", "04/24", "04/27", "04/30", "05/03", "05/06", "05/09", "05/12", "05/15", "05/18", "05/21", "05/24", "05/27", "05/30"},
		Datasets: []float64{100, 200, -600, -600, -600, -700, 800, 900, 900, 900, 900, -1100, 1300, 1300, 1300, -1400, -1400, 1600, 1600, 1600, 1600, -1800, 2000, 2000, 2000, -2100, -2100, -2300, -2300, -2300, -2300},
	}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeThreeMonths, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getLineChartData_WithSixMonths_ReturnMonthData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-08-31")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03": -100,
		"2024-04": -500,
		"2024-06": -700,
		"2024-07": -900,
		"2024-08": -1100,
	}

	s.mockTransactionRepo.On("GetMonthlyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"Mar", "Apr", "May", "Jun", "Jul", "Aug"},
		Datasets: []float64{-100, -500, -500, -700, -900, -1100},
	}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeSixMonths, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getLineChartData_WithOneYear_ReturnMonthData(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2025-02-28")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	DateToChartData := domain.DateToChartData{
		"2024-03": 100,
		"2024-04": 500,
		"2024-06": -700,
		"2024-07": 900,
		"2024-08": -1100,
		"2024-09": 1300,
		"2024-10": -1500,
		"2024-11": -1700,
		"2024-12": -1900,
		"2025-01": 2100,
		"2025-02": 2300,
	}

	s.mockTransactionRepo.On("GetMonthlyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(DateToChartData, nil).Once()

	// prepare expected result
	expResult := domain.ChartData{
		Labels:   []string{"Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "Jan", "Feb"},
		Datasets: []float64{100, 500, 500, -700, 900, -1100, 1300, -1500, -1700, -1900, 2100, 2300},
	}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneYear, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getLineChartData_GetChartDataFail_ReturnError(s *TransactionSuite, desc string) {
	start, err := time.Parse(time.DateOnly, "2024-03-17")
	s.Require().NoError(err)
	end, err := time.Parse(time.DateOnly, "2024-03-23")
	s.Require().NoError(err)

	chartDataRange := domain.ChartDateRange{
		Start: start,
		End:   end,
	}

	s.mockTransactionRepo.On("GetDailyLineChartData", mockCtx, chartDataRange, int64(1)).
		Return(domain.DateToChartData{}, errors.New("error")).Once()

	expResult := domain.ChartData{}

	result, err := s.transactionUC.GetLineChartData(mockCtx, chartDataRange, domain.TimeRangeTypeOneWeekDay, domain.User{ID: 1})
	s.Require().Equal(errors.New("error"), err, desc)
	s.Require().Equal(expResult, result, desc)
}

func (s *TransactionSuite) TestGetMonthlyData() {
	for scenario, fn := range map[string]func(s *TransactionSuite, desc string){
		"when 31 days in a month, return monthly data": getMonthlyData_31DaysInAMonth_ReturnMonthlyData,
		"when 30 days in a month, return monthly data": getMonthlyData_30DaysInAMonth_ReturnMonthlyData,
		"when 29 days in a month, return monthly data": getMonthlyData_29DaysInAMonth_ReturnMonthlyData,
		"when get monthly data fail, return error":     getMonthlyData_GetMonthlyDataFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getMonthlyData_31DaysInAMonth_ReturnMonthlyData(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-03-01")
	s.Require().NoError(err)
	endDate, err := time.Parse(time.DateOnly, "2024-03-31")
	s.Require().NoError(err)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	monthlyData := domain.MonthDayToTransactionType{
		5:  domain.TransactionTypeExpense,
		10: domain.TransactionTypeIncome,
		20: domain.TransactionTypeBoth,
	}

	expResult := []domain.TransactionType{
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeExpense, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeIncome, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeBoth, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified,
	}

	s.mockTransactionRepo.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
		Return(monthlyData, nil).Once()

	result, err := s.transactionUC.GetMonthlyData(mockCtx, dateRange, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getMonthlyData_30DaysInAMonth_ReturnMonthlyData(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-04-01")
	s.Require().NoError(err)
	endDate, err := time.Parse(time.DateOnly, "2024-04-30")
	s.Require().NoError(err)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	monthlyData := domain.MonthDayToTransactionType{
		3:  domain.TransactionTypeExpense,
		8:  domain.TransactionTypeIncome,
		10: domain.TransactionTypeBoth,
	}

	expResult := []domain.TransactionType{
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeExpense,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeIncome, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeBoth, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
	}

	s.mockTransactionRepo.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
		Return(monthlyData, nil).Once()

	result, err := s.transactionUC.GetMonthlyData(mockCtx, dateRange, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getMonthlyData_29DaysInAMonth_ReturnMonthlyData(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-02-01")
	s.Require().NoError(err)
	endDate, err := time.Parse(time.DateOnly, "2024-02-29")
	s.Require().NoError(err)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	monthlyData := domain.MonthDayToTransactionType{
		3:  domain.TransactionTypeExpense,
		8:  domain.TransactionTypeIncome,
		10: domain.TransactionTypeBoth,
	}

	expResult := []domain.TransactionType{
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeExpense,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeIncome, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeBoth, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
		domain.TransactionTypeUnSpecified, domain.TransactionTypeUnSpecified,
	}

	s.mockTransactionRepo.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
		Return(monthlyData, nil).Once()

	result, err := s.transactionUC.GetMonthlyData(mockCtx, dateRange, domain.User{ID: 1})
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, result, desc)
}

func getMonthlyData_GetMonthlyDataFail_ReturnError(s *TransactionSuite, desc string) {
	startDate, err := time.Parse(time.DateOnly, "2024-05-01")
	s.Require().NoError(err)
	endDate, err := time.Parse(time.DateOnly, "2024-05-31")
	s.Require().NoError(err)

	dateRange := domain.GetMonthlyDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	s.mockTransactionRepo.On("GetMonthlyData", mockCtx, dateRange, int64(1)).
		Return(nil, errors.New("error")).Once()

	result, err := s.transactionUC.GetMonthlyData(mockCtx, dateRange, domain.User{ID: 1})
	s.Require().Equal(errors.New("error"), err, desc)
	s.Require().Equal([]domain.TransactionType{}, result, desc)
}
