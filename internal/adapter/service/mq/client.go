package mq

import (
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// NewMQClient initializes a new MQ client.
func NewMQClient(url string) (interfaces.MQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("Failed to connect to message queue", "error", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to create channel", "error", err)
		return nil, err
	}

	return ch, nil
}
