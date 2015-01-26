package storage

import (
	"sync"
)

const (
	maxIndex = ^int64(0) // Max unit value
)

type (
	// Counter is responsible for operating queue read and write indexes
	counter struct {
		write int64 // Number of the record last written to the queue
		read  int64 // Number of the record last read from the queue
		// If write index is greater than read index then there are unread messages
		// If write index is less tham read index then max index was reached

		mutex     sync.Mutex
		stream    chan int64
		streaming *sync.Cond
	}
)

func newCounter(wi, ri int64) *counter {
	m := &sync.Mutex{}
	m.Lock()

	c := &counter{
		write:     wi,
		read:      ri,
		stream:    make(chan int64),
		streaming: sync.NewCond(m),
	}

	go c.increment()

	return c
}

func (c *counter) tryWrite(fn func(i int64) bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ok := fn(c.write + 1); ok {
		if c.write++; c.write < 0 {
			c.write = 0
		}

		c.streaming.Signal()
	}
}

func (c *counter) distance() int64 {
	d := c.write - c.read
	if d < 0 {
		d += maxIndex
	}
	return d
}

func (c *counter) increment() {
	for {
		if c.distance() == 0 {
			c.streaming.Wait()
		}

		next := c.read + 1
		if next < 0 {
			next = 0
		}

		c.stream <- next
		c.read = next
	}
}
