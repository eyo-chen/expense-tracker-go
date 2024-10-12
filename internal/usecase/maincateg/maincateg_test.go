package maincateg

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx = context.Background()
)

type MainCategSuite struct {
	suite.Suite
	uc                *UC
	mockIconRepo      *mocks.IconRepo
	mockMainCategRepo *mocks.MainCategRepo
	mockUserIconRepo  *mocks.UserIconRepo
	mockRedisService  *mocks.RedisService
	mockS3Service     *mocks.S3Service
}

func TestMainCategSuite(t *testing.T) {
	suite.Run(t, new(MainCategSuite))
}

func (s *MainCategSuite) SetupTest() {
	s.mockIconRepo = mocks.NewIconRepo(s.T())
	s.mockMainCategRepo = mocks.NewMainCategRepo(s.T())
	s.mockUserIconRepo = mocks.NewUserIconRepo(s.T())
	s.mockRedisService = mocks.NewRedisService(s.T())
	s.mockS3Service = mocks.NewS3Service(s.T())
	s.uc = New(s.mockMainCategRepo, s.mockIconRepo, s.mockUserIconRepo, s.mockRedisService, s.mockS3Service)
}

func (s *MainCategSuite) TearDownTest() {
	s.mockIconRepo.AssertExpectations(s.T())
	s.mockMainCategRepo.AssertExpectations(s.T())
	s.mockUserIconRepo.AssertExpectations(s.T())
	s.mockRedisService.AssertExpectations(s.T())
	s.mockS3Service.AssertExpectations(s.T())
}

func (s *MainCategSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, create successfully":        create_NoError_CreateSuccessfully,
		"when icon type unspecified, return error":  create_IconTypeUnspecified_ReturnError,
		"when default icon not exist, return error": create_DefaultIconNotFound_ReturnError,
		"when user icon not exist, return error":    create_UserIconNotFound_ReturnError,
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
	mockIconData := "https://example.com/icon.png"
	mockInput := domain.CreateMainCategInput{
		Name:     "Test",
		Type:     domain.TransactionTypeExpense,
		IconType: domain.IconTypeDefault,
		IconID:   1,
	}
	mockCateg := domain.MainCateg{
		Name:     "Test",
		Type:     domain.TransactionTypeExpense,
		IconType: domain.IconTypeDefault,
		IconData: mockIconData,
	}
	mockDefaultIcon := domain.DefaultIcon{ID: 1, URL: mockIconData}

	// prepare mock service
	s.mockIconRepo.On("GetByID", mockCtx, mockInput.IconID).Return(mockDefaultIcon, nil)
	s.mockMainCategRepo.On("Create", mockCtx, mockCateg, mockUserID).Return(nil)

	// action, assertion
	err := s.uc.Create(mockCtx, mockInput, mockUserID)
	s.Require().NoError(err, desc)
}

func create_IconTypeUnspecified_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockCateg := domain.CreateMainCategInput{
		IconType: domain.IconTypeUnspecified,
	}

	// action, assertion
	err := s.uc.Create(mockCtx, mockCateg, mockUserID)
	s.Require().ErrorIs(err, domain.ErrIconNotFound, desc)
}

func create_DefaultIconNotFound_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInput := domain.CreateMainCategInput{
		IconType: domain.IconTypeDefault,
		IconID:   1,
		Name:     "Test",
	}

	// prepare mock service
	s.mockIconRepo.On("GetByID", mockCtx, mockInput.IconID).Return(domain.DefaultIcon{}, domain.ErrIconNotFound)

	// action, assertion
	err := s.uc.Create(mockCtx, mockInput, mockUserID)
	s.Require().ErrorIs(err, domain.ErrIconNotFound, desc)
}

func create_UserIconNotFound_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInput := domain.CreateMainCategInput{
		IconType: domain.IconTypeCustom,
		IconID:   1,
		Name:     "Test",
	}

	// prepare mock service
	s.mockUserIconRepo.On("GetByID", mockCtx, mockInput.IconID, mockUserID).Return(domain.UserIcon{}, domain.ErrUserIconNotFound)

	// action, assertion
	err := s.uc.Create(mockCtx, mockInput, mockUserID)
	s.Require().ErrorIs(err, domain.ErrUserIconNotFound, desc)
}

func (s *MainCategSuite) TestGetAll() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when with only default icon, return main categories": getAll_WithOnlyDefaultIcon_ReturnMainCategories,
		"when with only user icon, return main categories":    getAll_WithOnlyUserIcon_ReturnMainCategories,
		"when get object url fail, return error":              getAll_GetObjectURLFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getAll_WithOnlyDefaultIcon_ReturnMainCategories(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockMainCategs := []domain.MainCateg{
		{ID: 1, Name: "Test1", IconType: domain.IconTypeDefault, IconData: "https://example.com/icon1.png"},
		{ID: 2, Name: "Test2", IconType: domain.IconTypeDefault, IconData: "https://example.com/icon2.png"},
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetAll", mockCtx, mockUserID, domain.TransactionTypeExpense).Return(mockMainCategs, nil)

	// action, assertion
	res, err := s.uc.GetAll(mockCtx, mockUserID, domain.TransactionTypeExpense)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockMainCategs, res, desc)
}

