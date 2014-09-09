package main

import (
	"sync"
)

type (
	request struct {
		queues     []string
		responseCh chan response
		abort      chan bool
		dead       bool
	}
	response struct {
		queue   string
		message message
	}
)

var (
	pool struct {
		requests []*request
		mutex    sync.Mutex
	}
)

func registerPublication(q string, msg message) bool {
	for _, r := range pool.requests {
		if r.dead {
			continue
		}
		for _, qname := range r.queues {
			if qname == q {
				rsp := response{queue: q, message: msg}
				ok := r.tryRespond(rsp)
				if ok {
					return true
				}
			}
		}
	}

	ok := getQueue(q).push(msg)
	return ok
}

func registerSubscription(r *request) {
	for _, qname := range r.queues {
		q := getQueue(qname)
		msg, ok := q.tryFetch(r.abort)
		if ok {
			rsp := response{queue: qname, message: msg}
			ok := r.tryRespond(rsp)
			if !ok {
				q.push(msg)
			}

			return
		}
	}

	pool.requests = append(pool.requests, r)
}

func (r *request) tryRespond(rsp response) bool {
	okch := make(chan bool)

	go func() {
		defer func() {
			err := recover()
			if err != nil { // Panic!
				r.dead = true
				okch <- false
			}
		}()

		r.responseCh <- rsp // If channel is already closed expect a panic
		okch <- true
	}()

	ok := <-okch
	return ok
}

func (r *request) purge() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	r.dead = true
	deleted := 0
	for i, req := range pool.requests {
		if req.dead {
			pool.requests = append(pool.requests[:i-deleted], pool.requests[i-deleted+1:]...)
			deleted++
		}
	}
}
