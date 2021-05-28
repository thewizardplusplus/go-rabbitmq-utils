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
		acknowledger.Logger.Logf("unable to handle the message: %v", err)

		requeue :=
			acknowledger.MessageHandling == TwiceMessageHandling && !message.Redelivered
		message.Reject(requeue)

		return
	}

	acknowledger.Logger.Log("message has been handled successfully")
	message.Ack(false /* multiple */)
}
