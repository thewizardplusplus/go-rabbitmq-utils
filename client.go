package rabbitmqutils

import (
	"github.com/streadway/amqp"
)

// Client ...
type Client struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}
