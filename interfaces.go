package log

// Logger defines the logger interface.
type Logger interface {
	Clone() Logger

	GetID() string
	GetLabels() Labels

	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Critical(format string, args ...interface{})

	Log(*Entry)
}

// Labels defines a log label map.
type Labels interface {
	Get(string) (string, bool)
	Set(string, string)
	Delete(string)
	Clone() Labels
}

// Handler defines a log handler.
type Handler interface {
	Handle(*Entry)
}

// null handler
type null struct{}

func (n *null) Handle(entry *Entry) {}
