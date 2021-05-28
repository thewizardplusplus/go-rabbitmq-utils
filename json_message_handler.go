package rabbitmqutils

import (
	"reflect"
)

// SpecificMessageHandler ...
type SpecificMessageHandler interface {
	MessageType() reflect.Type
	HandleMessage(message interface{}) error
}

// JSONMessageHandler ...
type JSONMessageHandler struct {
	MessageHandler SpecificMessageHandler
}
