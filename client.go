package rabbitmqutils

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// MessageBrokerConnection ...
type MessageBrokerConnection interface {
	Channel() (*amqp.Channel, error)
	Close() error
}

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

// NewClientWithOptions ...
func NewClientWithOptions(dsn string, options ...ClientOption) (Client, error) {
	var clientOptions ClientOptions
	for _, option := range options {
		option(&clientOptions)
	}

	client, err := NewClient(dsn)
	if err != nil {
		return Client{}, err
	}

	if clientOptions.maximalQueueSize != 0 {
		if err := client.channel.Qos(
			clientOptions.maximalQueueSize, // prefetch count
			0,                              // prefetch size
			false,                          // global
		); err != nil {
			return Client{}, errors.Wrap(err, "unable to set the maximal queue size")
		}
	}

	for _, queue := range clientOptions.queues {
		if _, err := client.channel.QueueDeclare(
			queue, // queue name
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no wait
			nil,   // arguments
		); err != nil {
			return Client{}, errors.Wrapf(err, "unable to declare queue %q", queue)
		}
	}

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

// ConsumeMessages ...
func (client Client) ConsumeMessages(queue string) (
	<-chan amqp.Delivery,
	error,
) {
	return client.channel.Consume(
		queue,                   // queue name
		makeConsumerName(queue), // consumer name
		false,                   // auto-acknowledge
		false,                   // exclusive
		false,                   // no local
		false,                   // no wait
		nil,                     // arguments
	)
}

// CancelConsuming ...
func (client Client) CancelConsuming(queue string) error {
	return client.channel.Cancel(
		makeConsumerName(queue), // consumer name
		false,                   // no wait
	)
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

func makeConsumerName(queue string) string {
	return queue + "_consumer"
}
