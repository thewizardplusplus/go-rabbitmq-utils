// +build integration

package rabbitmqutils

import (
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_PublishMessage_integration(test *testing.T) {
	type fields struct {
		clock ClockInterface
	}
	type args struct {
		queue       string
		messageID   string
		messageData interface{}
	}
	type testMessage struct {
		FieldOne int
		FieldTwo string
	}

	dsn, ok := os.LookupEnv("MESSAGE_BROKER_ADDRESS")
	if !ok {
		dsn = "amqp://rabbitmq:rabbitmq@localhost:5672"
	}

	for _, data := range []struct {
		name          string
		fields        fields
		args          args
		wantedMessage amqp.Delivery
	}{
		{
			name: "success",
			fields: fields{
				clock: func() ClockInterface {
					clock := new(MockClockInterface)
					clock.
						On("Time").
						Return(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))

					return clock
				}(),
			},
			args: args{
				queue:       "one",
				messageID:   "message-id",
				messageData: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedMessage: amqp.Delivery{
				MessageId:   "message-id",
				Timestamp:   time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				ContentType: "application/json",
				Body:        []byte(`{"FieldOne":23,"FieldTwo":"two"}`),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			// prepare the client
			client, err := NewClient(
				dsn,
				WithQueues([]string{data.args.queue}),
				WithClock(data.fields.clock.Time),
			)
			require.NoError(test, err)
			defer client.Close()

			// publish the message
			err = client.PublishMessage(
				data.args.queue,
				data.args.messageID,
				data.args.messageData,
			)
			require.NoError(test, err)

			// receive the message
			connection, err := amqp.Dial(dsn)
			require.NoError(test, err)
			defer connection.Close()

			channel, err := connection.Channel()
			require.NoError(test, err)
			defer channel.Close()

			receivedMessage, _, err := channel.Get(
				data.args.queue, // queue name
				true,            // auto-acknowledge
			)
			require.NoError(test, err)

			// clean the irrelevant message fields
			receivedMessage = amqp.Delivery{
				MessageId:   receivedMessage.MessageId,
				Timestamp:   receivedMessage.Timestamp.In(time.UTC),
				ContentType: receivedMessage.ContentType,
				Body:        receivedMessage.Body,
			}

			// check the results
			mock.AssertExpectationsForObjects(test, data.fields.clock)
			assert.Equal(test, data.wantedMessage, receivedMessage)
		})
	}
}

func TestClient_ConsumeMessages_integration(test *testing.T) {
	type fields struct {
		clock          ClockInterface
		messageHandler MessageHandler
	}
	type messageArgs struct {
		messageID   string
		messageData interface{}
	}
	type args struct {
		queue    string
		messages []messageArgs
	}

	dsn, ok := os.LookupEnv("MESSAGE_BROKER_ADDRESS")
	if !ok {
		dsn = "amqp://rabbitmq:rabbitmq@localhost:5672"
	}

	var waitGroupInstance *sync.WaitGroup
	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			waitGroupInstance = new(sync.WaitGroup)

			// prepare the client
			client, err := NewClient(
				dsn,
				WithQueues([]string{data.args.queue}),
				WithClock(data.fields.clock.Time),
			)
			require.NoError(test, err)
			defer client.Close()

			// start the message consuming
			messageConsumer, err := NewMessageConsumer(
				client,
				data.args.queue,
				data.fields.messageHandler,
			)
			go messageConsumer.StartConcurrently(runtime.NumCPU())
			defer messageConsumer.Stop()
			require.NoError(test, err)

			// publish the messages
			for _, message := range data.args.messages {
				waitGroupInstance.Add(1)

				err = client.PublishMessage(
					data.args.queue,
					message.messageID,
					message.messageData,
				)
				require.NoError(test, err)
			}
			waitGroupInstance.Wait()

			// check the results
			mock.AssertExpectationsForObjects(
				test,
				data.fields.clock,
				data.fields.messageHandler,
			)
		})
	}
}
