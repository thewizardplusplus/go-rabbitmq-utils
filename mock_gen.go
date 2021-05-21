package rabbitmqutils

//go:generate mockery --name=ContextCancellerInterface --inpackage --case=underscore --testonly

// ContextCancellerInterface ...
//
// It is used only for mock generating.
//
type ContextCancellerInterface interface {
	CancelContext()
}
