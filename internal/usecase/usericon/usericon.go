package usericon

import (
	"context"
	"fmt"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
)

type UC struct {
	s3 interfaces.S3Service
}

func New(s3 interfaces.S3Service) *UC {
	return &UC{
		s3: s3,
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
