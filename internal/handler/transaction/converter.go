package transaction

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

func cvtToGetTransactionResp(trans []domain.Transaction) getTransactionResp {
	resp := make([]transaction, 0, len(trans))

	for _, t := range trans {
		resp = append(resp, transaction{
			ID:   t.ID,
			Type: t.Type.ToString(),
			MainCateg: mainCateg{
				ID:       t.MainCateg.ID,
				Name:     t.MainCateg.Name,
				Type:     t.MainCateg.Type.ToString(),
				IconType: t.MainCateg.IconType.ToString(),
				IconData: t.MainCateg.IconData,
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

func cvtToGetMonthlyResp(data []domain.TransactionType) []string {
	resp := make([]string, len(data))

	for i, d := range data {
		if d == domain.TransactionTypeUnSpecified {
			resp[i] = "no data"
		} else {
			resp[i] = d.ToString()
		}
	}

	return resp
}
