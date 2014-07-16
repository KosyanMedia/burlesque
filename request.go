package main

import (
	"reflect"
)

type (
	Request struct {
		Queues   []string
		Callback func(*Response)
		Abort    chan bool
	}
	Response struct {
		Queue   string
		Message Message
	}
)

var (
	pool = []*Request{}
)

func Register(q string, msg Message) bool {
	for i, r := range pool {
		for _, queueName := range r.Queues {
			if queueName == q {
				go r.Callback(&Response{Queue: queueName, Message: msg})
				pool = append(pool[:i], pool[i+1:]...)
				return

				return true
			}
		}
	}

	ok := GetQueue(q).Push(msg)
	return ok
}

func Process(r *Request) {
	for _, queueName := range r.Queues {
		q := GetQueue(queueName)
		msg, ok := q.TryFetch(r.Abort)
		if ok {
			go r.Callback(&Response{Queue: queueName, Message: msg})
			return
		}
	}
	pool = append(pool, r)
}

func Purge(r *Request) {
	for i, req := range pool {
		if reflect.ValueOf(r).Pointer() == reflect.ValueOf(req).Pointer() {
			pool = append(pool[:i], pool[i+1:]...)
			return
		}
	}
}
