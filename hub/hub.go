package hub

import (
	"sync"
	"time"

	"github.com/KosyanMedia/burlesque/storage"
)

type (
	Hub struct {
		storage     *storage.Storage
		subscribers []*Subscription
		lock        sync.Mutex
	}
)

func New(st *storage.Storage) *Hub {
	h := &Hub{
		storage:     st,
		subscribers: []*Subscription{},
	}

	go h.cleanupEverySecond()

	return h
}

func (h *Hub) Pub(queue string, msg []byte) bool {
	for _, s := range h.subscribers {
		if ok := s.Need(queue); ok {
			select {
			case <-s.Done():
				continue
			default:
			}

			if ok := s.Send(Result{queue, msg}); ok {
				return true
			}
		}
	}

	err := h.storage.Put(queue, msg)

	return (err == nil)
}

func (h *Hub) Sub(s *Subscription) {
	for _, q := range s.Queues {
		if msg, ok := h.storage.Get(q); ok {
			s.Send(Result{q, msg})
			return
		}
	}

	h.subscribers = append(h.subscribers, s)
}

func (h *Hub) cleanupEverySecond() {
	t := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-t.C:
			h.cleanup()
		}
	}
}

func (h *Hub) cleanup() {
	h.lock.Lock()
	defer h.lock.Unlock()

	deleted := 0
	for i, s := range h.subscribers {
		select {
		case <-s.Done():
			h.subscribers = append(h.subscribers[:i-deleted], h.subscribers[i-deleted+1:]...)
			deleted++
		default:
		}
	}
}
