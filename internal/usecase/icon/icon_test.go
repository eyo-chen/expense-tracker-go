package icon

import (
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type IconSuite struct {
	suite.Suite
	iconUC        interfaces.IconUC
	mockIconModel *mocks.IconModel
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupTest() {
	s.mockIconModel = mocks.NewIconModel(s.T())
	s.iconUC = NewIconUC(s.mockIconModel)
}

func (s *IconSuite) TearDownTest() {
	s.mockIconModel.AssertExpectations(s.T())
}

func (s *IconSuite) TestList() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when no error, return icon list": list_NoError_ReturnIconList,
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
	mockIcons := []domain.Icon{
		{ID: 1, URL: "http://test.com/1"},
		{ID: 2, URL: "http://test.com/1"},
	}

	// mock the function
	s.mockIconModel.On("List").Return(mockIcons, nil)

	// test function
	icons, err := s.iconUC.List()
	s.Require().NoError(err)
	s.Require().Equal(mockIcons, icons)
}
