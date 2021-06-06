package rabbitmqutils

import (
	"time"

	"github.com/go-log/log"
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

//go:generate mockery --name=IDGeneratorInterface --inpackage --case=underscore --testonly

// IDGeneratorInterface ...
//
// It is used only for mock generating.
//
type IDGeneratorInterface interface {
	GenerateID() (string, error)
}

//go:generate mockery --name=ClockInterface --inpackage --case=underscore --testonly

// ClockInterface ...
//
// It is used only for mock generating.
//
type ClockInterface interface {
	Time() time.Time
}

//go:generate mockery --name=ContextCancellerInterface --inpackage --case=underscore --testonly

// ContextCancellerInterface ...
//
// It is used only for mock generating.
//
type ContextCancellerInterface interface {
	CancelContext()
}

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It is used only for mock generating.
//
type Logger interface {
	log.Logger
}

//go:generate mockery --name=AMQPAcknowledger --inpackage --case=underscore --testonly

// AMQPAcknowledger ...
//
// It is used only for mock generating.
//
type AMQPAcknowledger interface {
	amqp.Acknowledger
}
