package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

var (
	packageName = "adapter/service/s3"
)

type preSigner struct {
	bucket string
	client *s3.PresignClient
}

// NewPreSigner creates a new s3 pre-signer.
func NewPreSigner(bucket string, client *s3.PresignClient) *preSigner {
	return &preSigner{bucket: bucket, client: client}
}

func (p *preSigner) PutObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error) {
	req, err := p.client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &p.bucket,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		logger.Error("Failed to get presigned URL", "error", err, "package", packageName)
		return "", err
	}

	return req.URL, nil
}

func (p *preSigner) GetObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error) {
	req, err := p.client.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &p.bucket,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		logger.Error("Failed to get presigned URL", "error", err, "package", packageName)
		return "", err
	}

	return req.URL, nil
}

func (p *preSigner) DeleteObject(ctx context.Context, objectKey string) error {
	_, err := p.client.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &p.bucket,
		Key:    &objectKey,
	})
	if err != nil {
		logger.Error("Failed to get presigned URL", "error", err, "package", packageName)
		return err
	}

	return nil
}
