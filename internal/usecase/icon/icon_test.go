package icon

import (
	"context"
	"errors"
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
	mockRedisService *mocks.RedisService
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupTest() {
	s.mockIconRepo = mocks.NewIconRepo(s.T())
	s.mockRedisService = mocks.NewRedisService(s.T())
	s.uc = New(s.mockIconRepo, s.mockRedisService)
}

func (s *IconSuite) TearDownTest() {
	s.mockIconRepo.AssertExpectations(s.T())
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
