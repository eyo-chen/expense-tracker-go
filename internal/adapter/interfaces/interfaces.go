package interfaces

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
)

// S3PresignClient is the interface that wraps the basic methods for s3 presign client.
type S3PresignClient interface {
	// PresignPutObject returns a pre-signed URL to upload an object to S3.
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)

	// PresignGetObject returns a pre-signed URL to get an object from S3.
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// S3Client is the interface that wraps the basic methods for s3 client.
type S3Client interface {
	// DeleteObject deletes an object from S3.
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// MQClient is the interface that wraps the basic methods for message queue client.
type MQClient interface {
	// PublishWithContext publishes a message to a message queue.
	PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error

	// ConsumeWithContext consumes messages from a message queue.
	ConsumeWithContext(ctx context.Context, queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)

	// QueueDeclare declares a queue.
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)

	// Close closes the message queue client.
	Close() error
}
