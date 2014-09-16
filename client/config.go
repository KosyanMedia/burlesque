package client

import (
	"time"
)

type (
	Config struct {
		Host    string
		Port    int
		Timeout time.Duration
	}
)

func (c *Config) UseDefaults() {
	c.Host = "127.0.0.1"
	c.Port = 4401
	c.Timeout = 60 * time.Second
}
