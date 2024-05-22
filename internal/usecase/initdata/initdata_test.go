package initdata

import (
	"errors"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type InitDataSuite struct {
	suite.Suite
	uc            interfaces.InitDataUC
	mockIcon      *mocks.IconModel
	mockMainCateg *mocks.MainCategModel
	mockSubCateg  *mocks.SubCategModel
}

func TestInitDataSuite(t *testing.T) {
	suite.Run(t, new(InitDataSuite))
}

func (s *InitDataSuite) SetupTest() {
	s.mockIcon = mocks.NewIconModel(s.T())
	s.mockMainCateg = mocks.NewMainCategModel(s.T())
	s.mockSubCateg = mocks.NewSubCategModel(s.T())

	s.uc = NewInitDataUC(s.mockIcon, s.mockMainCateg, s.mockSubCateg)
}

func (s *InitDataSuite) TearDownTest() {
	s.mockIcon.AssertExpectations(s.T())
}

func (s *InitDataSuite) TestList() {
	for scenario, fn := range map[string]func(s *InitDataSuite, desc string){
		"when no error, return init data":  list_NoError_ReturnInitData,
		"when get icon fail, return error": list_GetIconFail_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func list_NoError_ReturnInitData(s *InitDataSuite, desc string) {
	mockIDToIcon := map[int64]domain.Icon{
		1:  {ID: 1, URL: "url1"},
		2:  {ID: 2, URL: "url2"},
		3:  {ID: 3, URL: "url3"},
		4:  {ID: 4, URL: "url4"},
		5:  {ID: 5, URL: "url5"},
		6:  {ID: 6, URL: "url6"},
		7:  {ID: 7, URL: "url7"},
		8:  {ID: 8, URL: "url8"},
		9:  {ID: 9, URL: "url9"},
		10: {ID: 10, URL: "url10"},
		11: {ID: 11, URL: "url11"},
		12: {ID: 12, URL: "url12"},
		14: {ID: 14, URL: "url14"},
		15: {ID: 15, URL: "url15"},
	}

	s.mockIcon.On("GetByIDs",
		[]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15},
	).Return(mockIDToIcon, nil)

	expRes := domain.InitData{
		Expense: []domain.InitDataMainCateg{
			{
				Name: "food",
				Icon: mockIDToIcon[1],
				SubCategs: []string{
					"breakfast", "brunch", "lunch", "dinner", "groceries", "drink", "snak",
				},
			},
			{
				Name: "transportation",
				Icon: mockIDToIcon[4],
				SubCategs: []string{
					"bus", "train", "MRT", "taxi", "uber", "gasoline", "parking fees", "repairs", "maintenance",
				},
			},
			{
				Name: "utilities",
				Icon: mockIDToIcon[9],
				SubCategs: []string{
					"electricity", "water", "internet", "phones", "garbage", "cable",
				},
			},
			{
				Name: "housing",
				Icon: mockIDToIcon[3],
				SubCategs: []string{
					"rent", "mortgage", "property taxes", "insurance", "repairs", "furnishings",
				},
			},
			{
				Name: "clothing",
				Icon: mockIDToIcon[2],
				SubCategs: []string{
					"shirts", "pants", "shoes", "accessories", "jewelry", "underwear", "socks",
				},
			},
			{
				Name: "entertainment",
				Icon: mockIDToIcon[6],
				SubCategs: []string{
					"movies", "concerts", "shows", "games", "toys", "hobbies", "books", "magazines", "music", "apps", "party", "vacations", "membership", "subscriptions",
				},
			},
			{
				Name: "gifts",
				Icon: mockIDToIcon[7],
				SubCategs: []string{
					"birthday", "wedding", "baby shower", "anniversary", "graduation", "holiday", "charities",
				},
			},
			{
				Name: "education",
				Icon: mockIDToIcon[5],
				SubCategs: []string{
					"tuition", "books", "course",
				},
			},
			{
				Name: "insurance",
				Icon: mockIDToIcon[10],
				SubCategs: []string{
					"health", "life", "auto", "home", "disability", "liability",
				},
			},
			{
				Name: "debt",
				Icon: mockIDToIcon[11],
				SubCategs: []string{
					"credit card", "student loans", "personal loans",
				},
			},
			{
				Name: "healthcare",
				Icon: mockIDToIcon[8],
				SubCategs: []string{
					"doctor", "dentist", "optometrist", "medication", "pharmacy", "hospital", "medical devices",
				},
			},
			{
				Name: "others",
				Icon: mockIDToIcon[14],
				SubCategs: []string{
					"others",
				},
			},
		},
		Income: []domain.InitDataMainCateg{
			{
				Name: "salary",
				Icon: mockIDToIcon[12],
				SubCategs: []string{
					"salary", "bonus", "commission", "tips",
				},
			},
			{
				Name: "investment",
				Icon: mockIDToIcon[15],
				SubCategs: []string{
					"dividends", "capital gains", "interest",
				},
			},
			{
				Name: "others",
				Icon: mockIDToIcon[14],
				SubCategs: []string{
					"others",
				},
			},
		},
	}

	res, err := s.uc.List()
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, res, desc)
}

func list_GetIconFail_ReturnError(s *InitDataSuite, desc string) {
	s.mockIcon.On("GetByIDs",
		[]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15},
	).Return(nil, errors.New("GetByIDs failed"))

	res, err := s.uc.List()
	s.Require().EqualError(err, "GetByIDs failed", desc)
	s.Require().Empty(res, desc)
}
