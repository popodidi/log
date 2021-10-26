package log

import "io"

// Logger defines the logger interface.
type Logger interface {
	Clone() Logger

	WithTag(tags ...string) Logger
	WithLabel(key, value string) Logger
	WithHandler(handlers ...Handler) Logger

	GetID() string
	GetTag() string
	GetLabels() Labels

	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Critical(format string, args ...interface{})

	Handler

	// Deprecated: Use Handler interface instead
	Log(*Entry)
}

// Labels defines a log label map.
type Labels interface {
	Set(string, string)
	Delete(string)

	RLabels
}

// RLabels defines a read only log label map.
type RLabels interface {
	Get(string) (string, bool)
	Clone() Labels
	CloneAsMap() map[string]string

	// Deprecated: Use CloneAsMap instead
	Map() map[string]string
}

// CloseHandler defines a log handler which is also an io.Closer.
type CloseHandler interface {
	io.Closer
	Handler
}

// Handler defines a log handler.
type Handler interface {
	Handle(*Entry)
}
