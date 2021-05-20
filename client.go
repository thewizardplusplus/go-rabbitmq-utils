package rabbitmqutils

import (
	"encoding/json"

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

// PublishMessage ...
func (client Client) PublishMessage(queue string, message interface{}) error {
	messageAsJSON, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "unable to marshal the message")
	}

	if err := client.channel.Publish(
		"",    // exchange
		queue, // queue name
		false, // mandatory
		false, // immediate
		amqp.Publishing{ContentType: "application/json", Body: messageAsJSON},
	); err != nil {
		return errors.Wrap(err, "unable to publish the message")
	}

	return nil
}

// Close ...
func (client Client) Close() error {
	if err := client.channel.Close(); err != nil {
		return errors.Wrap(err, "unable to close the channel")
	}

	if err := client.connection.Close(); err != nil {
		return errors.Wrap(err, "unable to close the connection")
	}

	return nil
}