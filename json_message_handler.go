package rabbitmqutils

import (
	"reflect"
)

// SpecificMessageHandler ...
type SpecificMessageHandler interface {
	MessageType() reflect.Type
	HandleMessage(message interface{}) error
}
