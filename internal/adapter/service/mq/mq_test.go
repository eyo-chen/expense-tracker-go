package mq

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX       = context.Background()
	mockQueueName = "test-queue"
)

type mqServiceSuite struct {
	suite.Suite
	service      *Service
	mockMQClient *mocks.MQClient
}

func TestMQServiceSuite(t *testing.T) {
	suite.Run(t, new(mqServiceSuite))
}

func (s *mqServiceSuite) SetupSuite() {
	logger.Register()
}

func (s *mqServiceSuite) SetupTest() {
	s.mockMQClient = new(mocks.MQClient)

	var args amqp.Table
	s.mockMQClient.On("QueueDeclare", mockQueueName, true, false, false, false, args).Return(amqp.Queue{Name: mockQueueName}, nil)
	s.service = New(mockQueueName, s.mockMQClient)
}

func (s *mqServiceSuite) TearDownTest() {
	s.mockMQClient.AssertExpectations(s.T())
}

func (s *mqServiceSuite) TestPublish() {
	for scenario, fn := range map[string]func(s *mqServiceSuite, desc string){
		"when no error, publish successfully": publish_NoError_ReturnSuccessfully,
		"when publish failed, return error":   publish_Error_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func publish_NoError_ReturnSuccessfully(s *mqServiceSuite, desc string) {
	// prepare mock data
	mockMessage := "test-message"
	mockBody, err := json.Marshal(mockMessage)
	s.Require().NoError(err, desc)
	mockPublishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        mockBody,
	}

	s.mockMQClient.On("PublishWithContext", mockCTX, "", mockQueueName, false, false, mockPublishing).Return(nil)

	err = s.service.Publish(mockCTX, mockMessage)
	s.Require().NoError(err, desc)
}

func publish_Error_ReturnError(s *mqServiceSuite, desc string) {
	// prepare mock data
	mockMessage := "test-message"
	mockBody, err := json.Marshal(mockMessage)
	s.Require().NoError(err, desc)
	mockPublishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        mockBody,
	}
	mockError := errors.New("test-error")

	s.mockMQClient.On("PublishWithContext", mockCTX, "", mockQueueName, false, false, mockPublishing).Return(mockError)

	err = s.service.Publish(mockCTX, mockMessage)

	s.Require().ErrorIs(err, mockError, desc)
}

func (s *mqServiceSuite) TestClose() {
	// mock service
	s.mockMQClient.On("Close").Return(nil)

	// action
	s.service.Close()
}
