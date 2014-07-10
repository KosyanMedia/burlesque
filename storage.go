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
		payload := <-saver
		i := payload.Queue.Counter.Write + 1
		key := NewKey(payload.Queue.Name, i)

		if err := storage.Set(key, payload.Message); err != nil {
			rollbar.Error("error", err)
		} else {
			payload.Queue.Counter.Incr()
		}
	}
}
