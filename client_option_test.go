package rabbitmqutils

import (
	"testing"

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

	assert.Equal(test, []string{"one", "two"}, clientConfig.queues)
}
