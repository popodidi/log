package file

import (
	"fmt"
	"io"
	"sync"
)

func newRotate(rotator Rotator) *rotate {
	return &rotate{
		rotator: rotator,
	}
}

var _ io.WriteCloser = (*rotate)(nil)

type rotate struct {
	sync.Once // to seek to last
	sync.Mutex

	rotator Rotator

	size   int
	single *single
}

func (r *rotate) Write(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()
	r.Do(func() { err = r.seek() })
	if err != nil {
		return
	}
	err = r.rotateIfNeed()
	if err != nil {
		return
	}
	n, err = r.single.Write(p)
	if err != nil {
		return
	}

	r.size += n
	return
}

func (r *rotate) Close() error {
	r.Lock()
	defer r.Unlock()
	if r.single == nil {
		return nil
	}
	return r.single.Close()
}

func (r *rotate) seek() (err error) {
	for {
		if r.single == nil {
			err = r.rotate()
			if err != nil {
				return
			}
			continue
		}
		var exist bool
		exist, err = r.single.exist()
		if err != nil {
			return
		}
		if !exist {
			return
		}
		err = r.rotate()
		if err != nil {
			return
		}
	}
}

func (r *rotate) rotateIfNeed() error {
	if r.single != nil && !r.rotator.ShouldRotate(r.size) {
		return nil
	}
	return r.rotate()
}

func (r *rotate) rotate() error {
	if r.single != nil {
		err := r.single.Close()
		if err != nil {
			return err
		}
	}
	r.size = 0
	r.single = newSingle(r.rotator.Next())
	return nil
}

// Rotator defines the file rotation behavior interface.
type Rotator interface {
	ShouldRotate(size int) bool
	Next() string
}

var _ Rotator = (*decorateNameRotator)(nil)

type decorateNameRotator struct {
	rot      Rotator
	decorate func(string) string
}

func (r *decorateNameRotator) ShouldRotate(size int) bool {
	return r.rot.ShouldRotate(size)
}

func (r *decorateNameRotator) Next() string {
	return r.decorate(r.rot.Next())
}

var _ Rotator = (*baseRotator)(nil)

type baseRotator struct {
	next func() string
	size int

	prefix string
	index  int
}

func (r *baseRotator) ShouldRotate(size int) bool {
	return r.size >= 0 && size >= r.size
}

func (r *baseRotator) Next() string {
	n := r.next()
	if n != r.prefix {
		r.prefix = n
		r.index = 0
	} else {
		r.index++
	}
	return fmt.Sprintf("%s.%d", r.prefix, r.index)
}
