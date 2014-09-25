package storage

import (
	"sync"
)

const (
	maxIndex = ^uint(0) // Max unit value
)

type (
	// Counter is responsible for operating queue read and write indexes
	counter struct {
		write uint // Number of the record last written to the queue
		read  uint // Number of the record last read from the queue
		// If write index is greater than read index then there are unread messages
		// If write index is less tham read index then max index was reached

		mutex     sync.Mutex
		stream    chan uint
		streaming *sync.Cond
	}
)

func newCounter(wi, ri uint) *counter {
	m := &sync.Mutex{}
	m.Lock()

	c := &counter{
		write:     wi,
		read:      ri,
		stream:    make(chan uint),
		streaming: sync.NewCond(m),
	}

	go c.increment()

	return c
}

func (c *counter) tryWrite(fn func(i uint) bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ok := fn(c.write + 1); ok {
		c.write++
		c.streaming.Signal()
	}
}

func (c *counter) distance() uint {
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

		c.stream <- c.read + 1
		c.read++
	}
}
