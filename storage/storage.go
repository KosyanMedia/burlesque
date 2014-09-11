package storage

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/ww/cabinet"
)

const (
	stateMetaKey      = "state"
	stateSaveInterval = 1 // seconds
)

type (
	Storage struct {
		kyoto    *cabinet.KCDB
		counters map[string]*counter
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

	return
}

func (s *Storage) Get(queue string) (message []byte, ok bool) {
	if _, exist := s.counters[queue]; !exist {
		return
	}
	if size := s.counters[queue].distance(); size == 0 {
		return
	}

	var index uint
	select {
	case index = <-s.counters[queue].stream:
	default:
		return
	}

	key := makeKey(queue, index)
	message, err := s.kyoto.Get(key)
	if err != nil {
		panic(err)
	}
	if err := s.kyoto.Remove(key); err != nil {
		panic(err)
	}
	ok = true

	return
}

func (s *Storage) Put(queue string, message []byte) (err error) {
	if _, ok := s.counters[queue]; !ok {
		s.counters[queue] = newCounter(0, 0)
	}

	s.counters[queue].tryWrite(func(index uint) bool {
		key := makeKey(queue, index)
		err = s.kyoto.Set(key, message)

		return (err == nil)
	})

	return
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
	state := make(map[string]map[string]uint)
	for queue, ctr := range s.counters {
		state[queue] = map[string]uint{
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
		state = make(map[string]map[string]uint)
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

func makeKey(queue string, index uint) []byte {
	return []byte(strings.Join([]string{queue, strconv.FormatUint(uint64(index), 10)}, "_"))
}
