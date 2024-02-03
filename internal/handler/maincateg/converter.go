package maincateg

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToGetAllMainCategResp(c []domain.MainCateg) GetAllMainCategResp {
	categs := make([]mainCateg, 0, len(c))

	for _, v := range c {
		categs = append(categs, mainCateg{
			ID:   v.ID,
			Name: v.Name,
			Type: v.Type.String(),
			Icon: icon{
				ID:  v.Icon.ID,
				URL: v.Icon.URL,
			},
		})
	}

	return GetAllMainCategResp{
		Categories: categs,
	}
}
