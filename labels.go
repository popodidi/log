package log

import "sync"

type labels sync.Map

// NewLabels returns a Labels instance.
func NewLabels() Labels {
	l := labels(sync.Map{})
	return &l
}

func (l *labels) Get(key string) (string, bool) {
	value, ok := (*sync.Map)(l).Load(key)
	return value.(string), ok
}

func (l *labels) Set(key string, value string) {
	(*sync.Map)(l).Store(key, value)
}

func (l *labels) Delete(key string) {
	(*sync.Map)(l).Delete(key)
}

func (l *labels) Clone() Labels {
	cloned := NewLabels()
	m := (*sync.Map)(l)
	m.Range(func(k, v interface{}) bool {
		cloned.Set(k.(string), v.(string))
		return true
	})

	return cloned
}
