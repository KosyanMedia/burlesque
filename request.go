package main

import (
	"sync"
)

type (
	Request struct {
		Queues   []string
		Callback func(*Response)
		Abort    chan bool
		Dead     bool
	}
	Response struct {
		Queue   string
		Message Message
	}
)

var (
	pool struct {
		Requests []*Request
		mutex    sync.Mutex
	}
)

func RegisterPublication(q string, msg Message) bool {
	pool.mutex.Lock()
	for i, r := range pool.Requests {
		for _, queueName := range r.Queues {
			if queueName == q {
				go r.Callback(&Response{Queue: queueName, Message: msg})
				pool.Requests = append(pool.Requests[:i], pool.Requests[i+1:]...)
				defer pool.mutex.Unlock()

				return true
			}
		}
	}
	pool.mutex.Unlock()

	ok := GetQueue(q).Push(msg)
	return ok
}

func RegisterSubscription(r *Request) {
	for _, queueName := range r.Queues {
		q := GetQueue(queueName)
		msg, ok := q.TryFetch(r.Abort)
		if ok {
			go r.Callback(&Response{Queue: queueName, Message: msg})
			return
		}
	}

	pool.mutex.Lock()
	pool.Requests = append(pool.Requests, r)
	pool.mutex.Unlock()
}

func (r *Request) Purge() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	r.Dead = true
	deleted := 0
	for i, req := range pool.Requests {
		if req.Dead {
			pool.Requests = append(pool.Requests[:i-deleted], pool.Requests[i-deleted+1:]...)
			deleted++
		}
	}
}
