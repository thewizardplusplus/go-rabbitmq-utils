package rabbitmqutils

// Dialer ...
type Dialer func(dsn string) (MessageBrokerConnection, error)

// IDGenerator ...
type IDGenerator func() (string, error)

// ClientConfig ...
type ClientConfig struct {
	dialer           Dialer
	maximalQueueSize int
	queues           []string
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
