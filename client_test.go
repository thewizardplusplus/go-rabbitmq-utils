package rabbitmqutils

import (
	"testing"

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
		// TODO: Add test cases.
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
