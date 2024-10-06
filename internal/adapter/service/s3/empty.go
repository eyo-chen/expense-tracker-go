package s3

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
This file contains empty implementations of S3 client interfaces.
These empty implementations serve as fallback options when the real S3 service
cannot be connected to or initialized.

If the S3 service connection fails during the application startup,
these empty clients will be used to replace the real S3 clients.
This allows the application to start and run, albeit with limited S3 functionality.

The empty clients implement the same interfaces as the real S3 clients,
but their methods do nothing and return zero values or empty results.
This approach helps prevent application crashes due to S3 connectivity issues,
while still maintaining the expected interface structure.

Usage of these empty clients should be logged appropriately, and the application
should handle the absence of actual S3 functionality gracefully where required.
*/

type emptyS3Client struct{}

// NewEmptyS3Client creates a new empty S3 client.
func NewEmptyS3Client() *emptyS3Client {
	return &emptyS3Client{}
}

func (e *emptyS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return nil, nil
}

type emptyS3PresignClient struct{}

// NewEmptyS3PresignClient creates a new empty S3 presign client.
func NewEmptyS3PresignClient() *emptyS3PresignClient {
	return &emptyS3PresignClient{}
}

func (e *emptyS3PresignClient) PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	return nil, nil
}

func (e *emptyS3PresignClient) PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	return nil, nil
}