func getAll_WithOnlyUserIcon_ReturnMainCategories(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockMainCategs := []domain.MainCateg{
		{ID: 1, Name: "Test1", IconType: domain.IconTypeCustom, IconData: "example.com/icon1.png"},
		{ID: 2, Name: "Test2", IconType: domain.IconTypeCustom, IconData: "example.com/icon2.png"},
	}
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockTTL := 7 * 24 * time.Hour
	mockKey1 := fmt.Sprintf("user_icon-%s", mockMainCategs[0].IconData)
	mockKey2 := fmt.Sprintf("user_icon-%s", mockMainCategs[1].IconData)

	// prepare mock service
	s.mockMainCategRepo.On("GetAll", mockCtx, mockUserID, domain.TransactionTypeExpense).Return(mockMainCategs, nil)
	s.mockRedisService.On("GetByFunc", mockCtx, mockKey1, mockTTL, mockGetFun).Return("https://example.com/icon1.png", nil)
	s.mockRedisService.On("GetByFunc", mockCtx, mockKey2, mockTTL, mockGetFun).Return("https://example.com/icon2.png", nil)

	// action, assertion
	res, err := s.uc.GetAll(mockCtx, mockUserID, domain.TransactionTypeExpense)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockMainCategs, res, desc)
}

func getAll_GetObjectURLFail_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockMainCategs := []domain.MainCateg{
		{ID: 1, Name: "Test1", IconType: domain.IconTypeCustom, IconData: "https://example.com/icon1.png"},
		{ID: 2, Name: "Test2", IconType: domain.IconTypeCustom, IconData: "https://example.com/icon2.png"},
	}
	mockKey := fmt.Sprintf("user_icon-%s", mockMainCategs[0].IconData)
	mockTTL := 7 * 24 * time.Hour
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockErr := errors.New("get object url failed")

	// prepare mock service
	s.mockMainCategRepo.On("GetAll", mockCtx, mockUserID, domain.TransactionTypeExpense).Return(mockMainCategs, nil)
	s.mockRedisService.On("GetByFunc", mockCtx, mockKey, mockTTL, mockGetFun).Return("", mockErr)

	// action, assertion
	res, err := s.uc.GetAll(mockCtx, mockUserID, domain.TransactionTypeExpense)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Empty(res, desc)
}

func (s *MainCategSuite) TestUpdate() {
	for scenario, fn := range map[string]func(s *MainCategSuite, desc string){
		"when no error, update successfully":         update_NoError_UpdateSuccessfully,
		"when main category not exist, return error": update_MainCategNotExist_ReturnError,
		"when icon type unspecified, return error":   update_IconTypeUnspecified_ReturnError,
		"when default icon not exist, return error":  update_DefaultIconNotFound_ReturnError,
		"when user icon not exist, return error":     update_UserIconNotFound_ReturnError,
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
	mockIconData := "https://example.com/icon.png"
	mockInput := domain.UpdateMainCategInput{
		ID:       1,
		Name:     "Test",
		IconType: domain.IconTypeDefault,
		IconID:   1,
	}
	mockCateg := domain.MainCateg{
		ID:       1,
		Name:     "Test",
		IconType: domain.IconTypeDefault,
		IconData: mockIconData,
	}
	mockDefaultIcon := domain.DefaultIcon{ID: 1, URL: mockIconData}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockInput.ID, mockUserID).Return(&domain.MainCateg{}, nil)
	s.mockIconRepo.On("GetByID", mockCtx, mockInput.IconID).Return(mockDefaultIcon, nil)
	s.mockMainCategRepo.On("Update", mockCtx, mockCateg).Return(nil)

	// action, assertion
	err := s.uc.Update(mockCtx, mockInput, mockUserID)
	s.Require().NoError(err, desc)
}

func update_MainCategNotExist_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInput := domain.UpdateMainCategInput{
		ID:       1,
		Name:     "Test",
		IconType: domain.IconTypeDefault,
		IconID:   1,
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockInput.ID, mockUserID).Return(nil, domain.ErrMainCategNotFound)

	// action, assertion
	err := s.uc.Update(mockCtx, mockInput, mockUserID)
	s.Require().EqualError(err, domain.ErrMainCategNotFound.Error(), desc)
}

func update_IconTypeUnspecified_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInput := domain.UpdateMainCategInput{
		IconType: domain.IconTypeUnspecified,
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockInput.ID, mockUserID).Return(&domain.MainCateg{}, nil)

	// action, assertion
	err := s.uc.Update(mockCtx, mockInput, mockUserID)
	s.Require().ErrorIs(err, domain.ErrIconNotFound, desc)
}

func update_DefaultIconNotFound_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInput := domain.UpdateMainCategInput{
		ID:       1,
		Name:     "Test",
		IconType: domain.IconTypeDefault,
		IconID:   1,
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockInput.ID, mockUserID).Return(&domain.MainCateg{}, nil)
	s.mockIconRepo.On("GetByID", mockCtx, mockInput.IconID).Return(domain.DefaultIcon{}, domain.ErrIconNotFound)

	// action, assertion
	err := s.uc.Update(mockCtx, mockInput, mockUserID)
	s.Require().ErrorIs(err, domain.ErrIconNotFound, desc)
}

func update_UserIconNotFound_ReturnError(s *MainCategSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockInput := domain.UpdateMainCategInput{
		ID:       1,
		Name:     "Test",
		IconType: domain.IconTypeCustom,
		IconID:   1,
	}

	// prepare mock service
	s.mockMainCategRepo.On("GetByID", mockInput.ID, mockUserID).Return(&domain.MainCateg{}, nil)
	s.mockUserIconRepo.On("GetByID", mockCtx, mockInput.IconID, mockUserID).Return(domain.UserIcon{}, domain.ErrUserIconNotFound)

	// action, assertion
	err := s.uc.Update(mockCtx, mockInput, mockUserID)
	s.Require().ErrorIs(err, domain.ErrUserIconNotFound, desc)
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
	err := s.uc.Delete(mockID)
	s.Require().NoError(err, desc)
}
