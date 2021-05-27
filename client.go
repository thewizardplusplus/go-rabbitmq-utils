package rabbitmqutils

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

//go:generate mockery --name=MessageBrokerConnection --inpackage --case=underscore --testonly

// MessageBrokerConnection ...
type MessageBrokerConnection interface {
	Channel() (*amqp.Channel, error)
	Close() error
}

//go:generate mockery --name=MessageBrokerChannel --inpackage --case=underscore --testonly

// MessageBrokerChannel ...
type MessageBrokerChannel interface {
	Qos(prefetchCount int, prefetchSize int, global bool) error
	QueueDeclare(
		queueName string,
		durable bool,
		autoDelete bool,
		exclusive bool,
		noWait bool,
		arguments amqp.Table,
	) (amqp.Queue, error)
	Publish(
		exchange string,
		queueName string,
		mandatory bool,
		immediate bool,
		message amqp.Publishing,
	) error
	Consume(
		queueName string,
		consumerName string,
		autoAcknowledge bool,
		exclusive bool,
		noLocal bool,
		noWait bool,
		arguments amqp.Table,
	) (<-chan amqp.Delivery, error)
	Cancel(consumerName string, noWait bool) error
	Close() error
}

// Client ...
type Client struct {
	connection MessageBrokerConnection
	channel    MessageBrokerChannel
}

// NewClient ...
func NewClient(dsn string, options ...ClientOption) (Client, error) {
	var clientOptions ClientOptions
	for _, option := range options {
		option(&clientOptions)
	}

	connection, err := amqp.Dial(dsn)
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to open the connection")
	}

	channel, err := connection.Channel()
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to open the channel")
	}

	if clientOptions.maximalQueueSize != 0 {
		if err := channel.Qos(
			clientOptions.maximalQueueSize, // prefetch count
			0,                              // prefetch size
			false,                          // global
		); err != nil {
			return Client{}, errors.Wrap(err, "unable to set the maximal queue size")
		}
	}

	for _, queue := range clientOptions.queues {
		if _, err := channel.QueueDeclare(
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
