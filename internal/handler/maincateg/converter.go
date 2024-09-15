package maincateg

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

func cvtToGetAllMainCategResp(c []domain.MainCateg) getAllMainCategResp {
	categs := make([]mainCateg, 0, len(c))

	for _, v := range c {
		categs = append(categs, mainCateg{
			ID:       v.ID,
			Name:     v.Name,
			Type:     v.Type.ToString(),
			IconType: v.IconType.ToString(),
			IconData: v.IconData,
		})
	}

	return getAllMainCategResp{
		Categories: categs,
	}
}
