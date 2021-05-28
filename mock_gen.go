package rabbitmqutils

import (
	"github.com/streadway/amqp"
)

//go:generate mockery --name=DialerInterface --inpackage --case=underscore --testonly

// DialerInterface ...
//
// It is used only for mock generating.
//
type DialerInterface interface {
	Dial(dsn string) (MessageBrokerConnection, error)
}

//go:generate mockery --name=ContextCancellerInterface --inpackage --case=underscore --testonly

// ContextCancellerInterface ...
//
// It is used only for mock generating.
//
type ContextCancellerInterface interface {
	CancelContext()
}

//go:generate mockery --name=AMQPAcknowledger --inpackage --case=underscore --testonly

// AMQPAcknowledger ...
//
// It is used only for mock generating.
//
type AMQPAcknowledger interface {
	amqp.Acknowledger
}
