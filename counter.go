package main

import (
	"sync"
)

const (
	MaxIndex = ^uint(0)
)

type (
	Counter struct {
		WriteIndex uint
		ReadIndex  uint
		mutex      sync.Mutex
		stream     chan uint
		streaming  bool
	}
)

func NewCounter(wi, ri uint) Counter {
	c := Counter{
		WriteIndex: wi,
		ReadIndex:  ri,
		stream:     make(chan uint),
		streaming:  false,
	}
	if c.Distance() > 0 {
		go c.Stream()
	}
	return c
}

func (c *Counter) Write(proc func(i uint) bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ok := proc(c.WriteIndex + 1)
	if ok {
		c.WriteIndex++
		if !c.streaming {
			go c.Stream()
		}
	}
}

func (c *Counter) Read() uint {
	return <-c.stream
}

func (c *Counter) Distance() uint {
	d := c.WriteIndex - c.ReadIndex
	if d < 0 {
		d += MaxIndex
	}
	return d
}

func (c *Counter) Stream() {
	c.streaming = true
	for c.WriteIndex > c.ReadIndex {
		c.stream <- c.ReadIndex + 1
		c.ReadIndex++
	}
	c.streaming = false
}
