package file

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

func newRotate(rotator Rotator) *rotate {
	return &rotate{
		rotator: rotator,
	}
}

type rotate struct {
	sync.Mutex

	rotator Rotator

	size   int
	single *single
}

func (r *rotate) Write(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()

	r.rotateIfNeed()
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
	return r.single.Close()
}

func (r *rotate) rotateIfNeed() {
	if r.single != nil && !r.rotator.ShouldRotate(r.size) {
		return
	}
	var err error
	if r.single != nil {
		err = r.single.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	r.size = 0
	r.single, err = newSingle(r.rotator.Next())
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Rotator defines the file rotation behavior interface.
type Rotator interface {
	ShouldRotate(size int) bool
	Next() string
}

// SecondRotator returns a rotator that rotates every second with file name as
// {UNIX_SECOND}.{IDX} and file size limited to `size` bytes.
func SecondRotator(size int) Rotator {
	return &baseRotator{
		next: func() string {
			return strconv.Itoa(int(time.Now().Unix()))
		},
		size: size,
	}
}

const dateFmt = "20060102"

// DateRotator returns a rotator that rotates every day with file name as
// {YYYYMMDD}.{IDX} and file size limited to `size` bytes.
func DateRotator(size int) Rotator {
	return &baseRotator{
		next: func() string {
			return time.Now().Format(dateFmt)
		},
		size: size,
	}
}

type baseRotator struct {
	next func() string
	size int

	prefix string
	index  int
}

func (r *baseRotator) ShouldRotate(size int) bool {
	return size >= r.size
}

func (r *baseRotator) Next() string {
	n := r.next()
	if n != r.prefix {
		r.prefix = n
		r.index = 0
	} else {
		r.index += 1
	}
	return fmt.Sprintf("%s.%d", r.prefix, r.index)
}
