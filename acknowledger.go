package rabbitmqutils

import (
	"github.com/go-log/log"
	"github.com/streadway/amqp"
)

// MessageHandling ...
type MessageHandling int

// ...
const (
	OnceMessageHandling MessageHandling = iota
	TwiceMessageHandling
)

//go:generate mockery --name=FailingMessageHandler --inpackage --case=underscore --testonly

// FailingMessageHandler ...
type FailingMessageHandler interface {
	HandleMessage(message amqp.Delivery) error
}

// Acknowledger ...
type Acknowledger struct {
	MessageHandling MessageHandling
	MessageHandler  FailingMessageHandler
	Logger          log.Logger
}

// HandleMessage ...
func (acknowledger Acknowledger) HandleMessage(message amqp.Delivery) {
	if err := acknowledger.MessageHandler.HandleMessage(message); err != nil {
		acknowledger.Logger.
			Logf("unable to handle the message %s: %v", message.MessageId, err)

		requeue :=
			acknowledger.MessageHandling == TwiceMessageHandling && !message.Redelivered
		message.Reject(requeue) // nolint: errcheck, gosec

		return
	}

	acknowledger.Logger.
		Logf("message %s has been handled successfully", message.MessageId)
	message.Ack(false /* multiple */) // nolint: errcheck, gosec
}
