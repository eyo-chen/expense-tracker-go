package subcateg

import (
	"errors"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type SubCategSuite struct {
	suite.Suite
	subCategUC        *SubCategUC
	mockSubCategRepo  *mocks.SubCategRepo
	mockMainCategRepo *mocks.MainCategRepo
}

func TestSubCategSuite(t *testing.T) {
	suite.Run(t, new(SubCategSuite))
}

func (s *SubCategSuite) SetupTest() {
	s.mockSubCategRepo = mocks.NewSubCategRepo(s.T())
	s.mockMainCategRepo = mocks.NewMainCategRepo(s.T())
	s.subCategUC = NewSubCategUC(s.mockSubCategRepo, s.mockMainCategRepo)
}

func (s *SubCategSuite) TearDownTest() {
	s.mockSubCategRepo.AssertExpectations(s.T())
	s.mockMainCategRepo.AssertExpectations(s.T())
}

func (s *SubCategSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when no error, create successfully":                create_NoError_CreateSuccessfully,
		"when main category not exist, create successfully": create_MainCategNotExist_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_CreateSuccessfully(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := &domain.SubCateg{
		MainCategID: 1,
		Name:        "Test",
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockCateg.MainCategID, mockUserID).Return(&domain.MainCateg{}, nil)

	s.mockSubCategRepo.On("Create", mockCateg, mockUserID).Return(nil)

	// action, assertion
	err := s.subCategUC.Create(mockCateg, mockUserID)
	s.Require().NoError(err, desc)
}

func create_MainCategNotExist_ReturnError(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := &domain.SubCateg{
		MainCategID: 1,
		Name:        "Test",
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockCateg.MainCategID, mockUserID).Return(nil, domain.ErrMainCategNotFound)

	// action, assertion
	err := s.subCategUC.Create(mockCateg, mockUserID)
	s.Require().EqualError(err, domain.ErrMainCategNotFound.Error(), desc)
}

func (s *SubCategSuite) TestGetByMainCategID() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when no error, return data":                      getByMainCategID_NoError_ReturnData,
		"when get by main category ID fail, return error": getByMainCategID_GetByMainCategIDFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByMainCategID_NoError_ReturnData(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockMainCategID := int64(1)
	mockSubCategs := []*domain.SubCateg{
		{ID: 1, MainCategID: 1, Name: "Test 1"},
		{ID: 2, MainCategID: 1, Name: "Test 2"},
	}

	// prepare mock service
	s.mockSubCategRepo.On("GetByMainCategID", mockUserID, mockMainCategID).Return(mockSubCategs, nil)

	// action, assertion
	subCategs, err := s.subCategUC.GetByMainCategID(mockUserID, mockMainCategID)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockSubCategs, subCategs, desc)
}

func getByMainCategID_GetByMainCategIDFail_ReturnError(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockMainCategID := int64(1)

	// prepare mock service
	s.mockSubCategRepo.On("GetByMainCategID", mockUserID, mockMainCategID).Return(nil, errors.New("getByMainCategID error"))

	// action, assertion
	subCategs, err := s.subCategUC.GetByMainCategID(mockUserID, mockMainCategID)
	s.Require().EqualError(err, "getByMainCategID error", desc)
	s.Require().Nil(subCategs, desc)
}

func (s *SubCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when no error, update successfully":            update_NoError_UpdateSuccessfully,
		"when sub category not exist, return error":     update_SubCategNotExist_ReturnError,
		"when main category ID not match, return error": update_MainCategIDNotMatch_ReturnError,
		"when update fail, return error":                update_UpdateFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func update_NoError_UpdateSuccessfully(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInputCateg := &domain.SubCateg{
		ID:          1,
		MainCategID: 1,
		Name:        "Test",
	}
	mockCateg := &domain.SubCateg{
		ID:          1,
		MainCategID: 1,
		Name:        "Test",
	}

	// prepare mock service
	s.mockSubCategRepo.On("GetByID", mockInputCateg.ID, mockUserID).Return(mockCateg, nil)
	s.mockSubCategRepo.On("Update", mockInputCateg).Return(nil)

	// action, assertion
	err := s.subCategUC.Update(mockInputCateg, mockUserID)
	s.Require().NoError(err, desc)
}

func update_SubCategNotExist_ReturnError(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInputCateg := &domain.SubCateg{
		ID:          1,
		MainCategID: 1,
		Name:        "Test",
	}

	// prepare mock service
	s.mockSubCategRepo.On("GetByID", mockInputCateg.ID, mockUserID).Return(nil, domain.ErrSubCategNotFound)

	// action, assertion
	err := s.subCategUC.Update(mockInputCateg, mockUserID)
	s.Require().EqualError(err, domain.ErrSubCategNotFound.Error(), desc)
}

func update_MainCategIDNotMatch_ReturnError(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInputCateg := &domain.SubCateg{
		ID:          1,
		MainCategID: 1,
		Name:        "Test",
	}
	mockSubCateg := &domain.SubCateg{
		ID:          1,
		MainCategID: 2, // different main category ID
		Name:        "Test",
	}

	// prepare mock service
	s.mockSubCategRepo.On("GetByID", mockInputCateg.ID, mockUserID).Return(mockSubCateg, nil)

	// action, assertion
	err := s.subCategUC.Update(mockInputCateg, mockUserID)
	s.Require().EqualError(err, domain.ErrMainCategNotFound.Error(), desc)
}

func update_UpdateFail_ReturnError(s *SubCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInputCateg := &domain.SubCateg{
		ID:          1,
		MainCategID: 1,
		Name:        "Test",
	}
	mockCateg := &domain.SubCateg{
		ID:          1,
		MainCategID: 1,
		Name:        "Test",
	}

	// prepare mock service
	s.mockSubCategRepo.On("GetByID", mockInputCateg.ID, mockUserID).Return(mockCateg, nil)
	s.mockSubCategRepo.On("Update", mockInputCateg).Return(errors.New("update error"))

	// action, assertion
	err := s.subCategUC.Update(mockInputCateg, mockUserID)
	s.Require().EqualError(err, "update error", desc)
}

func (s *SubCategSuite) TestDelete() {
	for scenario, fn := range map[string]func(s *SubCategSuite, desc string){
		"when no error, delete successfully": delete_NoError_DeleteSuccessfully,
		"when delete fail, return error":     delete_DeleteFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func delete_NoError_DeleteSuccessfully(s *SubCategSuite, desc string) {
	// prepare mock data
	mockID := int64(1)

	// prepare mock service
	s.mockSubCategRepo.On("Delete", mockID).Return(nil)

	// action, assertion
	err := s.subCategUC.Delete(mockID)
	s.Require().NoError(err, desc)
}

func delete_DeleteFail_ReturnError(s *SubCategSuite, desc string) {
	// prepare mock data
	mockID := int64(1)

	// prepare mock service
	s.mockSubCategRepo.On("Delete", mockID).Return(errors.New("delete error"))

	// action, assertion
	err := s.subCategUC.Delete(mockID)
	s.Require().EqualError(err, "delete error", desc)
}
