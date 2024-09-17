package maincateg

import (
	"context"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	MainCateg interfaces.MainCategRepo
	Icon      interfaces.IconRepo
	UserIcon  interfaces.UserIconRepo
}

func New(m interfaces.MainCategRepo, i interfaces.IconRepo, ui interfaces.UserIconRepo) *UC {
	return &UC{
		MainCateg: m,
		Icon:      i,
		UserIcon:  ui,
	}
}

func (u *UC) Create(categ domain.MainCateg, userID int64) error {
	if categ.IconType == domain.IconTypeUnspecified {
		return domain.ErrIconNotFound
	}

	ctx := context.Background()
	if categ.IconType == domain.IconTypeDefault {
		if _, err := u.Icon.GetByURL(ctx, categ.IconData); err != nil {
			return err
		}
	}

	if categ.IconType == domain.IconTypeCustom {
		if _, err := u.UserIcon.GetByObjectKeyAndUserID(ctx, categ.IconData, userID); err != nil {
			return err
		}
	}

	return u.MainCateg.Create(&categ, userID)
}

func (u *UC) GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error) {
	return u.MainCateg.GetAll(ctx, userID, transType)
}

func (u *UC) Update(categ domain.MainCateg, userID int64) error {
	// check if the main category exists
	if _, err := u.MainCateg.GetByID(categ.ID, userID); err != nil {
		return err
	}

	if categ.IconType == domain.IconTypeUnspecified {
		return domain.ErrIconNotFound
	}

	ctx := context.Background()
	if categ.IconType == domain.IconTypeDefault {
		if _, err := u.Icon.GetByURL(ctx, categ.IconData); err != nil {
			return err
		}
	}

	if categ.IconType == domain.IconTypeCustom {
		if _, err := u.UserIcon.GetByObjectKeyAndUserID(ctx, categ.IconData, userID); err != nil {
			return err
		}
	}

	return u.MainCateg.Update(&categ)
}

func (u *UC) Delete(id int64) error {
	return u.MainCateg.Delete(id)
}
