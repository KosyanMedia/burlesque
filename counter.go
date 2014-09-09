package main

import (
	"sync"
)

const (
	maxIndex = ^uint(0)
)

type (
	// Counter is responsible for operating queue read and write indexes
	counter struct {
		writeIndex uint // Number of the record last written to the queue
		readIndex  uint // Number of the record last read from the queue
		// If WriteIndex is greater than ReadIndex then there are unread messages
		// If WriteIndex is less tham ReadIndex then MaxIndex was reached

		read      chan uint
		mutex     sync.Mutex
		streaming *sync.Cond
	}
)

func newCounter(wi, ri uint) *counter {
	m := &sync.Mutex{}
	m.Lock()

	c := &counter{
		writeIndex: wi,
		readIndex:  ri,
		read:       make(chan uint),
		streaming:  sync.NewCond(m),
	}

	go c.stream()
	return c
}

func (c *counter) write(proc func(i uint) bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ok := proc(c.writeIndex + 1)
	if ok {
		c.writeIndex++
		c.streaming.Signal()
	}
}

func (c *counter) distance() uint {
	d := c.writeIndex - c.readIndex
	if d < 0 {
		d += maxIndex
	}
	return d
}

func (c *counter) stream() {
	for {
		if c.distance() == 0 {
			c.streaming.Wait()
		}

		c.read <- c.readIndex + 1
		c.readIndex++
	}
}
