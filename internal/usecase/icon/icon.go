package icon

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
)

type UC struct {
	icon  interfaces.IconRepo
	redis interfaces.RedisService
}

func New(i interfaces.IconRepo, r interfaces.RedisService) *UC {
	return &UC{
		icon:  i,
		redis: r,
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
