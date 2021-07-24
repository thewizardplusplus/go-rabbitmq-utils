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

// NewStartModeHolder ...
func NewStartModeHolder() *StartModeHolder {
	return &StartModeHolder{mode: NotStarted}
}

// GetStartMode ...
func (holder *StartModeHolder) GetStartMode() StartMode {
	holder.lock.RLock()
	defer holder.lock.RUnlock()

	return holder.mode
}

// SetStartModeOnce ...
func (holder *StartModeHolder) SetStartModeOnce(mode StartMode) {
	holder.lock.Lock()
	defer holder.lock.Unlock()

	if holder.mode == NotStarted {
		holder.mode = mode
	}
}
