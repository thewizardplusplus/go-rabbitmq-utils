package rabbitmqutils

import (
	"time"
)

// Dialer ...
type Dialer func(dsn string) (MessageBrokerConnection, error)

// IDGenerator ...
type IDGenerator func() (string, error)

// Clock ...
type Clock func() time.Time

// ClientConfig ...
type ClientConfig struct {
	dialer           Dialer
	maximalQueueSize int
	queues           []string
	idGenerator      IDGenerator
}

// ClientOption ...
type ClientOption func(clientConfig *ClientConfig)

// WithDialer ...
func WithDialer(dialer Dialer) ClientOption {
	return func(clientConfig *ClientConfig) {
		clientConfig.dialer = dialer
	}
}

// WithMaximalQueueSize ...
func WithMaximalQueueSize(maximalQueueSize int) ClientOption {
	return func(clientConfig *ClientConfig) {
		clientConfig.maximalQueueSize = maximalQueueSize
	}
}

// WithQueues ...
func WithQueues(queues []string) ClientOption {
	return func(clientConfig *ClientConfig) {
		clientConfig.queues = queues
	}
}

// WithIDGenerator ...
func WithIDGenerator(idGenerator IDGenerator) ClientOption {
	return func(clientConfig *ClientConfig) {
		clientConfig.idGenerator = idGenerator
	}
}
