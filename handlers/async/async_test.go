package async

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/popodidi/log"
)

type dummy struct {
	sync.WaitGroup
	count int
}

func (d *dummy) Handle(entry *log.Entry) {
	d.Wait()
	d.count++
}

func TestAsync(t *testing.T) {
	var count = 77

	d := &dummy{}
	async := newHandler(count-1, d)

	d.Add(1) // hold dummy handler

	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			async.Handle(nil)
		}()
	}
	wg.Wait()

	require.Equal(t, 0, d.count)
	// The first is out and blocked by the dummy handler.
	require.Len(t, async.ch, count-1)

	select {
	case async.ch <- nil:
		require.FailNow(t, "channel must be full")
	default:
	}

	d.Done() // release dummy handler
	time.Sleep(500 * time.Millisecond)
	require.Equal(t, count, d.count)
	require.Len(t, async.ch, 0)

	d.count = 0 // reset for the next test case
	d.Add(1)    // hold dummy handler
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			async.Handle(nil)
		}()
	}
	wg.Wait()

	require.Equal(t, 0, d.count)
	require.Len(t, async.ch, count-1)

	done := make(chan struct{})
	go func() {
		require.NoError(t, async.Close())
		done <- struct{}{}
	}()
	select {
	case <-done:
		require.FailNow(t, "async should be blocked")
	default:
	}
	d.Done() // release dummy handler
	time.Sleep(500 * time.Millisecond)
	select {
	case <-done:
	default:
		require.FailNow(t, "async should be done")
	}
	require.Equal(t, count, d.count)
	require.Len(t, async.ch, 0)
}
