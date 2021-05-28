package rabbitmqutils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestJSONMessageHandler_HandleMessage(test *testing.T) {
	type fields struct {
		MessageHandler SpecificMessageHandler
	}
	type args struct {
		message amqp.Delivery
	}
	type testMessage struct {
		FieldOne int
		FieldTwo string
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantedErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := JSONMessageHandler{
				MessageHandler: data.fields.MessageHandler,
			}
			receivedErr := handler.HandleMessage(data.args.message)

			data.wantedErr(test, receivedErr)
		})
	}
}
