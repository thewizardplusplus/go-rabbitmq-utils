package rabbitmqutils

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

//go:generate mockery --name=MessageConsumerClient --inpackage --case=underscore --testonly

// MessageConsumerClient ...
type MessageConsumerClient interface {
	ConsumeMessages(queue string) (<-chan amqp.Delivery, error)
	CancelConsuming(queue string) error
}

//go:generate mockery --name=MessageHandler --inpackage --case=underscore --testonly

// MessageHandler ...
type MessageHandler interface {
	HandleMessage(message amqp.Delivery)
}

type messageHandlerWrapper struct {
	messageHandler MessageHandler
}

func (wrapper messageHandlerWrapper) Handle(
	ctx context.Context,
	data interface{},
) {
	wrapper.messageHandler.HandleMessage(data.(amqp.Delivery))
}

// MessageConsumer ...
type MessageConsumer struct {
	// do not use embedding to hide the Handle() method
	concurrentHandler syncutils.ConcurrentHandler
	client            MessageConsumerClient
	queue             string
	stoppingCtx       context.Context
}

// NewMessageConsumer ...
func NewMessageConsumer(
	client MessageConsumerClient,
	queue string,
	messageHandler MessageHandler,
) (MessageConsumer, error) {
	messages, err := client.ConsumeMessages(queue)
	if err != nil {
		return MessageConsumer{}, errors.Wrap(err, "unable to start the consuming")
	}

	concurrentHandler := syncutils.NewConcurrentHandler(0, messageHandlerWrapper{
		messageHandler: messageHandler,
	})
	stoppingCtx, stoppingCtxCanceller := context.WithCancel(context.Background())
	go func() {
		defer stoppingCtxCanceller()

		for message := range messages {
			concurrentHandler.Handle(message)
		}
	}()

	messageConsumer := MessageConsumer{
		concurrentHandler: concurrentHandler,
		client:            client,
		queue:             queue,
		stoppingCtx:       stoppingCtx,
	}
	return messageConsumer, nil
}

// Start ...
func (consumer MessageConsumer) Start() {
	consumer.concurrentHandler.Start(context.Background())
}

// StartConcurrently ...
func (consumer MessageConsumer) StartConcurrently(concurrencyFactor int) {
	consumer.concurrentHandler.
		StartConcurrently(context.Background(), concurrencyFactor)
}

// Stop ...
func (consumer MessageConsumer) Stop() error {
	if err := consumer.client.CancelConsuming(consumer.queue); err != nil {
		return errors.Wrap(err, "unable to cancel the consuming")
	}
	<-consumer.stoppingCtx.Done()

	consumer.concurrentHandler.Stop()
	return nil
}
