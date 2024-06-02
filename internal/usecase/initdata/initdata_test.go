package initdata

import (
	"context"
	"errors"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/mocks"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx = context.Background()
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
	s.mockMainCateg.AssertExpectations(s.T())
	s.mockSubCateg.AssertExpectations(s.T())
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
	).Return(mockIDToIcon, nil).Once()

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
	).Return(nil, errors.New("GetByIDs failed")).Once()

	res, err := s.uc.List()
	s.Require().EqualError(err, "GetByIDs failed", desc)
	s.Require().Empty(res, desc)
}

func (s *InitDataSuite) TestCreate() {
	for scenario, fn := range map[string]func(s *InitDataSuite, desc string){
		"when no error, create successfully": create_NoError_CreateSuccessfully,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func create_NoError_CreateSuccessfully(s *InitDataSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockData := domain.InitData{
		Expense: []domain.InitDataMainCateg{
			{
				Name: "food",
				Icon: domain.Icon{ID: 1, URL: "url1"},
				SubCategs: []string{
					"breakfast", "brunch", "lunch", "dinner", "groceries", "drink", "snak",
				},
			},
			{
				Name: "transportation",
				Icon: domain.Icon{ID: 4, URL: "url4"},
				SubCategs: []string{
					"bus", "train", "MRT", "taxi", "uber", "gasoline", "parking fees", "repairs", "maintenance",
				},
			},
			{
				Name: "utilities",
				Icon: domain.Icon{ID: 9, URL: "url9"},
				SubCategs: []string{
					"electricity", "water", "internet", "phones", "garbage", "cable",
				},
			},
			{
				Name: "housing",
				Icon: domain.Icon{ID: 3, URL: "url3"},
				SubCategs: []string{
					"rent", "mortgage", "property taxes", "insurance", "repairs", "furnishings",
				},
			},
			{
				Name: "clothing",
				Icon: domain.Icon{ID: 2, URL: "url2"},
				SubCategs: []string{
					"shirts", "pants", "shoes", "accessories", "jewelry", "underwear", "socks",
				},
			},
			{
				Name: "entertainment",
				Icon: domain.Icon{ID: 6, URL: "url6"},
				SubCategs: []string{
					"movies", "concerts", "shows", "games", "toys", "hobbies", "books", "magazines", "music", "apps", "party", "vacations", "membership", "subscriptions",
				},
			},
			{
				Name: "gifts",
				Icon: domain.Icon{ID: 7, URL: "url7"},
				SubCategs: []string{
					"birthday", "wedding", "baby shower", "anniversary", "graduation", "holiday", "charities",
				},
			},
			{
				Name: "education",
				Icon: domain.Icon{ID: 5, URL: "url5"},
				SubCategs: []string{
					"tuition", "books", "course",
				},
			},
			{
				Name: "insurance",
				Icon: domain.Icon{ID: 10, URL: "url10"},
				SubCategs: []string{
					"health", "life", "auto", "home", "disability", "liability",
				},
			},
			{
				Name: "debt",
				Icon: domain.Icon{ID: 11, URL: "url11"},
				SubCategs: []string{
					"credit card", "student loans", "personal loans",
				},
			},
			{
				Name: "healthcare",
				Icon: domain.Icon{ID: 8, URL: "url8"},
				SubCategs: []string{
					"doctor", "dentist", "optometrist", "medication", "pharmacy", "hospital", "medical devices",
				},
			},
			{
				Name: "others",
				Icon: domain.Icon{ID: 14, URL: "url14"},
				SubCategs: []string{
					"others",
				},
			},
		},
		Income: []domain.InitDataMainCateg{
			{
				Name: "salary",
				Icon: domain.Icon{ID: 12, URL: "url12"},
				SubCategs: []string{
					"salary", "bonus", "commission", "tips",
				},
			},
			{
				Name: "investment",
				Icon: domain.Icon{ID: 15, URL: "url15"},
				SubCategs: []string{
					"dividends", "capital gains", "interest",
				},
			},
			{
				Name: "others",
				Icon: domain.Icon{ID: 14, URL: "url14"},
				SubCategs: []string{
					"others",
				},
			},
		},
	}
	mockMainCategs := []domain.MainCateg{
		{Name: "food", Icon: domain.Icon{ID: 1, URL: "url1"}, Type: domain.TransactionTypeExpense},
		{Name: "transportation", Icon: domain.Icon{ID: 4, URL: "url4"}, Type: domain.TransactionTypeExpense},
		{Name: "utilities", Icon: domain.Icon{ID: 9, URL: "url9"}, Type: domain.TransactionTypeExpense},
		{Name: "housing", Icon: domain.Icon{ID: 3, URL: "url3"}, Type: domain.TransactionTypeExpense},
		{Name: "clothing", Icon: domain.Icon{ID: 2, URL: "url2"}, Type: domain.TransactionTypeExpense},
		{Name: "entertainment", Icon: domain.Icon{ID: 6, URL: "url6"}, Type: domain.TransactionTypeExpense},
		{Name: "gifts", Icon: domain.Icon{ID: 7, URL: "url7"}, Type: domain.TransactionTypeExpense},
		{Name: "education", Icon: domain.Icon{ID: 5, URL: "url5"}, Type: domain.TransactionTypeExpense},
		{Name: "insurance", Icon: domain.Icon{ID: 10, URL: "url10"}, Type: domain.TransactionTypeExpense},
		{Name: "debt", Icon: domain.Icon{ID: 11, URL: "url11"}, Type: domain.TransactionTypeExpense},
		{Name: "healthcare", Icon: domain.Icon{ID: 8, URL: "url8"}, Type: domain.TransactionTypeExpense},
		{Name: "others", Icon: domain.Icon{ID: 14, URL: "url14"}, Type: domain.TransactionTypeExpense},
		{Name: "salary", Icon: domain.Icon{ID: 12, URL: "url12"}, Type: domain.TransactionTypeIncome},
		{Name: "investment", Icon: domain.Icon{ID: 15, URL: "url15"}, Type: domain.TransactionTypeIncome},
		{Name: "others", Icon: domain.Icon{ID: 14, URL: "url14"}, Type: domain.TransactionTypeIncome},
	}
	mockMainCategsWithID := []domain.MainCateg{
		{ID: 1, Name: "food", Icon: domain.Icon{ID: 1, URL: "url1"}, Type: domain.TransactionTypeExpense},
		{ID: 2, Name: "transportation", Icon: domain.Icon{ID: 4, URL: "url4"}, Type: domain.TransactionTypeExpense},
		{ID: 3, Name: "utilities", Icon: domain.Icon{ID: 9, URL: "url9"}, Type: domain.TransactionTypeExpense},
		{ID: 4, Name: "housing", Icon: domain.Icon{ID: 3, URL: "url3"}, Type: domain.TransactionTypeExpense},
		{ID: 5, Name: "clothing", Icon: domain.Icon{ID: 2, URL: "url2"}, Type: domain.TransactionTypeExpense},
		{ID: 6, Name: "entertainment", Icon: domain.Icon{ID: 6, URL: "url6"}, Type: domain.TransactionTypeExpense},
		{ID: 7, Name: "gifts", Icon: domain.Icon{ID: 7, URL: "url7"}, Type: domain.TransactionTypeExpense},
		{ID: 8, Name: "education", Icon: domain.Icon{ID: 5, URL: "url5"}, Type: domain.TransactionTypeExpense},
		{ID: 9, Name: "insurance", Icon: domain.Icon{ID: 10, URL: "url10"}, Type: domain.TransactionTypeExpense},
		{ID: 10, Name: "debt", Icon: domain.Icon{ID: 11, URL: "url11"}, Type: domain.TransactionTypeExpense},
		{ID: 11, Name: "healthcare", Icon: domain.Icon{ID: 8, URL: "url8"}, Type: domain.TransactionTypeExpense},
		{ID: 12, Name: "others", Icon: domain.Icon{ID: 14, URL: "url14"}, Type: domain.TransactionTypeExpense},
		{ID: 13, Name: "salary", Icon: domain.Icon{ID: 12, URL: "url12"}, Type: domain.TransactionTypeIncome},
		{ID: 14, Name: "investment", Icon: domain.Icon{ID: 15, URL: "url15"}, Type: domain.TransactionTypeIncome},
		{ID: 15, Name: "others", Icon: domain.Icon{ID: 14, URL: "url14"}, Type: domain.TransactionTypeIncome},
	}
	mockSubCategs := []domain.SubCateg{
		{Name: "breakfast", MainCategID: 1},
		{Name: "brunch", MainCategID: 1},
		{Name: "lunch", MainCategID: 1},
		{Name: "dinner", MainCategID: 1},
		{Name: "groceries", MainCategID: 1},
		{Name: "drink", MainCategID: 1},
		{Name: "snak", MainCategID: 1},
		{Name: "bus", MainCategID: 2},
		{Name: "train", MainCategID: 2},
		{Name: "MRT", MainCategID: 2},
		{Name: "taxi", MainCategID: 2},
		{Name: "uber", MainCategID: 2},
		{Name: "gasoline", MainCategID: 2},
		{Name: "parking fees", MainCategID: 2},
		{Name: "repairs", MainCategID: 2},
		{Name: "maintenance", MainCategID: 2},
		{Name: "electricity", MainCategID: 3},
		{Name: "water", MainCategID: 3},
		{Name: "internet", MainCategID: 3},
		{Name: "phones", MainCategID: 3},
		{Name: "garbage", MainCategID: 3},
		{Name: "cable", MainCategID: 3},
		{Name: "rent", MainCategID: 4},
		{Name: "mortgage", MainCategID: 4},
		{Name: "property taxes", MainCategID: 4},
		{Name: "insurance", MainCategID: 4},
		{Name: "repairs", MainCategID: 4},
		{Name: "furnishings", MainCategID: 4},
		{Name: "shirts", MainCategID: 5},
		{Name: "pants", MainCategID: 5},
		{Name: "shoes", MainCategID: 5},
		{Name: "accessories", MainCategID: 5},
		{Name: "jewelry", MainCategID: 5},
		{Name: "underwear", MainCategID: 5},
		{Name: "socks", MainCategID: 5},
		{Name: "movies", MainCategID: 6},
		{Name: "concerts", MainCategID: 6},
		{Name: "shows", MainCategID: 6},
		{Name: "games", MainCategID: 6},
		{Name: "toys", MainCategID: 6},
		{Name: "hobbies", MainCategID: 6},
		{Name: "books", MainCategID: 6},
		{Name: "magazines", MainCategID: 6},
		{Name: "music", MainCategID: 6},
		{Name: "apps", MainCategID: 6},
		{Name: "party", MainCategID: 6},
		{Name: "vacations", MainCategID: 6},
		{Name: "membership", MainCategID: 6},
		{Name: "subscriptions", MainCategID: 6},
		{Name: "birthday", MainCategID: 7},
		{Name: "wedding", MainCategID: 7},
		{Name: "baby shower", MainCategID: 7},
		{Name: "anniversary", MainCategID: 7},
		{Name: "graduation", MainCategID: 7},
		{Name: "holiday", MainCategID: 7},
		{Name: "charities", MainCategID: 7},
		{Name: "tuition", MainCategID: 8},
		{Name: "books", MainCategID: 8},
		{Name: "course", MainCategID: 8},
		{Name: "health", MainCategID: 9},
		{Name: "life", MainCategID: 9},
		{Name: "auto", MainCategID: 9},
		{Name: "home", MainCategID: 9},
		{Name: "disability", MainCategID: 9},
		{Name: "liability", MainCategID: 9},
		{Name: "credit card", MainCategID: 10},
		{Name: "student loans", MainCategID: 10},
		{Name: "personal loans", MainCategID: 10},
		{Name: "doctor", MainCategID: 11},
		{Name: "dentist", MainCategID: 11},
		{Name: "optometrist", MainCategID: 11},
		{Name: "medication", MainCategID: 11},
		{Name: "pharmacy", MainCategID: 11},
		{Name: "hospital", MainCategID: 11},
		{Name: "medical devices", MainCategID: 11},
		{Name: "others", MainCategID: 12},
		{Name: "salary", MainCategID: 13},
		{Name: "bonus", MainCategID: 13},
		{Name: "commission", MainCategID: 13},
		{Name: "tips", MainCategID: 13},
		{Name: "dividends", MainCategID: 14},
		{Name: "capital gains", MainCategID: 14},
		{Name: "interest", MainCategID: 14},
		{Name: "others", MainCategID: 15},
	}

	// mock service
	s.mockMainCateg.On("BatchCreate", mockMainCategs, mockUserID).Return(nil).Once()

	s.mockMainCateg.On("GetAll", mockUserID, domain.TransactionTypeUnSpecified).Return(mockMainCategsWithID, nil).Once()

	s.mockSubCateg.On("BatchCreate", mockCtx, mockSubCategs, mockUserID).Return(nil).Once()

	err := s.uc.Create(mockCtx, mockData, mockUserID)
	s.Require().NoError(err, desc)
}
