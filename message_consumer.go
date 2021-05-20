package rabbitmqutils

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// MessageHandler ...
type MessageHandler interface {
	HandleMessage(message amqp.Delivery)
}

// MessageConsumer ...
type MessageConsumer struct {
	client               Client
	queue                string
	messages             <-chan amqp.Delivery
	messageHandler       MessageHandler
	stoppingCtx          context.Context
	stoppingCtxCanceller context.CancelFunc
}

// NewMessageConsumer ...
func NewMessageConsumer(
	client Client,
	queue string,
	messageHandler MessageHandler,
) (MessageConsumer, error) {
	messages, err := client.ConsumeMessages(queue)
	if err != nil {
		return MessageConsumer{}, errors.Wrap(err, "unable to start the consuming")
	}

	stoppingCtx, stoppingCtxCanceller := context.WithCancel(context.Background())
	messageConsumer := MessageConsumer{
		client:               client,
		queue:                queue,
		messages:             messages,
		messageHandler:       messageHandler,
		stoppingCtx:          stoppingCtx,
		stoppingCtxCanceller: stoppingCtxCanceller,
	}
	return messageConsumer, nil
}

// Start ...
func (consumer MessageConsumer) Start() {
	for message := range consumer.messages {
		consumer.messageHandler.HandleMessage(message)
	}
}

// StartConcurrently ...
func (consumer MessageConsumer) StartConcurrently(concurrency int) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer waitGroup.Done()

			consumer.Start()
		}()
	}

	waitGroup.Wait()
	consumer.stoppingCtxCanceller()
}

// Stop ...
func (consumer MessageConsumer) Stop() error {
	if err := consumer.client.CancelConsuming(consumer.queue); err != nil {
		return errors.Wrap(err, "unable to cancel the consuming")
	}

	<-consumer.stoppingCtx.Done()
	return nil
}