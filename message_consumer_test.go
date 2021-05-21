package rabbitmqutils

import (
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
