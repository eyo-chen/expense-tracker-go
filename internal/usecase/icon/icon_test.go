package icon

import (
	"context"
	"errors"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
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
	iconUC           interfaces.IconUC
	mockIconModel    *mocks.IconModel
	mockRedisService *mocks.RedisService
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupTest() {
	s.mockIconModel = mocks.NewIconModel(s.T())
	s.mockRedisService = mocks.NewRedisService(s.T())
	s.iconUC = NewIconUC(s.mockIconModel, s.mockRedisService)
}

func (s *IconSuite) TearDownTest() {
	s.mockIconModel.AssertExpectations(s.T())
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

	// mock the function
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockGetFun).Return(mockIconsStr, nil)

	// prepare expected result
	expResp := []domain.Icon{
		{ID: 1, URL: "http://test.com/1"},
		{ID: 2, URL: "http://test.com/1"},
	}

	// test function
	icons, err := s.iconUC.List()
	s.Require().NoError(err)
	s.Require().Equal(expResp, icons)
}

func list_CacheFailed_ReturnError(s *IconSuite, desc string) {
	// prepare mock data
	mockGetFun := mock.AnythingOfType("func() (string, error)")
	s.mockRedisService.On("GetByFunc", mockCTX, "icons", mockGetFun).Return("", errors.New("cache failed"))

	// prepare expected result
	expResp := []domain.Icon(nil)

	// test function
	icons, err := s.iconUC.List()
	s.Require().EqualError(err, "cache failed")
	s.Require().Equal(expResp, icons)
}
