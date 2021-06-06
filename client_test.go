package rabbitmqutils

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// nolint: gocyclo
func TestNewClient(test *testing.T) {
	type args struct {
		dsn     string
		dialer  DialerInterface
		options []ClientOption
	}

	for _, data := range []struct {
		name         string
		args         args
		wantedClient func(test *testing.T, client Client)
		wantedErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success without options",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					channel := new(MockMessageBrokerChannel)

					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(channel, nil)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: nil,
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.NotNil(test, field)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "success with the setting of a maximal queue size",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Qos",
							23,    // prefetch count
							0,     // prefetch size
							false, // global
						).
						Return(nil)

					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(channel, nil)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: []ClientOption{WithMaximalQueueSize(23)},
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.NotNil(test, field)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "success with the declaring of queues",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					channel := new(MockMessageBrokerChannel)
					for _, queue := range []string{"one", "two"} {
						channel.
							On(
								"QueueDeclare",
								queue,           // queue name
								true,            // durable
								false,           // auto-delete
								false,           // exclusive
								false,           // no wait
								amqp.Table(nil), // arguments
							).
							Return(amqp.Queue{}, nil)
					}

					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(channel, nil)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: []ClientOption{WithQueues([]string{"one", "two"})},
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.NotNil(test, field)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "success with the setting of an ID generator",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					channel := new(MockMessageBrokerChannel)

					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(channel, nil)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: []ClientOption{
					WithIDGenerator(func() (string, error) {
						panic("it should not be called")
					}),
				},
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.NotNil(test, field)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "success with the setting of a clock",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					channel := new(MockMessageBrokerChannel)

					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(channel, nil)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: []ClientOption{
					WithClock(func() time.Time {
						panic("it should not be called")
					}),
				},
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.NotNil(test, field)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with the opening of the connection",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(nil, iotest.ErrTimeout)

					return dialer
				}(),
				options: nil,
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.Nil(test, field)
				}
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with the opening of the channel",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(nil, iotest.ErrTimeout)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: nil,
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.Nil(test, field)
				}
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with the setting of a maximal queue size",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Qos",
							23,    // prefetch count
							0,     // prefetch size
							false, // global
						).
						Return(iotest.ErrTimeout)

					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(channel, nil)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: []ClientOption{WithMaximalQueueSize(23)},
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.Nil(test, field)
				}
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with the declaring of queues",
			args: args{
				dsn: "dsn",
				dialer: func() DialerInterface {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"QueueDeclare",
							"one",           // queue name
							true,            // durable
							false,           // auto-delete
							false,           // exclusive
							false,           // no wait
							amqp.Table(nil), // arguments
						).
						Return(amqp.Queue{}, iotest.ErrTimeout)

					connection := new(MockMessageBrokerConnection)
					connection.On("Channel").Return(channel, nil)

					dialer := new(MockDialerInterface)
					dialer.On("Dial", "dsn").Return(connection, nil)

					return dialer
				}(),
				options: []ClientOption{WithQueues([]string{"one", "two"})},
			},
			wantedClient: func(test *testing.T, client Client) {
				for _, field := range []interface{}{
					client.connection,
					client.channel,
					client.idGenerator,
					client.clock,
				} {
					assert.Nil(test, field)
				}
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			options := append(data.args.options, WithDialer(data.args.dialer.Dial))
			receivedClient, receivedErr := NewClient(data.args.dsn, options...)

			mock.AssertExpectationsForObjects(test, data.args.dialer)
			if receivedClient.connection != nil {
				mock.AssertExpectationsForObjects(test, receivedClient.connection)
			}
			if receivedClient.channel != nil {
				mock.AssertExpectationsForObjects(test, receivedClient.channel)
			}
			data.wantedClient(test, receivedClient)
			data.wantedErr(test, receivedErr)
		})
	}
}

func TestClient_PublishMessage(test *testing.T) {
	type fields struct {
		channel     MessageBrokerChannel
		idGenerator IDGeneratorInterface
	}
	type args struct {
		queue   string
		message interface{}
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
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Publish",
							"",    // exchange
							"one", // queue name
							false, // mandatory
							false, // immediate
							amqp.Publishing{
								MessageId:   "message-id",
								ContentType: "application/json",
								Body:        []byte(`{"FieldOne":23,"FieldTwo":"two"}`),
							},
						).
						Return(nil)

					return channel
				}(),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("message-id", nil)

					return idGenerator
				}(),
			},
			args: args{
				queue:   "one",
				message: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with message ID generating",
			fields: fields{
				channel: new(MockMessageBrokerChannel),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("", iotest.ErrTimeout)

					return idGenerator
				}(),
			},
			args: args{
				queue:   "one",
				message: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with marshalling",
			fields: fields{
				channel: new(MockMessageBrokerChannel),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("message-id", nil)

					return idGenerator
				}(),
			},
			args: args{
				queue:   "one",
				message: func() {},
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with publishing",
			fields: fields{
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.
						On(
							"Publish",
							"",    // exchange
							"one", // queue name
							false, // mandatory
							false, // immediate
							amqp.Publishing{
								MessageId:   "message-id",
								ContentType: "application/json",
								Body:        []byte(`{"FieldOne":23,"FieldTwo":"two"}`),
							},
						).
						Return(iotest.ErrTimeout)

					return channel
				}(),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("message-id", nil)

					return idGenerator
				}(),
			},
			args: args{
				queue:   "one",
				message: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := Client{
				channel:     data.fields.channel,
				idGenerator: data.fields.idGenerator.GenerateID,
			}
			receivedErr := client.PublishMessage(data.args.queue, data.args.message)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.channel,
				data.fields.idGenerator,
			)
			data.wantedErr(test, receivedErr)
		})
	}
}

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

func TestClient_Close(test *testing.T) {
	type fields struct {
		connection MessageBrokerConnection
		channel    MessageBrokerChannel
	}

	for _, data := range []struct {
		name      string
		fields    fields
		wantedErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				connection: func() MessageBrokerConnection {
					connection := new(MockMessageBrokerConnection)
					connection.On("Close").Return(nil)

					return connection
				}(),
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.On("Close").Return(nil)

					return channel
				}(),
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with the channel",
			fields: fields{
				connection: new(MockMessageBrokerConnection),
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.On("Close").Return(iotest.ErrTimeout)

					return channel
				}(),
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with the connection",
			fields: fields{
				connection: func() MessageBrokerConnection {
					connection := new(MockMessageBrokerConnection)
					connection.On("Close").Return(iotest.ErrTimeout)

					return connection
				}(),
				channel: func() MessageBrokerChannel {
					channel := new(MockMessageBrokerChannel)
					channel.On("Close").Return(nil)

					return channel
				}(),
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := Client{
				connection: data.fields.connection,
				channel:    data.fields.channel,
			}
			receivedErr := client.Close()

			mock.AssertExpectationsForObjects(
				test,
				data.fields.connection,
				data.fields.channel,
			)
			data.wantedErr(test, receivedErr)
		})
	}
}
