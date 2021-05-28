package rabbitmqutils

import (
	"testing"

	"github.com/go-log/log"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

func TestAcknowledger_HandleMessage(test *testing.T) {
	type fields struct {
		MessageHandling MessageHandling
		MessageHandler  FailingMessageHandler
		Logger          log.Logger
	}
	type args struct {
		message amqp.Delivery
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			acknowledger := Acknowledger{
				MessageHandling: data.fields.MessageHandling,
				MessageHandler:  data.fields.MessageHandler,
				Logger:          data.fields.Logger,
			}
			acknowledger.HandleMessage(data.args.message)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.MessageHandler,
				data.fields.Logger,
				data.args.message.Acknowledger,
			)
		})
	}
}
