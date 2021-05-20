package rabbitmqutils

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// MessageHandler ...
type MessageHandler interface {
	HandleMessage(message amqp.Delivery)
}

// MessageConsumer ...
type MessageConsumer struct {
	client               Client
	queue                string
	messages             <-chan amqp.Delivery
	messageHandler       MessageHandler
	stoppingCtx          context.Context
	stoppingCtxCanceller context.CancelFunc
}

// NewMessageConsumer ...
func NewMessageConsumer(
	client Client,
	queue string,
	messageHandler MessageHandler,
) (MessageConsumer, error) {
	messages, err := client.ConsumeMessages(queue)
	if err != nil {
		return MessageConsumer{}, errors.Wrap(err, "unable to start the consuming")
	}

	stoppingCtx, stoppingCtxCanceller := context.WithCancel(context.Background())
	messageConsumer := MessageConsumer{
		client:               client,
		queue:                queue,
		messages:             messages,
		messageHandler:       messageHandler,
		stoppingCtx:          stoppingCtx,
		stoppingCtxCanceller: stoppingCtxCanceller,
	}
	return messageConsumer, nil
}
