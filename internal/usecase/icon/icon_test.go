package icon

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
	mockCTX = context.Background()
)

type IconSuite struct {
	suite.Suite
	uc               *UC
	mockIconRepo     *mocks.IconRepo
	mockUserIconRepo *mocks.UserIconRepo
	mockRedisService *mocks.RedisService
	mockS3Service    *mocks.S3Service
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupTest() {
	s.mockIconRepo = mocks.NewIconRepo(s.T())
	s.mockUserIconRepo = mocks.NewUserIconRepo(s.T())
	s.mockRedisService = mocks.NewRedisService(s.T())
	s.mockS3Service = mocks.NewS3Service(s.T())
	s.uc = New(s.mockIconRepo, s.mockUserIconRepo, s.mockRedisService, s.mockS3Service)
}

func (s *IconSuite) TearDownTest() {
	s.mockIconRepo.AssertExpectations(s.T())
	s.mockUserIconRepo.AssertExpectations(s.T())
	s.mockRedisService.AssertExpectations(s.T())
	s.mockS3Service.AssertExpectations(s.T())
}

func (s *IconSuite) TestList() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when no error, return icon list": list_NoError_ReturnIconList,
		"when cache failed, return error": list_CacheFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func list_NoError_ReturnIconList(s *IconSuite, desc string) {
	// prepare mock data
	mockIconsStr := `[{"id":1,"url":"http://test.com/1"},{"id":2,"url":"http://test.com/1"}]`
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockTTL := 7 * 24 * time.Hour

	// mock service
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockTTL, mockGetFun).Return(mockIconsStr, nil)

	// prepare expected result
	expResp := []domain.DefaultIcon{
		{ID: 1, URL: "http://test.com/1"},
		{ID: 2, URL: "http://test.com/1"},
	}

	// test function
	icons, err := s.uc.List()
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, icons, desc)
}

func list_CacheFailed_ReturnError(s *IconSuite, desc string) {
	// prepare mock data
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockTTL := 7 * 24 * time.Hour
	mockErr := errors.New("cache failed")

	// mock service
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockTTL, mockGetFun).Return("", mockErr)

	// test function
	icons, err := s.uc.List()
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Nil(icons, desc)
}

func (s *IconSuite) TestListByUserID() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when no error, return icon list":              listByUserID_NoError_ReturnSuccessfully,
		"when list default icons failed, return error": listByUserID_ListDefaultIconsFailed_ReturnError,
		"when get user icons failed, return error":     listByUserID_GetUserIconsFailed_ReturnError,
		"when get object url failed, return error":     listByUserID_GetObjectUrlFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func listByUserID_NoError_ReturnSuccessfully(s *IconSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockIconsStr := `[{"id":1,"url":"http://test.com/1"},{"id":2,"url":"http://test.com/2"}]`
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockTTL := 7 * 24 * time.Hour
	mockUserIcons := []domain.UserIcon{{ID: 1, UserID: 1, ObjectKey: "test/1"}, {ID: 2, UserID: 1, ObjectKey: "test/2"}}
	mockKey1 := fmt.Sprintf("user_icon-%s", mockUserIcons[0].ObjectKey)
	mockKey2 := fmt.Sprintf("user_icon-%s", mockUserIcons[1].ObjectKey)

	// prepare mock service
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockTTL, mockGetFun).Return(mockIconsStr, nil)
	s.mockUserIconRepo.On("GetByUserID", mockCTX, mockUserID).Return(mockUserIcons, nil)
	s.mockRedisService.On("GetByFunc", mockCTX, mockKey1, mockTTL, mockGetFun).Return("http://s3.com/1", nil)
	s.mockRedisService.On("GetByFunc", mockCTX, mockKey2, mockTTL, mockGetFun).Return("http://s3.com/2", nil)

	// prepare expected result
	expResp := []domain.Icon{
		{ID: 1, Type: domain.IconTypeCustom, URL: "http://s3.com/1"},
		{ID: 2, Type: domain.IconTypeCustom, URL: "http://s3.com/2"},
		{ID: 1, Type: domain.IconTypeDefault, URL: "http://test.com/1"},
		{ID: 2, Type: domain.IconTypeDefault, URL: "http://test.com/2"},
	}

	// test function
	icons, err := s.uc.ListByUserID(mockCTX, mockUserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResp, icons, desc)
}

func listByUserID_ListDefaultIconsFailed_ReturnError(s *IconSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockTTL := 7 * 24 * time.Hour
	mockErr := errors.New("list default icons failed")

	// prepare mock service
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockTTL, mockGetFun).Return("", mockErr)

	// test function
	icons, err := s.uc.ListByUserID(mockCTX, mockUserID)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Nil(icons, desc)
}

func listByUserID_GetUserIconsFailed_ReturnError(s *IconSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockIconsStr := `[{"id":1,"url":"http://test.com/1"},{"id":2,"url":"http://test.com/2"}]`
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockTTL := 7 * 24 * time.Hour
	mockErr := errors.New("get user icons failed")

	// prepare mock service
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockTTL, mockGetFun).Return(mockIconsStr, nil)
	s.mockUserIconRepo.On("GetByUserID", mockCTX, mockUserID).Return(nil, mockErr)

	// test function
	icons, err := s.uc.ListByUserID(mockCTX, mockUserID)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Nil(icons, desc)
}

func listByUserID_GetObjectUrlFailed_ReturnError(s *IconSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockIconsStr := `[{"id":1,"url":"http://test.com/1"},{"id":2,"url":"http://test.com/2"}]`
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	mockTTL := 7 * 24 * time.Hour
	mockUserIcons := []domain.UserIcon{{ID: 1, UserID: 1, ObjectKey: "test/1"}, {ID: 2, UserID: 1, ObjectKey: "test/2"}}
	mockKey1 := fmt.Sprintf("user_icon-%s", mockUserIcons[0].ObjectKey)
	mockErr := errors.New("get object url failed")

	// prepare mock service
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockTTL, mockGetFun).Return(mockIconsStr, nil)
	s.mockUserIconRepo.On("GetByUserID", mockCTX, mockUserID).Return(mockUserIcons, nil)
	s.mockRedisService.On("GetByFunc", mockCTX, mockKey1, mockTTL, mockGetFun).Return("", mockErr)
	// test function
	icons, err := s.uc.ListByUserID(mockCTX, mockUserID)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Nil(icons, desc)
}
