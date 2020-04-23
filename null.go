package log

// Null defines the null logger.
var Null Logger = NewLogger(Config{
	Tag:       "null",
	Threshold: last,
	Handler:   (*null)(nil),
})

// null handler
type null struct{}

func (n *null) Handle(entry *Entry) {}
