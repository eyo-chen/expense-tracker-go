package maincateg

import (
	"context"
	"fmt"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	MainCateg interfaces.MainCategRepo
	Icon      interfaces.IconRepo
	UserIcon  interfaces.UserIconRepo
	Redis     interfaces.RedisService
	S3        interfaces.S3Service
}

func New(m interfaces.MainCategRepo, i interfaces.IconRepo, ui interfaces.UserIconRepo, r interfaces.RedisService, s interfaces.S3Service) *UC {
	return &UC{
		MainCateg: m,
		Icon:      i,
		UserIcon:  ui,
		Redis:     r,
		S3:        s,
	}
}

func (u *UC) Create(ctx context.Context, categ domain.CreateMainCategInput, userID int64) error {
	if !categ.IconType.IsValid() {
		return domain.ErrIconNotFound
	}

	var iconData string
	if categ.IconType == domain.IconTypeDefault {
		i, err := u.Icon.GetByID(ctx, categ.IconID)
		if err != nil {
			return err
		}
		iconData = i.URL
	}

	if categ.IconType == domain.IconTypeCustom {
		ui, err := u.UserIcon.GetByID(ctx, categ.IconID, userID)
		if err != nil {
			return err
		}
		iconData = ui.ObjectKey
	}

	c := domain.MainCateg{
		Name:     categ.Name,
		Type:     categ.Type,
		IconType: categ.IconType,
		IconData: iconData,
	}
	return u.MainCateg.Create(ctx, c, userID)
}

func (u *UC) GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error) {
	categs, err := u.MainCateg.GetAll(ctx, userID, transType)
	if err != nil {
		return nil, err
	}

	// get and cache presigned URL of custom icons
	for i, categ := range categs {
		if categ.IconType != domain.IconTypeCustom {
			continue
		}

		key := fmt.Sprintf("user_icon-%s", categ.IconData)
		url, err := u.Redis.GetByFunc(ctx, key, 7*24*time.Hour, func() (string, error) {
			presignedURL, err := u.S3.GetObjectUrl(ctx, categ.IconData, int64((7 * 24 * time.Hour).Seconds()))
			if err != nil {
				return "", err
			}

			return presignedURL, nil
		})
		if err != nil {
			return nil, err
		}

		categs[i].IconData = url

	}

	return categs, nil
}

func (u *UC) Update(ctx context.Context, categ domain.UpdateMainCategInput, userID int64) error {
	// check if the main category exists
	if _, err := u.MainCateg.GetByID(categ.ID, userID); err != nil {
		return err
	}

	if !categ.IconType.IsValid() {
		return domain.ErrIconNotFound
	}

	var iconData string
	if categ.IconType == domain.IconTypeDefault {
		i, err := u.Icon.GetByID(ctx, categ.IconID)
		if err != nil {
			return err
		}
		iconData = i.URL
	}

	if categ.IconType == domain.IconTypeCustom {
		ui, err := u.UserIcon.GetByID(ctx, categ.IconID, userID)
		if err != nil {
			return err
		}
		iconData = ui.ObjectKey
	}

	c := domain.MainCateg{
		ID:       categ.ID,
		Name:     categ.Name,
		Type:     categ.Type,
		IconType: categ.IconType,
		IconData: iconData,
	}
	return u.MainCateg.Update(ctx, c)
}

func (u *UC) Delete(id int64) error {
	return u.MainCateg.Delete(id)
}
