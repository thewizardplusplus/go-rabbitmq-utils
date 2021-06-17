// +build integration

package rabbitmqutils

import (
	"os"
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
		// TODO: Add test cases.
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
