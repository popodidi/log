package multi

import "github.com/popodidi/log"

type handler []log.Handler

// New return a multi handler.
func New(handlers ...log.Handler) log.Handler {
	h := handler(handlers)
	h.expand()
	return &h
}
func (h *handler) expand() {
	expanded := handler{}
	for _, sub := range *h {
		if mh, ok := sub.(*handler); ok {
			mh.expand()
			expanded = append(expanded, *mh...)
			continue
		}
		expanded = append(expanded, sub)
	}
	*h = expanded
}

func (h *handler) Handle(entry *log.Entry) {
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
