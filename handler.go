package log

import (
	"fmt"
)

var _ CloseHandler = (*multiHandler)(nil)

type multiHandler []Handler

// MultiHandler return a multi handler.
func MultiHandler(handlers ...Handler) CloseHandler {
	h := multiHandler(handlers)
	h.expand()
	return &h
}

func (h *multiHandler) expand() {
	expanded := multiHandler{}
	for _, sub := range *h {
		if mh, ok := sub.(*multiHandler); ok {
			mh.expand()
			expanded = append(expanded, *mh...)
			continue
		}
		expanded = append(expanded, sub)
	}
	*h = expanded
}

func (h *multiHandler) Handle(entry *Entry) {
	l := len(*h)
	if l == 0 {
		return
	} else if l == 1 {
		(*h)[0].Handle(entry)
		return
	}
	for i := range *h {
		(*h)[i].Handle(entry)
	}
}

func (h *multiHandler) Close() error {
	l := len(*h)
	if l == 0 {
		return nil
	}
	// best effort to close every handler.
	var errs []error
	for i := range *h {
		if closer, ok := (*h)[i].(CloseHandler); ok {
			err := closer.Close()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}
