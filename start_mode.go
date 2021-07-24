package rabbitmqutils

import (
	"sync"
)

// StartMode ...
type StartMode int

// ...
const (
	NotStarted StartMode = iota
	Started
	StartedConcurrently
)

// StartModeHolder ...
type StartModeHolder struct {
	lock sync.RWMutex
	mode StartMode
}
