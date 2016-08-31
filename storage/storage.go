package storage

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"sync"

	"github.com/ezotrank/cabinetgo"
)

const (
	stateMetaKey      = "state"
	stateSaveInterval = 1 // seconds
)

type (
	Storage struct {
		kyoto    *cabinet.KCDB
		counters map[string]*counter
		sync.Mutex
	}
)

func New(path string) (s *Storage, err error) {
	kyoto := cabinet.New()
	if err = kyoto.Open(path, cabinet.KCOWRITER|cabinet.KCOCREATE); err != nil {
		return
	}

	s = &Storage{
		kyoto:    kyoto,
		counters: make(map[string]*counter),
	}
	s.loadState()
	go s.keepStatePersisted()

	return
}

func (s *Storage) Get(queue string, done <-chan struct{}) (message []byte, ok bool) {
	defer s.Unlock()
	s.Lock()
	if _, exist := s.counters[queue]; !exist {
		return
	}
	if size := s.counters[queue].distance(); size == 0 {
		return
	}

	var index int64
	select {
	case index = <-s.counters[queue].stream:
	case <-done:
		return
	}

	key := makeKey(queue, index)
	var err error
	if message, err = s.kyoto.Get(key); err != nil {
		return
	}
	if err := s.kyoto.Remove(key); err != nil {
		return
	}
	ok = true

	return
}

func (s *Storage) Put(queue string, message []byte) (err error) {
	defer s.Unlock()
	s.Lock()
	if _, ok := s.counters[queue]; !ok {
		s.counters[queue] = newCounter(0, 0)
	}

	s.counters[queue].tryWrite(func(index int64) bool {
		key := makeKey(queue, index)
		err = s.kyoto.Set(key, message)

		return (err == nil)
	})

	return
}

func (s *Storage) Flush(queue string) (messages [][]byte) {
	done := make(chan struct{})

	for {
		if msg, ok := s.Get(queue, done); ok {
			messages = append(messages, msg)
		} else {
			return
		}
	}

	return
}

func (s *Storage) QueueSizes() map[string]int64 {
	info := make(map[string]int64)

	for queue, c := range s.counters {
		info[queue] = c.distance()
	}

	return info
}

func (s *Storage) Info() map[string]interface{} {
	info := make(map[string]interface{})
	status, err := s.kyoto.Status()
	if err != nil {
		panic(err)
	}

	status = status[:len(status)-1] // Removing trailing new line
	tokens := strings.Split(status, "\n")
	for _, t := range tokens {
		tt := strings.Split(t, "\t")
		num, err := strconv.Atoi(tt[1])
		if err != nil {
			info[tt[0]] = tt[1]
		} else {
			info[tt[0]] = num
		}
	}

	return info
}

func (s *Storage) Close() (err error) {
	if err = s.kyoto.Sync(true); err != nil {
		return
	}
	err = s.kyoto.Close()

	return
}

// State

func (s *Storage) saveState() (err error) {
	state := make(map[string]map[string]int64)
	for queue, ctr := range s.counters {
		state[queue] = map[string]int64{
			"wi": ctr.write,
			"ri": ctr.read,
		}
	}

	jsn, _ := json.Marshal(state)
	err = s.kyoto.Set([]byte(stateMetaKey), jsn)

	return
}

func (s *Storage) loadState() (err error) {
	var (
		jsn   []byte
		state = make(map[string]map[string]int64)
	)

	if jsn, err = s.kyoto.Get([]byte(stateMetaKey)); err != nil {
		return
	}
	if err = json.Unmarshal(jsn, &state); err != nil {
		return
	}

	for queue, meta := range state {
		s.counters[queue] = newCounter(meta["wi"], meta["ri"])
	}

	return
}

func (s *Storage) keepStatePersisted() {
	t := time.NewTicker(stateSaveInterval * time.Second)

	for {
		select {
		case <-t.C:
			if err := s.saveState(); err != nil {
				panic("Failed to persist state")
			}
			if err := s.kyoto.Sync(false); err != nil {
				panic("Failed to sync storage")
			}
		}
	}
}

func makeKey(queue string, index int64) []byte {
	return []byte(strings.Join([]string{
		queue,
		strconv.FormatInt(index, 10),
	}, "_"))
}
