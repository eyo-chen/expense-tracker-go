package mq

import (
	"context"
	"encoding/json"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	packageName = "adapter/service/mq"
)

type Service struct {
	QueueName string
	MQClient  interfaces.MQClient
}

// NewMQService initializes a new MQ service.
func New(queueName string, mqClient interfaces.MQClient) *Service {
	queue, err := mqClient.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		logger.Fatal("Failed to declare queue", "error", err, "package", packageName)
		return nil
	}

	return &Service{QueueName: queue.Name, MQClient: mqClient}
}

func (s *Service) Publish(ctx context.Context, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Failed to marshal message", "error", err, "package", packageName)
		return err
	}

	return s.MQClient.PublishWithContext(
		ctx,
		"", // default exchange
		s.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (s *Service) Close() {
	s.MQClient.Close()
}
