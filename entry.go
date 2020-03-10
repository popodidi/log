package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// Entry defines a log entry.
type Entry struct {
	StackNum int
	Level    Level
	Tag      string
	Labels   Labels
	Log      string
	Time     time.Time
}

// DebugInfo returns the debug info of the caller function.
func (e *Entry) DebugInfo() string {
	funcName := "???"
	pc, file, line, ok := runtime.Caller(e.StackNum)
	if !ok {
		file = "???"
		line = -1
	} else {
		funcName = runtime.FuncForPC(pc).Name()
		file = filepath.Base(file)
	}
	return fmt.Sprintf("%s:%d:%s", file, line, funcName)
}
