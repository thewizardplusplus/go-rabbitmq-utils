package rabbitmqutils

import (
	"github.com/streadway/amqp"
)

// MessageHandler ...
type MessageHandler interface {
	HandleMessage(message amqp.Delivery)
}
