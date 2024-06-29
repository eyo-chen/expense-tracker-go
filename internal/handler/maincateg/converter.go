package maincateg

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

func cvtToGetAllMainCategResp(c []domain.MainCateg) getAllMainCategResp {
	categs := make([]mainCateg, 0, len(c))

	for _, v := range c {
		categs = append(categs, mainCateg{
			ID:   v.ID,
			Name: v.Name,
			Type: v.Type.ToString(),
			Icon: icon{
				ID:  v.Icon.ID,
				URL: v.Icon.URL,
			},
		})
	}

	return getAllMainCategResp{
		Categories: categs,
	}
}
