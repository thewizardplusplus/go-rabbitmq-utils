package rabbitmqutils

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
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

// HandleMessage ...
func (handler JSONMessageHandler) HandleMessage(message amqp.Delivery) error {
	data := reflect.New(handler.MessageHandler.MessageType())
	if err := json.Unmarshal(message.Body, data.Interface()); err != nil {
		return errors.Wrap(err, "unable to unmarshal the data")
	}

	err := handler.MessageHandler.HandleMessage(data.Elem().Interface())
	if err != nil {
		return errors.Wrap(err, "unable to handle the data")
	}

	return nil
}
