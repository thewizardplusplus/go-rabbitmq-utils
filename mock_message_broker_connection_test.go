// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package rabbitmqutils

import (
	amqp "github.com/streadway/amqp"
	mock "github.com/stretchr/testify/mock"
)

// MockMessageBrokerConnection is an autogenerated mock type for the MessageBrokerConnection type
type MockMessageBrokerConnection struct {
	mock.Mock
}

// Channel provides a mock function with given fields:
func (_m *MockMessageBrokerConnection) Channel() (*amqp.Channel, error) {
	ret := _m.Called()

	var r0 *amqp.Channel
	if rf, ok := ret.Get(0).(func() *amqp.Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*amqp.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *MockMessageBrokerConnection) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
