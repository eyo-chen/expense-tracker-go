package mq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type emptyMQClient struct{}

func NewEmptyMQClient() *emptyMQClient {
	return &emptyMQClient{}
}

// PublishWithContext publishes a message to a message queue.
func (e *emptyMQClient) PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return nil
}

// ConsumeWithContext consumes messages from a message queue.
func (e *emptyMQClient) ConsumeWithContext(ctx context.Context, queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, nil
}

// QueueDeclare declares a queue.
func (e *emptyMQClient) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

// Close closes the message queue client.
func (e *emptyMQClient) Close() error {
	return nil
}
