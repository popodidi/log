package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// Entry defines a log entry.
type Entry struct {
	Level     Level
	Tag       string
	Labels    Labels
	Log       string
	Time      time.Time
	DebugInfo string
}

// SetDebugInfo returns the debug info of the caller function.
func (e *Entry) SetDebugInfo(stackNum int) {
	funcName := "???"
	pc, file, line, ok := runtime.Caller(stackNum)
	if !ok {
		file = "???"
		line = -1
	} else {
		funcName = runtime.FuncForPC(pc).Name()
		file = filepath.Base(file)
	}
	e.DebugInfo = fmt.Sprintf("%s:%d:%s", file, line, funcName)
}
