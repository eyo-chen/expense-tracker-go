package icon

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToDomainIcon(m Icon) domain.Icon {
	return domain.Icon{
		ID:  m.ID,
		URL: m.URL,
	}
}

func cvtToIDToDomainIcon(icons []Icon) map[int64]domain.Icon {
	result := make(map[int64]domain.Icon)

	for _, i := range icons {
		result[i.ID] = cvtToDomainIcon(i)
	}

	return result
}
