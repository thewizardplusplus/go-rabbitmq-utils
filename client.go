package rabbitmqutils

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Client ...
type Client struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewClient ...
func NewClient(dsn string) (Client, error) {
	connection, err := amqp.Dial(dsn)
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to open the connection")
	}

	channel, err := connection.Channel()
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to open the channel")
	}

	client := Client{connection: connection, channel: channel}
	return client, nil
}
