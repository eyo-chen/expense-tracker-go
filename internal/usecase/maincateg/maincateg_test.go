package maincateg

import (
	"context"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx = context.Background()
)

type MainCategSuite struct {
	suite.Suite
	mainCategUC       interfaces.MainCategUC
	mockIconRepo      *mocks.IconRepo
	mockMainCategRepo *mocks.MainCategRepo
}

func TestMainCategSuite(t *testing.T) {
	suite.Run(t, new(MainCategSuite))
}

func (s *MainCategSuite) SetupTest() {
	s.mockIconRepo = mocks.NewIconRepo(s.T())
	s.mockMainCategRepo = mocks.NewMainCategRepo(s.T())
	s.mainCategUC = NewMainCategUC(s.mockMainCategRepo, s.mockIconRepo)
}

func (s *MainCategSuite) TearDownTest() {
	s.mockIconRepo.AssertExpectations(s.T())
	s.mockMainCategRepo.AssertExpectations(s.T())
}

func (s *MainCategSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, create successfully": create_NoError_CreateSuccessfully,
		"when icon not exist, return error":  create_IconNotExist_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_CreateSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := domain.MainCateg{
		Icon: domain.Icon{ID: 1},
		Name: "Test",
	}

	// prepare mock service
	s.mockIconRepo.On("GetByID", mockCateg.Icon.ID).Return(domain.Icon{}, nil)
	s.mockMainCategRepo.On("Create", &mockCateg, mockUserID).Return(nil)

	// action, assertion
	err := s.mainCategUC.Create(mockCateg, mockUserID)
	s.Require().NoError(err, desc)
}

func create_IconNotExist_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := domain.MainCateg{
		Icon: domain.Icon{ID: 1},
		Name: "Test",
	}

	// prepare mock service
	s.mockIconRepo.On("GetByID", mockCateg.Icon.ID).Return(domain.Icon{}, domain.ErrIconNotFound)

	// action, assertion
	err := s.mainCategUC.Create(mockCateg, mockUserID)
	s.Require().EqualError(err, domain.ErrIconNotFound.Error(), desc)
}

func (s *MainCategSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, return main categories": getAll_NoError_ReturnMainCategories,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAll_NoError_ReturnMainCategories(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockMainCategs := []domain.MainCateg{
		{ID: 1, Name: "Test1"},
		{ID: 2, Name: "Test2"},
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetAll", mockCtx, mockUserID, domain.TransactionTypeExpense).Return(mockMainCategs, nil)

	// action, assertion
	res, err := s.mainCategUC.GetAll(mockCtx, mockUserID, domain.TransactionTypeExpense)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockMainCategs, res, desc)
}

func (s *MainCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, update successfully":         update_NoError_UpdateSuccessfully,
		"when main category not exist, return error": update_MainCategNotExist_ReturnError,
		"when icon not exist, return error":          update_IconNotExist_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoError_UpdateSuccessfully(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := domain.MainCateg{
		ID:   1,
		Icon: domain.Icon{ID: 1},
		Name: "Test",
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockCateg.ID, mockUserID).Return(&domain.MainCateg{}, nil)
	s.mockIconRepo.On("GetByID", mockCateg.Icon.ID).Return(domain.Icon{}, nil)
	s.mockMainCategRepo.On("Update", &mockCateg).Return(nil)

	// action, assertion
	err := s.mainCategUC.Update(mockCateg, mockUserID)
	s.Require().NoError(err, desc)
}

func update_MainCategNotExist_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := domain.MainCateg{
		ID:   1,
		Icon: domain.Icon{ID: 1},
		Name: "Test",
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockCateg.ID, mockUserID).Return(nil, domain.ErrMainCategNotFound)

	// action, assertion
	err := s.mainCategUC.Update(mockCateg, mockUserID)
	s.Require().EqualError(err, domain.ErrMainCategNotFound.Error(), desc)
}

func update_IconNotExist_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := domain.MainCateg{
		ID:   1,
		Icon: domain.Icon{ID: 1},
		Name: "Test",
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockCateg.ID, mockUserID).Return(&domain.MainCateg{}, nil)
	s.mockIconRepo.On("GetByID", mockCateg.Icon.ID).Return(domain.Icon{}, domain.ErrIconNotFound)

	// action, assertion
	err := s.mainCategUC.Update(mockCateg, mockUserID)
	s.Require().EqualError(err, domain.ErrIconNotFound.Error(), desc)
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
	// prepare mock data
	mockID := int64(1)

	// prepare mock service
	s.mockMainCategRepo.On("Delete", mockID).Return(nil)

	// action, assertion
	err := s.mainCategUC.Delete(mockID)
	s.Require().NoError(err, desc)
}
