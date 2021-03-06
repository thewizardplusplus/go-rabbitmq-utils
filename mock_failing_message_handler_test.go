// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package rabbitmqutils

import (
	amqp "github.com/streadway/amqp"
	mock "github.com/stretchr/testify/mock"
)

// MockFailingMessageHandler is an autogenerated mock type for the FailingMessageHandler type
type MockFailingMessageHandler struct {
	mock.Mock
}

// HandleMessage provides a mock function with given fields: message
func (_m *MockFailingMessageHandler) HandleMessage(message amqp.Delivery) error {
	ret := _m.Called(message)

	var r0 error
	if rf, ok := ret.Get(0).(func(amqp.Delivery) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
