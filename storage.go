package main

import (
	"github.com/ezotrank/cabinetgo"
	"strconv"
	"strings"
)

type (
	Message []byte
	Key     []byte
	Payload struct {
		Queue   *Queue
		Message Message
	}
)

var (
	storage  = cabinet.New()
	payloads = make(chan Payload, 1000)
)

func NewKey(queue string, index uint) Key {
	istr := strconv.FormatUint(uint64(index), 10)
	key := strings.Join([]string{queue, istr}, "_")
	return Key(key)
}

func SetupStorage() {
	err := storage.Open(cfg.Storage, cabinet.KCOWRITER|cabinet.KCOCREATE)
	if err != nil {
		Error(err, "Failed to open database '%s'", cfg.Storage)
	}
}

func Persist(p Payload) {
	payloads <- p
}

func PersistMessages() {
	for {
		p := <-payloads

		p.Queue.Counter.Write(func(i uint) bool {
			key := NewKey(p.Queue.Name, i)
			err := storage.Set(key, p.Message)
			if err != nil {
				Error(err, "Failed to write %d bytes to record '%s'", len(p.Message), key)
			}

			return (err == nil)
		})
	}
}
