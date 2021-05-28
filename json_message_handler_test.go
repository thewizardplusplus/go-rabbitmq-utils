package rabbitmqutils

import (
	"reflect"
	"testing"
	"testing/iotest"

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
		{
			name: "success",
			fields: fields{
				MessageHandler: func() SpecificMessageHandler {
					messageHandler := new(MockSpecificMessageHandler)
					messageHandler.On("MessageType").Return(reflect.TypeOf(testMessage{}))
					messageHandler.
						On("HandleMessage", testMessage{FieldOne: 23, FieldTwo: "one"}).
						Return(nil)

					return messageHandler
				}(),
			},
			args: args{
				message: amqp.Delivery{Body: []byte(`{"FieldOne":23,"FieldTwo":"one"}`)},
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with unmarshalling",
			fields: fields{
				MessageHandler: func() SpecificMessageHandler {
					messageHandler := new(MockSpecificMessageHandler)
					messageHandler.On("MessageType").Return(reflect.TypeOf(func() {}))
					messageHandler.
						On("HandleMessage", testMessage{FieldOne: 23, FieldTwo: "one"}).
						Return(nil)

					return messageHandler
				}(),
			},
			args: args{
				message: amqp.Delivery{Body: []byte(`{"FieldOne":23,"FieldTwo":"one"}`)},
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with handling",
			fields: fields{
				MessageHandler: func() SpecificMessageHandler {
					messageHandler := new(MockSpecificMessageHandler)
					messageHandler.On("MessageType").Return(reflect.TypeOf(testMessage{}))
					messageHandler.
						On("HandleMessage", testMessage{FieldOne: 23, FieldTwo: "one"}).
						Return(iotest.ErrTimeout)

					return messageHandler
				}(),
			},
			args: args{
				message: amqp.Delivery{Body: []byte(`{"FieldOne":23,"FieldTwo":"one"}`)},
			},
			wantedErr: assert.Error,
		},
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
