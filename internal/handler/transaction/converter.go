package transaction

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToGetTransactionResp(trans []domain.Transaction) getTransactionResp {
	resp := make([]transaction, 0, len(trans))

	for _, t := range trans {
		resp = append(resp, transaction{
			ID: t.ID,
			MainCateg: mainCateg{
				ID:   t.MainCateg.ID,
				Name: t.MainCateg.Name,
				Type: t.MainCateg.Type.ToString(),
				Icon: icon{
					ID:  t.MainCateg.Icon.ID,
					URL: t.MainCateg.Icon.URL,
				},
			},
			SubCateg: subCateg{
				ID:   t.SubCateg.ID,
				Name: t.SubCateg.Name,
			},
			Price: t.Price,
			Note:  t.Note,
			Date:  t.Date,
		})
	}

	return getTransactionResp{Transactions: resp}
}
