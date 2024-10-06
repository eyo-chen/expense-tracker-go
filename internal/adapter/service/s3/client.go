package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

// NewS3Clients creates a new S3 presign client and S3 client.
func NewS3Clients(region, key, secret string) (interfaces.S3Client, interfaces.S3PresignClient) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, "")),
	)
	if err != nil {
		logger.Error("Unable to create S3 clients. Use empty S3 client and presign client", "error", err)
		return NewEmptyS3Client(), NewEmptyS3PresignClient()
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)
	return s3Client, presignClient
}
