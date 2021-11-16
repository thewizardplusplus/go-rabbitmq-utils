package rabbitmqutils

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewMessageConsumer(test *testing.T) {
	type args struct {
		client         MessageConsumerClient
		queue          string
		messageHandler MessageHandler
	}

	for _, data := range []struct {
		name                  string
		args                  args
		wantedMessageConsumer func(test *testing.T, messageConsumer MessageConsumer)
		wantedErr             assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				client: func() MessageConsumerClient {
					client := new(MockMessageConsumerClient)
					client.
						On("ConsumeMessages", "test").
						Return(make(<-chan amqp.Delivery), nil)

					return client
				}(),
				queue:          "test",
				messageHandler: new(MockMessageHandler),
			},
			wantedMessageConsumer: func(
				test *testing.T,
				messageConsumer MessageConsumer,
			) {
				assert.IsType(test, new(MockMessageConsumerClient), messageConsumer.client)
				assert.Equal(test, "test", messageConsumer.queue)

				for _, field := range []interface{}{
					messageConsumer.client,
					messageConsumer.stoppingCtx,
				} {
					assert.NotNil(test, field)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				client: func() MessageConsumerClient {
					client := new(MockMessageConsumerClient)
					client.On("ConsumeMessages", "test").Return(nil, iotest.ErrTimeout)

					return client
				}(),
				queue:          "test",
				messageHandler: new(MockMessageHandler),
			},
			wantedMessageConsumer: func(
				test *testing.T,
				messageConsumer MessageConsumer,
			) {
				assert.Zero(test, messageConsumer)
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			receivedMessageConsumer, receivedErr := NewMessageConsumer(
				data.args.client,
				data.args.queue,
				data.args.messageHandler,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.client,
				data.args.messageHandler,
			)
			data.wantedMessageConsumer(test, receivedMessageConsumer)
			data.wantedErr(test, receivedErr)
		})
	}
}

func TestMessageConsumer_starting(test *testing.T) {
	type args struct {
		client         MessageConsumerClient
		queue          string
		messageHandler MessageHandler
	}

	for _, data := range []struct {
		name          string
		args          args
		startConsumer func(consumer MessageConsumer)
	}{
		{
			name: "with the Start() method",
			args: args{
				client: func() MessageConsumerClient {
					messagesAsSlice := []amqp.Delivery{
						{Body: []byte("one")},
						{Body: []byte("two")},
					}
					messages := make(chan amqp.Delivery, len(messagesAsSlice))
					for _, message := range messagesAsSlice {
						messages <- message
					}
					close(messages)

					client := new(MockMessageConsumerClient)
					client.
						On("ConsumeMessages", "test").
						Return((<-chan amqp.Delivery)(messages), nil)
					client.On("CancelConsuming", "test").Return(nil)

					return client
				}(),
				queue: "test",
				messageHandler: func() MessageHandler {
					messageHandler := new(MockMessageHandler)
					messageHandler.
						On("HandleMessage", amqp.Delivery{Body: []byte("one")}).
						Return()
					messageHandler.
						On("HandleMessage", amqp.Delivery{Body: []byte("two")}).
						Return()

					return messageHandler
				}(),
			},
			startConsumer: func(consumer MessageConsumer) {
				consumer.Start()
			},
		},
		{
			name: "with the StartConcurrently() method",
			args: args{
				client: func() MessageConsumerClient {
					messagesAsSlice := []amqp.Delivery{
						{Body: []byte("one")},
						{Body: []byte("two")},
					}
					messages := make(chan amqp.Delivery, len(messagesAsSlice))
					for _, message := range messagesAsSlice {
						messages <- message
					}
					close(messages)

					client := new(MockMessageConsumerClient)
					client.
						On("ConsumeMessages", "test").
						Return((<-chan amqp.Delivery)(messages), nil)
					client.On("CancelConsuming", "test").Return(nil)

					return client
				}(),
				queue: "test",
				messageHandler: func() MessageHandler {
					messageHandler := new(MockMessageHandler)
					messageHandler.
						On("HandleMessage", amqp.Delivery{Body: []byte("one")}).
						Return()
					messageHandler.
						On("HandleMessage", amqp.Delivery{Body: []byte("two")}).
						Return()

					return messageHandler
				}(),
			},
			startConsumer: func(consumer MessageConsumer) {
				consumer.StartConcurrently(10)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			messageConsumer, err := NewMessageConsumer(
				data.args.client,
				data.args.queue,
				data.args.messageHandler,
			)
			require.NoError(test, err)

			go data.startConsumer(messageConsumer)
			messageConsumer.Stop() // nolint: errcheck

			mock.AssertExpectationsForObjects(
				test,
				data.args.client,
				data.args.messageHandler,
			)
		})
	}
}

func TestMessageConsumer_Stop(test *testing.T) {
	type fields struct {
		stoppingCtx context.Context
	}
	type args struct {
		client         MessageConsumerClient
		queue          string
		messageHandler MessageHandler
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantedErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				stoppingCtx: func() context.Context {
					ctx, canceller := context.WithCancel(context.Background())
					canceller()

					return ctx
				}(),
			},
			args: args{
				client: func() MessageConsumerClient {
					messages := make(chan amqp.Delivery)
					close(messages)

					client := new(MockMessageConsumerClient)
					client.
						On("ConsumeMessages", "test").
						Return((<-chan amqp.Delivery)(messages), nil)
					client.On("CancelConsuming", "test").Return(nil)

					return client
				}(),
				queue:          "test",
				messageHandler: new(MockMessageHandler),
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				stoppingCtx: context.Background(),
			},
			args: args{
				client: func() MessageConsumerClient {
					messages := make(chan amqp.Delivery)
					close(messages)

					client := new(MockMessageConsumerClient)
					client.
						On("ConsumeMessages", "test").
						Return((<-chan amqp.Delivery)(messages), nil)
					client.On("CancelConsuming", "test").Return(iotest.ErrTimeout)

					return client
				}(),
				queue:          "test",
				messageHandler: new(MockMessageHandler),
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			messageConsumer, err := NewMessageConsumer(
				data.args.client,
				data.args.queue,
				data.args.messageHandler,
			)
			require.NoError(test, err)
			messageConsumer.stoppingCtx = data.fields.stoppingCtx

			go messageConsumer.Start()
			receivedErr := messageConsumer.Stop()

			mock.AssertExpectationsForObjects(
				test,
				data.args.client,
				data.args.messageHandler,
			)
			data.wantedErr(test, receivedErr)
		})
	}
}
