package main

import (
	"reflect"
)

type (
	Request struct {
		Queues   []string
		Callback func(*Response)
	}
	Response struct {
		Queue   string
		Message Message
	}
)

var (
	pool = []*Request{}
)

func Register(q string, msg Message) {
	for i, r := range pool {
		for _, queueName := range r.Queues {
			if queueName == q {
				go r.Callback(&Response{Queue: queueName, Message: msg})
				pool = append(pool[:i], pool[i+1:]...)
				return
			}
		}
	}
	GetQueue(q).Push(msg)
}

func Process(r *Request) {
	for _, queueName := range r.Queues {
		q := GetQueue(queueName)
		if q.Counter.Distance() > 0 {
			if msg, err := q.Fetch(); err != nil {
				go r.Callback(nil)
			} else {
				go r.Callback(&Response{Queue: queueName, Message: msg})
			}
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
