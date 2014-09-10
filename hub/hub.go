package hub

import (
	"github.com/KosyanMedia/burlesque/storage"
)

type (
	Hub struct {
		storage     *storage.Storage
		subscribers []*Subscription
	}
)

func New(st *storage.Storage) *Hub {
	return &Hub{
		storage:     st,
		subscribers: []*Subscription{},
	}
}

func (h *Hub) Pub(queue string, msg []byte) bool {
	for _, s := range h.subscribers {
		if s.Queue == queue {
			select {
			case <-s.Done():
				continue
			default:
			}

			if ok := s.Send(msg); ok {
				return true
			}
		}
	}

	err := h.storage.Put(queue, msg)

	return (err == nil)
}

func (h *Hub) Sub(s *Subscription) {
	if msg, ok := h.storage.Get(s.Queue); ok {
		s.Send(msg)
	} else {
		h.subscribers = append(h.subscribers, s)
	}
}
