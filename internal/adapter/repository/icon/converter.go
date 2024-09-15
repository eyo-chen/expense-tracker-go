package icon

import "github.com/eyo-chen/expense-tracker-go/internal/domain"

func cvtToDomainIcon(m Icon) domain.DefaultIcon {
	return domain.DefaultIcon{
		ID:  m.ID,
		URL: m.URL,
	}
}

func cvtToIDToDomainIcon(icons []Icon) map[int64]domain.DefaultIcon {
	result := make(map[int64]domain.DefaultIcon)

	for _, i := range icons {
		result[i.ID] = cvtToDomainIcon(i)
	}

	return result
}
