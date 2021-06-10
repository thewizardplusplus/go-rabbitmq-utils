package rabbitmqutils

import (
	"testing"
	"testing/iotest"
	"time"

	mapset "github.com/deckarep/golang-set"
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

				assert.Equal(test, mapset.NewSet(), client.queues)
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

				assert.Equal(test, mapset.NewSet(), client.queues)
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

				assert.Equal(test, mapset.NewSet("one", "two"), client.queues)
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

				assert.Equal(test, mapset.NewSet(), client.queues)
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

				assert.Equal(test, mapset.NewSet(), client.queues)
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
					client.queues,
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
					client.queues,
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
					client.queues,
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
							// queue name
							mock.MatchedBy(func(queue string) bool {
								return queue == "one" || queue == "two"
							}),
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
					client.queues,
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
		queues      mapset.Set
		idGenerator IDGeneratorInterface
		clock       ClockInterface
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

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantedErr assert.ErrorAssertionFunc
	}{
		{
			name: "success without message ID generating",
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
								Timestamp:   time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
								ContentType: "application/json",
								Body:        []byte(`{"FieldOne":23,"FieldTwo":"two"}`),
							},
						).
						Return(nil)

					return channel
				}(),
				queues: mapset.NewSet("one"),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("message-id", nil)

					return idGenerator
				}(),
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
				messageID:   "",
				messageData: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedErr: assert.NoError,
		},
		{
			name: "success with message ID generating",
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
								Timestamp:   time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
								ContentType: "application/json",
								Body:        []byte(`{"FieldOne":23,"FieldTwo":"two"}`),
							},
						).
						Return(nil)

					return channel
				}(),
				queues:      mapset.NewSet("one"),
				idGenerator: new(MockIDGeneratorInterface),
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
			wantedErr: assert.NoError,
		},
		{
			name: "error with an unknown queue",
			fields: fields{
				channel:     new(MockMessageBrokerChannel),
				queues:      mapset.NewSet("one"),
				idGenerator: new(MockIDGeneratorInterface),
				clock:       new(MockClockInterface),
			},
			args: args{
				queue:       "unknown",
				messageID:   "",
				messageData: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedErr: func(
				test assert.TestingT,
				err error,
				msgAndArgs ...interface{},
			) bool {
				return assert.Equal(test, errUnknownQueue, err, msgAndArgs...)
			},
		},
		{
			name: "error with message ID generating",
			fields: fields{
				channel: new(MockMessageBrokerChannel),
				queues:  mapset.NewSet("one"),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("", iotest.ErrTimeout)

					return idGenerator
				}(),
				clock: new(MockClockInterface),
			},
			args: args{
				queue:       "one",
				messageID:   "",
				messageData: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with marshalling",
			fields: fields{
				channel: new(MockMessageBrokerChannel),
				queues:  mapset.NewSet("one"),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("message-id", nil)

					return idGenerator
				}(),
				clock: new(MockClockInterface),
			},
			args: args{
				queue:       "one",
				messageID:   "",
				messageData: func() {},
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
								Timestamp:   time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
								ContentType: "application/json",
								Body:        []byte(`{"FieldOne":23,"FieldTwo":"two"}`),
							},
						).
						Return(iotest.ErrTimeout)

					return channel
				}(),
				queues: mapset.NewSet("one"),
				idGenerator: func() IDGeneratorInterface {
					idGenerator := new(MockIDGeneratorInterface)
					idGenerator.On("GenerateID").Return("message-id", nil)

					return idGenerator
				}(),
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
				messageID:   "",
				messageData: testMessage{FieldOne: 23, FieldTwo: "two"},
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := Client{
				channel:     data.fields.channel,
				queues:      data.fields.queues,
				idGenerator: data.fields.idGenerator.GenerateID,
				clock:       data.fields.clock.Time,
			}
			receivedErr := client.PublishMessage(
				data.args.queue,
				data.args.messageID,
				data.args.messageData,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.channel,
				data.fields.idGenerator,
				data.fields.clock,
			)
			data.wantedErr(test, receivedErr)
		})
	}
}

func TestClient_ConsumeMessages(test *testing.T) {
	type fields struct {
		channel MessageBrokerChannel
		queues  mapset.Set
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
				queues: mapset.NewSet("test"),
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
			name: "error with an unknown queue",
			fields: fields{
				channel: new(MockMessageBrokerChannel),
				queues:  mapset.NewSet("test"),
			},
			args: args{
				queue: "unknown",
			},
			wantedMessages: nil,
			wantedErr: func(
				test assert.TestingT,
				err error,
				msgAndArgs ...interface{},
			) bool {
				return assert.Equal(test, errUnknownQueue, err, msgAndArgs...)
			},
		},
		{
			name: "error with starting of message consuming",
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
				queues: mapset.NewSet("test"),
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
				queues:  data.fields.queues,
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
		queues  mapset.Set
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
				queues: mapset.NewSet("test"),
			},
			args: args{
				queue: "test",
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with an unknown queue",
			fields: fields{
				channel: new(MockMessageBrokerChannel),
				queues:  mapset.NewSet("test"),
			},
			args: args{
				queue: "unknown",
			},
			wantedErr: func(
				test assert.TestingT,
				err error,
				msgAndArgs ...interface{},
			) bool {
				return assert.Equal(test, errUnknownQueue, err, msgAndArgs...)
			},
		},
		{
			name: "error with cancelling of message consuming",
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
				queues: mapset.NewSet("test"),
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
				queues:  data.fields.queues,
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
