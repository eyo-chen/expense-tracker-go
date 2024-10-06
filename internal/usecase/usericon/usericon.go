package usericon

import (
	"context"
	"fmt"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	s3 interfaces.S3Service
	ui interfaces.UserIconRepo
}

func New(s3 interfaces.S3Service, ui interfaces.UserIconRepo) *UC {
	return &UC{
		s3: s3,
		ui: ui,
	}
}

func (u *UC) GetPutObjectURL(ctx context.Context, fileName string, userID int64) (string, error) {
	objectKey := fmt.Sprintf("user_icons/%d/%s", userID, fileName)
	ttl := 60 * time.Second
	url, err := u.s3.PutObjectUrl(ctx, objectKey, int64(ttl.Seconds()))
	if err != nil {
		return "", err
	}

	return url, nil
}

func (u *UC) Create(ctx context.Context, fileName string, userID int64) error {
	objectKey := fmt.Sprintf("user_icons/%d/%s", userID, fileName)

	ui := domain.UserIcon{
		UserID:    userID,
		ObjectKey: objectKey,
	}

	return u.ui.Create(ctx, ui)
}
