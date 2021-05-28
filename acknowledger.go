package rabbitmqutils

import (
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
