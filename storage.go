package main

import (
	"github.com/ezotrank/cabinetgo"
	"github.com/stvp/rollbar"
)

type (
	Payload struct {
		Queue   *Queue
		Message Message
	}
)

var (
	storage = cabinet.New()
	saver   = make(chan Payload, 1000)
)

func SetupStorage() {
	err := storage.Open(cfg.Storage, cabinet.KCOWRITER|cabinet.KCOCREATE)
	if err != nil {
		panic(err)
	}
}

func PersistMessages() {
	for {
		p := <-saver

		p.Queue.Counter.Write(func(i uint) bool {
			key := NewKey(p.Queue.Name, i)
			err := storage.Set(key, p.Message)
			if err != nil {
				Error(err, "Failed to write %d bytes to record '%s'", len(p.Message), key)
			}

			return (err != nil)
		})
	}
}
