package multi

import (
	"github.com/popodidi/log"
)

// New return a multi handler.
//
// Deprecated: use log.MultiHandler instead.
func New(handlers ...log.Handler) log.CloseHandler {
	return log.MultiHandler(handlers...)
}
