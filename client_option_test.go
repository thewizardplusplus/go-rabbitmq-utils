package rabbitmqutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMaximalQueueSize(test *testing.T) {
	var clientOptions ClientOptions
	clientOption := WithMaximalQueueSize(23)
	clientOption(&clientOptions)

	assert.Equal(test, 23, clientOptions.maximalQueueSize)
}
