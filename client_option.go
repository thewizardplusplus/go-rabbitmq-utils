package rabbitmqutils

// ClientConfig ...
type ClientConfig struct {
	maximalQueueSize int
	queues           []string
}

// ClientOption ...
type ClientOption func(clientConfig *ClientConfig)

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
