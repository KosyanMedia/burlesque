package main

import (
	"sync"
)

type (
	Request struct {
		Queues     []string
		ResponseCh chan Response
		Abort      chan bool
		Dead       bool
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
	for _, r := range pool.Requests {
		if r.Dead {
			continue
		}
		for _, qname := range r.Queues {
			if qname == q {
				rsp := Response{Queue: q, Message: msg}
				ok := r.TryRespond(rsp)
				if ok {
					return true
				}
			}
		}
	}

	ok := GetQueue(q).Push(msg)
	return ok
}

func RegisterSubscription(r *Request) {
	for _, qname := range r.Queues {
		q := GetQueue(qname)
		msg, ok := q.TryFetch(r.Abort)
		if ok {
			rsp := Response{Queue: qname, Message: msg}
			ok := r.TryRespond(rsp)
			if !ok {
				q.Push(msg)
			}

			return
		}
	}

	pool.Requests = append(pool.Requests, r)
}

func (r *Request) TryRespond(rsp Response) bool {
	okch := make(chan bool)

	go func() {
		defer func() {
			err := recover()
			if err != nil { // Panic!
				r.Dead = true
				okch <- false
			}
		}()

		r.ResponseCh <- rsp // If channel is already closed expect a panic
		okch <- true
	}()

	ok := <-okch
	return ok
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
