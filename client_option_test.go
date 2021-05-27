package rabbitmqutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
