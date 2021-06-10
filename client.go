package rabbitmqutils

import (
	"encoding/json"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var errUnknownQueue = errors.New("unknown queue")

//go:generate mockery --name=MessageBrokerConnection --inpackage --case=underscore --testonly

// MessageBrokerConnection ...
type MessageBrokerConnection interface {
	Channel() (MessageBrokerChannel, error)
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
	connection  MessageBrokerConnection
	channel     MessageBrokerChannel
	queues      mapset.Set
	idGenerator IDGenerator
	clock       Clock
}

// NewClient ...
func NewClient(dsn string, options ...ClientOption) (Client, error) {
	clientConfig := ClientConfig{
		dialer: func(dsn string) (MessageBrokerConnection, error) {
			connection, err := amqp.Dial(dsn)
			if err != nil {
				return nil, err
			}

			wrapper := ConnectionWrapper{Connection: connection}
			return wrapper, nil
		},
		queues: mapset.NewSet(),
		idGenerator: func() (string, error) {
			uuid, err := uuid.NewRandom()
			if err != nil {
				return "", err
			}

			return uuid.String(), nil
		},
		clock: time.Now,
	}
	for _, option := range options {
		option(&clientConfig)
	}

	connection, err := clientConfig.dialer(dsn)
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to open the connection")
	}

	channel, err := connection.Channel()
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to open the channel")
	}

	if clientConfig.maximalQueueSize != 0 {
		if err := channel.Qos(
			clientConfig.maximalQueueSize, // prefetch count
			0,                             // prefetch size
			false,                         // global
		); err != nil {
			return Client{}, errors.Wrap(err, "unable to set the maximal queue size")
		}
	}

	var queueErr error
	clientConfig.queues.Each(func(queue interface{}) bool {
		if _, err := channel.QueueDeclare(
			queue.(string), // queue name
			true,           // durable
			false,          // auto-delete
			false,          // exclusive
			false,          // no wait
			nil,            // arguments
		); err != nil {
			queueErr = errors.Wrapf(err, "unable to declare queue %q", queue)
			return true
		}

		return false
	})
	if queueErr != nil {
		return Client{}, queueErr
	}

	client := Client{
		connection:  connection,
		channel:     channel,
		queues:      clientConfig.queues,
		idGenerator: clientConfig.idGenerator,
		clock:       clientConfig.clock,
	}
	return client, nil
}

// PublishMessage ...
func (client Client) PublishMessage(
	queue string,
	messageID string,
	messageData interface{},
) error {
	if !client.queues.Contains(queue) {
		return errUnknownQueue
	}

	if messageID == "" {
		var err error
		messageID, err = client.idGenerator()
		if err != nil {
			return errors.Wrap(err, "unable to generate the message ID")
		}
	}

	messageDataAsJSON, err := json.Marshal(messageData)
	if err != nil {
		return errors.Wrap(err, "unable to marshal the message")
	}

	messageTimestamp := client.clock()
	if err := client.channel.Publish(
		"",    // exchange
		queue, // queue name
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			MessageId:   messageID,
			Timestamp:   messageTimestamp,
			ContentType: "application/json",
			Body:        messageDataAsJSON,
		},
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
	if !client.queues.Contains(queue) {
		return nil, errUnknownQueue
	}

	messages, err := client.channel.Consume(
		queue,                   // queue name
		makeConsumerName(queue), // consumer name
		false,                   // auto-acknowledge
		false,                   // exclusive
		false,                   // no local
		false,                   // no wait
		nil,                     // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to start message consuming")
	}

	return messages, nil
}

// CancelConsuming ...
func (client Client) CancelConsuming(queue string) error {
	if !client.queues.Contains(queue) {
		return errUnknownQueue
	}

	err := client.channel.Cancel(
		makeConsumerName(queue), // consumer name
		false,                   // no wait
	)
	if err != nil {
		return errors.Wrap(err, "unable to cancel message consuming")
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

func makeConsumerName(queue string) string {
	return queue + "_consumer"
}
