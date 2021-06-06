package rabbitmqutils

import (
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
)

func TestWithDialer(test *testing.T) {
	var clientConfig ClientConfig
	option := WithDialer(func(dsn string) (MessageBrokerConnection, error) {
		panic("it should not be called")
	})
	option(&clientConfig)

	assert.NotNil(test, clientConfig.dialer)
}

func TestWithMaximalQueueSize(test *testing.T) {
	var clientConfig ClientConfig
	option := WithMaximalQueueSize(23)
	option(&clientConfig)

	assert.Equal(test, 23, clientConfig.maximalQueueSize)
}

func TestWithQueues(test *testing.T) {
	var clientConfig ClientConfig
	option := WithQueues([]string{"one", "two"})
	option(&clientConfig)

	assert.Equal(test, mapset.NewSet("one", "two"), clientConfig.queues)
}

func TestWithIDGenerator(test *testing.T) {
	var clientConfig ClientConfig
	option := WithIDGenerator(func() (string, error) {
		panic("it should not be called")
	})
	option(&clientConfig)

	assert.NotNil(test, clientConfig.idGenerator)
}

func TestWithClock(test *testing.T) {
	var clientConfig ClientConfig
	option := WithClock(func() time.Time {
		panic("it should not be called")
	})
	option(&clientConfig)

	assert.NotNil(test, clientConfig.clock)
}
