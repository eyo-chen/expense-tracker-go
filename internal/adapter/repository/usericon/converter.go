package usericon

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func cvtToDomainUserIcon(userIcon userIcon) domain.UserIcon {
	return domain.UserIcon{
		ID:        userIcon.ID,
		UserID:    userIcon.UserID,
		ObjectKey: userIcon.ObjectKey,
	}
}
