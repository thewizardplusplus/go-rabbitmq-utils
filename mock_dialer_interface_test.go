// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package rabbitmqutils

import mock "github.com/stretchr/testify/mock"

// MockDialerInterface is an autogenerated mock type for the DialerInterface type
type MockDialerInterface struct {
	mock.Mock
}

// Dial provides a mock function with given fields: dsn
func (_m *MockDialerInterface) Dial(dsn string) (MessageBrokerConnection, error) {
	ret := _m.Called(dsn)

	var r0 MessageBrokerConnection
	if rf, ok := ret.Get(0).(func(string) MessageBrokerConnection); ok {
		r0 = rf(dsn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(MessageBrokerConnection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dsn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}