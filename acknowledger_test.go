package rabbitmqutils

import (
	"reflect"
	"testing"
	"testing/iotest"

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
		{
			name: "success",
			fields: fields{
				MessageHandling: OnceMessageHandling,
				MessageHandler: func() FailingMessageHandler {
					messageHandler := new(MockFailingMessageHandler)
					messageHandler.
						On("HandleMessage", mock.MatchedBy(func(message amqp.Delivery) bool {
							message.Acknowledger = nil

							return reflect.DeepEqual(message, amqp.Delivery{
								DeliveryTag: 23,
								Body:        []byte("test"),
							})
						})).
						Return(nil)

					return messageHandler
				}(),
				Logger: func() log.Logger {
					logger := new(MockLogger)
					logger.On("Log", "message has been handled successfully").Return()

					return logger
				}(),
			},
			args: args{
				message: amqp.Delivery{
					Acknowledger: func() amqp.Acknowledger {
						amqpAcknowledger := new(MockAMQPAcknowledger)
						amqpAcknowledger.On("Ack", uint64(23), false).Return(nil)

						return amqpAcknowledger
					}(),
					DeliveryTag: 23,
					Body:        []byte("test"),
				},
			},
		},
		{
			name: "error with once message handling",
			fields: fields{
				MessageHandling: OnceMessageHandling,
				MessageHandler: func() FailingMessageHandler {
					messageHandler := new(MockFailingMessageHandler)
					messageHandler.
						On("HandleMessage", mock.MatchedBy(func(message amqp.Delivery) bool {
							message.Acknowledger = nil

							return reflect.DeepEqual(message, amqp.Delivery{
								DeliveryTag: 23,
								Body:        []byte("test"),
							})
						})).
						Return(iotest.ErrTimeout)

					return messageHandler
				}(),
				Logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On("Logf", "unable to handle the message: %v", iotest.ErrTimeout).
						Return()

					return logger
				}(),
			},
			args: args{
				message: amqp.Delivery{
					Acknowledger: func() amqp.Acknowledger {
						amqpAcknowledger := new(MockAMQPAcknowledger)
						amqpAcknowledger.On("Reject", uint64(23), false).Return(nil)

						return amqpAcknowledger
					}(),
					DeliveryTag: 23,
					Body:        []byte("test"),
				},
			},
		},
		{
			name: "error with twice message handling",
			fields: fields{
				MessageHandling: TwiceMessageHandling,
				MessageHandler: func() FailingMessageHandler {
					messageHandler := new(MockFailingMessageHandler)
					messageHandler.
						On("HandleMessage", mock.MatchedBy(func(message amqp.Delivery) bool {
							message.Acknowledger = nil

							return reflect.DeepEqual(message, amqp.Delivery{
								DeliveryTag: 23,
								Body:        []byte("test"),
							})
						})).
						Return(iotest.ErrTimeout)

					return messageHandler
				}(),
				Logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On("Logf", "unable to handle the message: %v", iotest.ErrTimeout).
						Return()

					return logger
				}(),
			},
			args: args{
				message: amqp.Delivery{
					Acknowledger: func() amqp.Acknowledger {
						amqpAcknowledger := new(MockAMQPAcknowledger)
						amqpAcknowledger.On("Reject", uint64(23), true).Return(nil)

						return amqpAcknowledger
					}(),
					DeliveryTag: 23,
					Body:        []byte("test"),
				},
			},
		},
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
