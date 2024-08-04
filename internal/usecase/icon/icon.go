package icon

import (
	"context"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/interfaces"
)

type IconUC struct {
	icon  interfaces.IconModel
	redis interfaces.RedisService
}

func NewIconUC(i interfaces.IconModel, redis interfaces.RedisService) *IconUC {
	return &IconUC{
		icon:  i,
		redis: redis,
	}
}

func (i *IconUC) List() ([]domain.Icon, error) {
	ctx := context.Background()

	res, err := i.redis.GetByFunc(ctx, "icons", func() (string, error) {
		icons, err := i.icon.List()
		if err != nil {
			return "", err
		}

		return domain.CvtIconsToJSON(icons)
	})
	if err != nil {
		return nil, err
	}

	return domain.CvtJSONToIcons(res)
}
