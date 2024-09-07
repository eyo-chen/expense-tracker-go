package icon

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/interfaces"
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
)

type IconUC struct {
	icon  interfaces.IconRepo
	redis interfaces.RedisService
}

func NewIconUC(i interfaces.IconRepo, r interfaces.RedisService) *IconUC {
	return &IconUC{
		icon:  i,
		redis: r,
	}
}

func (i *IconUC) List() ([]domain.Icon, error) {
	ctx := context.Background()

	res, err := i.redis.GetByFunc(ctx, "icons", 7*24*time.Hour, func() (string, error) {
		icons, err := i.icon.List()
		if err != nil {
			return "", err
		}

		return jsonutil.CvtToJSON(icons)
	})
	if err != nil {
		return nil, err
	}

	return jsonutil.CvtFromJSON[[]domain.Icon](res)
}
