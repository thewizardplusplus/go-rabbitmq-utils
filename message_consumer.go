package rabbitmqutils

import (
	"context"

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
