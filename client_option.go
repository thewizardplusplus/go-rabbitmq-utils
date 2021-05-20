package rabbitmqutils

// ClientOptions ...
type ClientOptions struct {
	maximalQueueSize int
	queues           []string
}

// ClientOption ...
type ClientOption func(options *ClientOptions)
