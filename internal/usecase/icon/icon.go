package icon

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
)

type UC struct {
	icon     interfaces.IconRepo
	userIcon interfaces.UserIconRepo
	redis    interfaces.RedisService
	s3       interfaces.S3Service
}

func New(i interfaces.IconRepo,
	ui interfaces.UserIconRepo,
	r interfaces.RedisService,
	s3 interfaces.S3Service) *UC {
	return &UC{
		icon:     i,
		userIcon: ui,
		redis:    r,
		s3:       s3,
	}
}

func (u *UC) List() ([]domain.DefaultIcon, error) {
	ctx := context.Background()

	res, err := u.redis.GetByFunc(ctx, "icons", 7*24*time.Hour, func() (string, error) {
		icons, err := u.icon.List()
		if err != nil {
			return "", err
		}

		return jsonutil.CvtToJSON(icons)
	})
	if err != nil {
		return nil, err
	}

	return jsonutil.CvtFromJSON[[]domain.DefaultIcon](res)
}

func (u *UC) ListByUserID(ctx context.Context, userID int64) ([]domain.Icon, error) {
	defaultIcons, err := u.List()
	if err != nil {
		return nil, err
	}

	userIcons, err := u.userIcon.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	icons := make([]domain.Icon, 0, len(defaultIcons)+len(userIcons))
	for _, ui := range userIcons {
		key := domain.GenUserIconCacheKey(ui.ObjectKey)
		url, err := u.redis.GetByFunc(ctx, key, 7*24*time.Hour, func() (string, error) {
			presignedURL, err := u.s3.GetObjectUrl(ctx, ui.ObjectKey, int64((7 * 24 * time.Hour).Seconds()))
			if err != nil {
				return "", err
			}

			return presignedURL, nil
		})
		if err != nil {
			return nil, err
		}

		icons = append(icons, domain.Icon{
			ID:   ui.ID,
			Type: domain.IconTypeCustom,
			URL:  url,
		})
	}

	for _, di := range defaultIcons {
		icons = append(icons, domain.Icon{
			ID:   di.ID,
			Type: domain.IconTypeDefault,
			URL:  di.URL,
		})
	}

	return icons, nil
}
