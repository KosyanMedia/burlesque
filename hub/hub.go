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

	h.subscribers = append(h.subscribers, s)
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
		inRate, outRate := h.statistics.Rates(queue)
		info[queue] = map[string]int64{
			"messages":      size,
			"subscriptions": 0,
			"in_rate":       inRate,
			"out_rate":      outRate,
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

func (h *Hub) RateHistory() map[string]map[string][]int64 {
	hist := map[string]map[string][]int64{}
	for queue, _ := range h.storage.QueueSizes() {
		hist[queue] = h.statistics.RateHistory(queue)
	}

	return hist
}

func (h *Hub) StorageInfo() map[string]interface{} {
	return h.storage.Info()
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
