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
