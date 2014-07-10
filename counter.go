package main

import ()

const (
	MaxIndex = ^uint(0)
)

type (
	Counter struct {
		Write  uint
		Read   uint
		stream chan uint
		inLoop bool
	}
)

func NewCounter(wi, ri uint) *Counter {
	c := &Counter{Write: wi, Read: ri}
	c.stream = make(chan uint)
	go c.Loop()
	return c
}

func (c *Counter) Incr() {
	c.Write++
	if !c.inLoop {
		c.inLoop = true
		go c.Loop()
	}
}

func (c *Counter) Next() uint {
	return <-c.stream
}

func (c *Counter) Distance() uint {
	d := c.Write - c.Read
	if d < 0 {
		d += MaxIndex
	}
	return d
}

func (c *Counter) Loop() {
	for c.Write > c.Read {
		c.stream <- c.Read + 1
		c.Read++
	}
	c.inLoop = false
}
