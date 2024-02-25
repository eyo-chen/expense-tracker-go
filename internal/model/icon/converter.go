package icon

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func cvtToDomainIcon(m Icon) domain.Icon {
	return domain.Icon{
		ID:  m.ID,
		URL: m.URL,
	}
}
