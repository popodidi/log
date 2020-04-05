package async

import (
	"context"

	"github.com/popodidi/log"
)

const defaultBufsize = 1 << 10

// New returns an asynchronous handler with default buffer size.
func New(handler log.Handler) log.CloseHandler {
	return newHandler(defaultBufsize, handler)
}

// NewWithBuf returns an asynchronous handler with specific buffer size.
func NewWithBuf(bufsize int, handler log.Handler) log.CloseHandler {
	return newHandler(bufsize, handler)
}

func newHandler(bufsize int, handler log.Handler) *asyncHandler {
	h := &asyncHandler{
		done:    make(chan struct{}),
		ch:      make(chan *log.Entry, bufsize),
		handler: handler,
	}
	h.ctx, h.cancel = context.WithCancel(context.Background())
	go h.run()
	return h
}

var _ log.CloseHandler = (*asyncHandler)(nil)

type asyncHandler struct {
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan struct{}
	ch      chan *log.Entry
	handler log.Handler
}

func (h *asyncHandler) Close() error {
	h.cancel()
	<-h.done
	if closer, ok := h.handler.(log.CloseHandler); ok {
		return closer.Close()
	}
	return nil
}

func (h *asyncHandler) Handle(entry *log.Entry) {
	if h.ctx.Err() != nil {
		return
	}
	select {
	case h.ch <- entry:
	case <-h.ctx.Done():
	}
}

func (h *asyncHandler) run() {
	for {
		select {
		case entry := <-h.ch:
			h.handler.Handle(entry)
		case <-h.ctx.Done():
			for len(h.ch) > 0 { // drain the channel
				entry := <-h.ch
				h.handler.Handle(entry)
			}
			h.done <- struct{}{}
			return
		}
	}
}
