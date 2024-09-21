package icon

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func cvtToIcon(icons []domain.Icon) []icon {
	iconList := make([]icon, len(icons))
	for i, ic := range icons {
		iconList[i] = icon{
			ID:   ic.ID,
			Type: ic.Type.ToString(),
			URL:  ic.URL,
		}
	}
	return iconList
}
