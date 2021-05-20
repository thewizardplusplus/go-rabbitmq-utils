package rabbitmqutils

// ClientOptions ...
type ClientOptions struct {
	maximalQueueSize int
	queues           []string
}

// ClientOption ...
type ClientOption func(options *ClientOptions)

// WithMaximalQueueSize ...
func WithMaximalQueueSize(maximalQueueSize int) ClientOption {
	return func(options *ClientOptions) {
		options.maximalQueueSize = maximalQueueSize
	}
}
