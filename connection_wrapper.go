package rabbitmqutils

import (
	"github.com/streadway/amqp"
)

// ConnectionWrapper ...
type ConnectionWrapper struct {
	*amqp.Connection
}

// Channel ...
func (wrapper ConnectionWrapper) Channel() (MessageBrokerChannel, error) {
	return wrapper.Connection.Channel()
}
