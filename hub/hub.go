package hub

import (
	"sync"
	"time"

	"github.com/KosyanMedia/burlesque/stats"
	"github.com/KosyanMedia/burlesque/storage"
)

type (
	Hub struct {
		storage     *storage.Storage
		subscribers []*Subscription
		lock        sync.Mutex
		statistics  *stats.Stats
	}
	Message struct {
		Queue   string
		Message []byte
	}
	MessageDump struct {
		Queue   string `json:"queue"`
		Message string `json:"message"`
	}
)

func New(st *storage.Storage) *Hub {
	h := &Hub{
		storage:     st,
		subscribers: []*Subscription{},
		statistics:  stats.New(),
	}

	go h.cleanupEverySecond()

	return h
}

func (h *Hub) Pub(queue string, msg []byte) bool {
	h.statistics.AddMessage(queue)
	for _, s := range h.subscribers {
		if ok := s.Need(queue); ok {
			// Check if subscription is already served
			select {
			case <-s.Done():
				continue
			default:
			}

			if ok := s.Send(Message{queue, msg}); ok {
				h.statistics.AddDelivery(queue)
				return true
			}
		}
	}

	err := h.storage.Put(queue, msg)

	return (err == nil)
}

func (h *Hub) Sub(s *Subscription) {
	for _, queue := range s.Queues {
		if msg, okGot := h.storage.Get(queue, s.Done()); okGot {
			if okSent := s.Send(Message{queue, msg}); okSent {
				h.statistics.AddDelivery(queue)
				return
			}
		}
	}

	h.lock.Lock()
	h.subscribers = append(h.subscribers, s)
	h.lock.Unlock()
}

func (h *Hub) Flush(queues []string) []MessageDump {
	messages := []MessageDump{}

	for _, queue := range queues {
		for _, msg := range h.storage.Flush(queue) {
			messages = append(messages, MessageDump{queue, string(msg)})
		}
	}

	return messages
}

func (h *Hub) Info() map[string]map[string]int64 {
	info := make(map[string]map[string]int64)

	for queue, size := range h.storage.QueueSizes() {
		info[queue] = map[string]int64{
			"messages":      size,
			"subscriptions": 0,
		}
	}
	for _, sub := range h.subscribers {
		for _, queue := range sub.Queues {
			if _, ok := info[queue]; !ok {
				info[queue] = map[string]int64{"messages": 0}
			}
			if _, ok := info[queue]["subscriptions"]; !ok {
				info[queue]["subscriptions"] = 0
			}
			info[queue]["subscriptions"]++
		}
	}

	return info
}

func (h *Hub) Rates(queue string) (in, out int64) {
	return h.statistics.Rates(queue)
}

func (h *Hub) RateHistory(queue string) (in, out []int64) {
	return h.statistics.RateHistory(queue)
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

	tmp := h.subscribers[:0]
	for _, s := range h.subscribers {
		select {
		case <-s.Done():
			continue
		default:
		}
		tmp = append(tmp, s)
		for _, queue := range s.Queues {
			if msg, okGot := h.storage.Get(queue, s.Done()); okGot {
				if okSent := s.Send(Message{queue, msg}); okSent {
					h.statistics.AddDelivery(queue)
					return
				}
			}
		}
	}
	h.subscribers = tmp
}
