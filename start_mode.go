package rabbitmqutils

import (
	"sync"
)

type startMode int

const (
	notStarted startMode = iota
	started
	startedConcurrently
)

// StartModeHolder ...
type StartModeHolder struct {
	lock sync.RWMutex
	mode startMode
}

// NewStartModeHolder ...
func NewStartModeHolder() *StartModeHolder {
	return &StartModeHolder{mode: notStarted}
}

// GetStartMode ...
func (holder *StartModeHolder) GetStartMode() startMode {
	holder.lock.RLock()
	defer holder.lock.RUnlock()

	return holder.mode
}

// SetStartModeOnce ...
func (holder *StartModeHolder) SetStartModeOnce(mode startMode) {
	holder.lock.Lock()
	defer holder.lock.Unlock()

	if holder.mode == notStarted {
		holder.mode = mode
	}
}
