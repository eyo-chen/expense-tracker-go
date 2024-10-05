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
	categs, err := u.MainCateg.GetAll(ctx, userID, transType)
	if err != nil {
		return nil, err
	}

	// get and cache presigned URL of custom icons
	for _, categ := range categs {
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

		categ.IconData = url
	}

	return categs, nil
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
