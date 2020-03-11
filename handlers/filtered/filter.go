package filtered

import "github.com/popodidi/log"

// Debug only handles debug logs.
func Debug(handler log.Handler) log.Handler {
	return Levels(handler, log.Debug)
}

// Info only handles info logs.
func Info(handler log.Handler) log.Handler {
	return Levels(handler, log.Info)
}

// Notice only handles notice logs.
func Notice(handler log.Handler) log.Handler {
	return Levels(handler, log.Notice)
}

// Warn only handles warn logs.
func Warn(handler log.Handler) log.Handler {
	return Levels(handler, log.Warn)
}

// Error only handles error logs.
func Error(handler log.Handler) log.Handler {
	return Levels(handler, log.Error)
}

// Critical only handles critical logs.
func Critical(handler log.Handler) log.Handler {
	return Levels(handler, log.Critical)
}

// Levels only handles logs with level.
func Levels(handler log.Handler, levels ...log.Level) log.Handler {
	if len(levels) == 0 {
		return handler
	}
	m := map[log.Level]struct{}{}
	for _, l := range levels {
		m[l] = struct{}{}
	}
	return Handler(handler, func(entry *log.Entry) bool {
		_, ok := m[entry.Level]
		return ok
	})
}

// Handler returns a filtered handler with shouldLog filter.
func Handler(handler log.Handler, shouldLog func(*log.Entry) bool) log.Handler {
	return &filterBase{
		shouldLog: shouldLog,
		handler:   handler,
	}
}

type filterBase struct {
	shouldLog func(*log.Entry) bool
	handler   log.Handler
}

func (b *filterBase) Close() error {
	if closer, ok := b.handler.(log.CloseHandler); ok {
		return closer.Close()
	}
	return nil
}

func (b *filterBase) Handle(entry *log.Entry) {
	if b.shouldLog(entry) {
		b.handler.Handle(entry)
	}
}
