package rabbitmqutils

// StartMode ...
type StartMode int

// ...
const (
	NotStarted StartMode = iota
	Started
	StartedConcurrently
)
