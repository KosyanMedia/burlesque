package main

import (
	"github.com/stvp/rollbar"
)

type (
	Queue struct {
		Name    string
		Counter Counter
	}
)

var (
	queues = make(map[string]*Queue)
)

func (q *Queue) Push(msg Message) {
	saver <- Payload{Queue: q, Message: msg}
}

func (q *Queue) Fetch() (Message, error) {
	i := q.Counter.Next()
	key := NewKey(q.Name, i)

	msg, err := storage.Get(key)
	if err != nil {
		rollbar.Error("error", err)
		return msg, err
	}

	defer func() {
		if err := storage.Remove(key); err != nil {
			rollbar.Error("error", err)
		}
	}()

	return msg, nil
}

func GetQueue(name string) *Queue {
	if _, ok := queues[name]; !ok {
		RegisterQueue(name, 0, 0)
	}
	return queues[name]
}

func RegisterQueue(name string, wi, ri uint) {
	queues[name] = &Queue{Name: name, Counter: NewCounter(wi, ri)}
}
