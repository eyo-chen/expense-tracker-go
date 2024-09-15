package icon

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

func cvtToDomainDefaultIcon(m Icon) domain.DefaultIcon {
	return domain.DefaultIcon{
		ID:  m.ID,
		URL: m.URL,
	}
}

func cvtToIDToDomainDefaultIcon(icons []Icon) map[int64]domain.DefaultIcon {
	result := make(map[int64]domain.DefaultIcon)

	for _, i := range icons {
		result[i.ID] = cvtToDomainDefaultIcon(i)
	}

	return result
}
