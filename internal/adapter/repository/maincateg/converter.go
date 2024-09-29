package maincateg

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func cvtToDomainMainCateg(c MainCateg) domain.MainCateg {
	return domain.MainCateg{
		ID:       c.ID,
		Name:     c.Name,
		Type:     domain.CvtToTransactionType(c.Type),
		IconType: domain.CvtToIconType(c.IconType),
		IconData: c.IconData,
	}
}

func cvtToMainCateg(c domain.MainCateg, userID int64) MainCateg {
	return MainCateg{
		ID:       c.ID,
		Name:     c.Name,
		Type:     c.Type.ToModelValue(),
		UserID:   userID,
		IconType: c.IconType.ToModelValue(),
		IconData: c.IconData,
	}
}
