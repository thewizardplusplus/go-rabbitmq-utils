package rabbitmqutils

import (
	"context"
	"fmt"
	"testing"
	"testing/iotest"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				assert.Len(test, messageConsumer.messages, 0)
				assert.Equal(test, new(MockMessageHandler), messageConsumer.messageHandler)
				assert.Equal(
					test,
					&StartModeHolder{mode: NotStarted},
					messageConsumer.startMode,
				)

				for _, field := range []interface{}{
					messageConsumer.client,
					messageConsumer.messages,
					messageConsumer.stoppingCtx,
					messageConsumer.stoppingCtxCanceller,
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

func TestMessageConsumer_Start(test *testing.T) {
	type fields struct {
		messages             <-chan amqp.Delivery
		messageHandler       MessageHandler
		startMode            *StartModeHolder
		stoppingCtxCanceller ContextCancellerInterface
	}

	for _, data := range []struct {
		name   string
		fields fields
	}{
		{
			name: "success with the messages",
			fields: fields{
				messages: func() <-chan amqp.Delivery {
					messagesAsSlice := []amqp.Delivery{
						{Body: []byte("one")},
						{Body: []byte("two")},
					}
					messages := make(chan amqp.Delivery, len(messagesAsSlice))
					for _, message := range messagesAsSlice {
						messages <- message
					}

					close(messages)
					return messages
				}(),
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
				startMode: &StartModeHolder{mode: NotStarted},
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return()

					return stoppingCtxCanceller
				}(),
			},
		},
		{
			name: "success without messages",
			fields: fields{
				messages: func() <-chan amqp.Delivery {
					messages := make(chan amqp.Delivery)
					close(messages)

					return messages
				}(),
				messageHandler: new(MockMessageHandler),
				startMode:      &StartModeHolder{mode: NotStarted},
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return()

					return stoppingCtxCanceller
				}(),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			consumer := MessageConsumer{
				messages:             data.fields.messages,
				messageHandler:       data.fields.messageHandler,
				startMode:            data.fields.startMode,
				stoppingCtxCanceller: data.fields.stoppingCtxCanceller.CancelContext,
			}
			consumer.Start()

			mock.AssertExpectationsForObjects(
				test,
				data.fields.messageHandler,
				data.fields.stoppingCtxCanceller,
			)
		})
	}
}

func TestMessageConsumer_StartConcurrently(test *testing.T) {
	type fields struct {
		messages             <-chan amqp.Delivery
		messageHandler       MessageHandler
		stoppingCtxCanceller ContextCancellerInterface
	}
	type args struct {
		concurrency int
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success with the messages and without concurrency",
			fields: fields{
				messages: func() <-chan amqp.Delivery {
					var messagesAsSlice []amqp.Delivery
					for i := 0; i < 100; i++ {
						messagesAsSlice = append(messagesAsSlice, amqp.Delivery{
							Body: []byte(fmt.Sprintf("message #%d", i)),
						})
					}

					messages := make(chan amqp.Delivery, len(messagesAsSlice))
					for _, message := range messagesAsSlice {
						messages <- message
					}

					close(messages)
					return messages
				}(),
				messageHandler: func() MessageHandler {
					messageHandler := new(MockMessageHandler)
					for i := 0; i < 100; i++ {
						messageHandler.
							On("HandleMessage", amqp.Delivery{
								Body: []byte(fmt.Sprintf("message #%d", i)),
							}).
							Return()
					}

					return messageHandler
				}(),
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return()

					return stoppingCtxCanceller
				}(),
			},
			args: args{
				concurrency: 1,
			},
		},
		{
			name: "success with the messages and concurrency",
			fields: fields{
				messages: func() <-chan amqp.Delivery {
					var messagesAsSlice []amqp.Delivery
					for i := 0; i < 100; i++ {
						messagesAsSlice = append(messagesAsSlice, amqp.Delivery{
							Body: []byte(fmt.Sprintf("message #%d", i)),
						})
					}

					messages := make(chan amqp.Delivery, len(messagesAsSlice))
					for _, message := range messagesAsSlice {
						messages <- message
					}

					close(messages)
					return messages
				}(),
				messageHandler: func() MessageHandler {
					messageHandler := new(MockMessageHandler)
					for i := 0; i < 100; i++ {
						messageHandler.
							On("HandleMessage", amqp.Delivery{
								Body: []byte(fmt.Sprintf("message #%d", i)),
							}).
							Return()
					}

					return messageHandler
				}(),
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return()

					return stoppingCtxCanceller
				}(),
			},
			args: args{
				concurrency: 5,
			},
		},
		{
			name: "success without messages and concurrency",
			fields: fields{
				messages: func() <-chan amqp.Delivery {
					messages := make(chan amqp.Delivery)
					close(messages)

					return messages
				}(),
				messageHandler: new(MockMessageHandler),
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return()

					return stoppingCtxCanceller
				}(),
			},
			args: args{
				concurrency: 1,
			},
		},
		{
			name: "success without messages and with concurrency",
			fields: fields{
				messages: func() <-chan amqp.Delivery {
					messages := make(chan amqp.Delivery)
					close(messages)

					return messages
				}(),
				messageHandler: new(MockMessageHandler),
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return()

					return stoppingCtxCanceller
				}(),
			},
			args: args{
				concurrency: 5,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			consumer := MessageConsumer{
				messages:             data.fields.messages,
				messageHandler:       data.fields.messageHandler,
				stoppingCtxCanceller: data.fields.stoppingCtxCanceller.CancelContext,
			}
			consumer.StartConcurrently(data.args.concurrency)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.messageHandler,
				data.fields.stoppingCtxCanceller,
			)
		})
	}
}

func TestMessageConsumer_Stop(test *testing.T) {
	type fields struct {
		client      MessageConsumerClient
		queue       string
		stoppingCtx context.Context
	}

	for _, data := range []struct {
		name      string
		fields    fields
		wantedErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				client: func() MessageConsumerClient {
					client := new(MockMessageConsumerClient)
					client.On("CancelConsuming", "test").Return(nil)

					return client
				}(),
				queue: "test",
				stoppingCtx: func() context.Context {
					ctx, canceller := context.WithCancel(context.Background())
					canceller()

					return ctx
				}(),
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				client: func() MessageConsumerClient {
					client := new(MockMessageConsumerClient)
					client.On("CancelConsuming", "test").Return(iotest.ErrTimeout)

					return client
				}(),
				queue:       "test",
				stoppingCtx: context.Background(),
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			consumer := MessageConsumer{
				client:      data.fields.client,
				queue:       data.fields.queue,
				stoppingCtx: data.fields.stoppingCtx,
			}
			receivedErr := consumer.Stop()

			mock.AssertExpectationsForObjects(test, data.fields.client)
			data.wantedErr(test, receivedErr)
		})
	}
}
