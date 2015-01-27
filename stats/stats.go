package stats

import (
	"time"
)

const (
	maxHistoryLength = 300 // Five minutes
)

type (
	Stats struct {
		queues map[string]*meta
	}
	meta struct {
		cur    point
		points []point
	}
	point struct {
		in  int64
		out int64
	}
)

func New() *Stats {
	s := &Stats{
		q: map[string]*meta{},
	}

	go s.loopCollectSeconds()

	return s
}

func (s *Stats) AddMessage(queue string) {
	s.metaFor(queue).cur.in++
}

func (s *Stats) AddDelivery(queue string) {
	s.metaFor(queue).cur.out++
}

func (s *Stats) Rates(queue string) (in, out int64) {
	m := s.metaFor(queue)
	p := point{}
	if len(m.points) > 0 {
		p = m.points[len(m.points)-1]
	}

	return p.in, p.out
}

func (s *Stats) RateHistory(queue string) map[string][]int64 {
	hist := map[string][]int64{
		"in":  []int64{},
		"out": []int64{},
	}

	for _, p := range s.metaFor(queue).points {
		hist["in"] = append(hist["in"], p.in)
		hist["out"] = append(hist["out"], p.out)
	}

	return hist
}

func (s *Stats) loopCollectSeconds() {
	t := time.NewTicker(1 * time.Second)
	for {
		<-t.C
		s.collectSeconds()
	}
}

func (s *Stats) collectSeconds() {
	for _, m := range s.queues {
		m.points = append(m.points, m.cur)
		m.cur.in = 0
		m.cur.out = 0
		if len(m.points) > maxHistoryLength {
			m.points = m.points[1:]
		}
	}
}

func (s *Stats) metaFor(queue string) *meta {
	m, ok := s.queues[queue]
	if !ok {
		m = &meta{
			points: []point{},
		}
		s.queues[queue] = m
	}

	return m
}
