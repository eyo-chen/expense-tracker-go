package maincateg

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
)

func cvtToDomainMainCateg(c *MainCateg, i *icon.Icon) *domain.MainCateg {
	return &domain.MainCateg{
		ID:   c.ID,
		Name: c.Name,
		Type: domain.CvtToMainCategType(c.Type),
		Icon: cvtToDomainIcon(i),
	}
}

func cvtToDomainIcon(i *icon.Icon) *domain.Icon {
	if i == nil {
		return nil
	}

	return &domain.Icon{
		ID:  i.ID,
		URL: i.URL,
	}
}

func cvtToMainCateg(c *domain.MainCateg, userID int64) *MainCateg {
	return &MainCateg{
		ID:     c.ID,
		Name:   c.Name,
		Type:   c.Type.ModelValue(),
		IconID: c.Icon.ID,
		UserID: userID,
	}
}
