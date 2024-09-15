package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

var (
	packageName = "adapter/service/s3"
)

type Service struct {
	bucket        string
	s3Client      interfaces.S3Client
	presignClient interfaces.S3PresignClient
}

// New creates a new s3 service.
func New(bucket string, s3Client interfaces.S3Client, presignClient interfaces.S3PresignClient) *Service {
	return &Service{bucket: bucket, s3Client: s3Client, presignClient: presignClient}
}

func (s *Service) PutObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error) {
	req, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
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

func (s *Service) GetObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error) {
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
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

func (s *Service) DeleteObject(ctx context.Context, objectKey string) error {
	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &objectKey,
	})
	if err != nil {
		logger.Error("Failed to delete object", "error", err, "package", packageName)
		return err
	}

	return nil
}
