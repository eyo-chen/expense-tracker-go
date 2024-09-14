package initdata

import (
	"context"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	icon      interfaces.IconRepo
	mainCateg interfaces.MainCategRepo
	subCateg  interfaces.SubCategRepo
	user      interfaces.UserRepo
}

func New(
	i interfaces.IconRepo,
	m interfaces.MainCategRepo,
	s interfaces.SubCategRepo,
	u interfaces.UserRepo,
) *UC {
	return &UC{
		icon:      i,
		mainCateg: m,
		subCateg:  s,
		user:      u,
	}
}

// [food 1], [transportation 4], [utilities 9], [housing 3], [clothing 2], [entertainment 6], [gifts 7], [education 5], [insurance 10], [debt 11], [healthcare 8], [others 14]
// [salary 12], [investment 15], [others 14]
func (u *UC) List() (domain.InitData, error) {
	iconIDs := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15}
	idToIcon, err := u.icon.GetByIDs(iconIDs)
	if err != nil {
		return domain.InitData{}, err
	}

	return domain.InitData{
		Expense: []domain.InitDataMainCateg{
			{
				Name: "food",
				Icon: idToIcon[1],
				SubCategs: []string{
					"breakfast", "brunch", "lunch", "dinner", "groceries", "drink", "snak",
				},
			},
			{
				Name: "transportation",
				Icon: idToIcon[4],
				SubCategs: []string{
					"bus", "train", "MRT", "taxi", "uber", "gasoline", "parking fees", "repairs", "maintenance",
				},
			},
			{
				Name: "utilities",
				Icon: idToIcon[9],
				SubCategs: []string{
					"electricity", "water", "internet", "phones", "garbage", "cable",
				},
			},
			{
				Name: "housing",
				Icon: idToIcon[3],
				SubCategs: []string{
					"rent", "mortgage", "property taxes", "insurance", "repairs", "furnishings",
				},
			},
			{
				Name: "clothing",
				Icon: idToIcon[2],
				SubCategs: []string{
					"shirts", "pants", "shoes", "accessories", "jewelry", "underwear", "socks",
				},
			},
			{
				Name: "entertainment",
				Icon: idToIcon[6],
				SubCategs: []string{
					"movies", "concerts", "shows", "games", "toys", "hobbies", "books", "magazines", "music", "apps", "party", "vacations", "membership", "subscriptions",
				},
			},
			{
				Name: "gifts",
				Icon: idToIcon[7],
				SubCategs: []string{
					"birthday", "wedding", "baby shower", "anniversary", "graduation", "holiday", "charities",
				},
			},
			{
				Name: "education",
				Icon: idToIcon[5],
				SubCategs: []string{
					"tuition", "books", "course",
				},
			},
			{
				Name: "insurance",
				Icon: idToIcon[10],
				SubCategs: []string{
					"health", "life", "auto", "home", "disability", "liability",
				},
			},
			{
				Name: "debt",
				Icon: idToIcon[11],
				SubCategs: []string{
					"credit card", "student loans", "personal loans",
				},
			},
			{
				Name: "healthcare",
				Icon: idToIcon[8],
				SubCategs: []string{
					"doctor", "dentist", "optometrist", "medication", "pharmacy", "hospital", "medical devices",
				},
			},
			{
				Name: "others",
				Icon: idToIcon[14],
				SubCategs: []string{
					"others",
				},
			},
		},
		Income: []domain.InitDataMainCateg{
			{
				Name: "salary",
				Icon: idToIcon[12],
				SubCategs: []string{
					"salary", "bonus", "commission", "tips",
				},
			},
			{
				Name: "investment",
				Icon: idToIcon[15],
				SubCategs: []string{
					"dividends", "capital gains", "interest",
				},
			},
			{
				Name: "others",
				Icon: idToIcon[14],
				SubCategs: []string{
					"others",
				},
			},
		},
	}, nil
}

func (u *UC) Create(ctx context.Context, data domain.InitData, userID int64) error {
	mainCategs := genAllMainCategs(data)
	if len(mainCategs) == 0 {
		return nil
	}

	if err := u.mainCateg.BatchCreate(ctx, mainCategs, userID); err != nil {
		return err
	}

	allCategs, err := u.mainCateg.GetAll(ctx, userID, domain.TransactionTypeUnSpecified)
	if err != nil {
		return err
	}

	subCategs := genAllSubCategs(data, allCategs)
	if len(subCategs) == 0 {
		return nil
	}

	if err := u.subCateg.BatchCreate(ctx, subCategs, userID); err != nil {
		return err
	}

	t := true
	opt := domain.UpdateUserOpt{IsSetInitCategory: &t}
	return u.user.Update(ctx, userID, opt)
}
