package rabbitmqutils

import (
	"testing"
	"testing/iotest"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ConsumeMessages(test *testing.T) {
	type fields struct {
		channel MessageBrokerChannel
	}
	type args struct {
		queue string
	}

	for _, data := range []struct {
		name           string
		fields         fields
		args           args
		wantedMessages []amqp.Delivery
		wantedErr      assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				channel: func() MessageBrokerChannel {
					messagesAsSlice := []amqp.Delivery{
						{Body: []byte("one")},
						{Body: []byte("two")},
					}
					messages := make(chan amqp.Delivery, len(messagesAsSlice))
					for _, message := range messagesAsSlice {
						messages <- message
					}
					close(messages)

					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Consume",
							"test",          // queue name
							"test_consumer", // consumer name
							false,           // auto-acknowledge
							false,           // exclusive
							false,           // no local
							false,           // no wait
							amqp.Table(nil), // arguments
						).
						Return((<-chan amqp.Delivery)(messages), nil)

					return channel
				}(),
			},
			args: args{
				queue: "test",
			},
			wantedMessages: []amqp.Delivery{
				{Body: []byte("one")},
				{Body: []byte("two")},
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Consume",
							"test",          // queue name
							"test_consumer", // consumer name
							false,           // auto-acknowledge
							false,           // exclusive
							false,           // no local
							false,           // no wait
							amqp.Table(nil), // arguments
						).
						Return(nil, iotest.ErrTimeout)

					return channel
				}(),
			},
			args: args{
				queue: "test",
			},
			wantedMessages: nil,
			wantedErr:      assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := Client{
				channel: data.fields.channel,
			}
			receivedMessagesChannel, receivedErr :=
				client.ConsumeMessages(data.args.queue)

			var receivedMessages []amqp.Delivery
			if receivedMessagesChannel != nil {
				for message := range receivedMessagesChannel {
					receivedMessages = append(receivedMessages, message)
				}
			}

			mock.AssertExpectationsForObjects(test, data.fields.channel)
			assert.Equal(test, data.wantedMessages, receivedMessages)
			data.wantedErr(test, receivedErr)
		})
	}
}

func TestClient_CancelConsuming(test *testing.T) {
	type fields struct {
		channel MessageBrokerChannel
	}
	type args struct {
		queue string
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
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Cancel",
							"test_consumer", // consumer name
							false,           // no wait
						).
						Return(nil)

					return channel
				}(),
			},
			args: args{
				queue: "test",
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Cancel",
							"test_consumer", // consumer name
							false,           // no wait
						).
						Return(iotest.ErrTimeout)

					return channel
				}(),
			},
			args: args{
				queue: "test",
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := Client{
				channel: data.fields.channel,
			}
			receivedErr := client.CancelConsuming(data.args.queue)

			mock.AssertExpectationsForObjects(test, data.fields.channel)
			data.wantedErr(test, receivedErr)
		})
	}
}
