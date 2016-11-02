package storage

import (
	"github.com/siddontang/ledisdb/ledis"
	"github.com/siddontang/ledisdb/config"
	"sort"
)

type (
	Storage struct {
		l		 *ledis.Ledis
		db		*ledis.DB
		keys map[string]bool
	}
)

func New(path string) (s *Storage, err error) {
	var (
		l *ledis.Ledis
		db *ledis.DB
	)

	cfg := config.NewConfigDefault()
	cfg.DBName = "leveldb"
	cfg.DataDir = path

	cfg.LevelDB.Compression = true
	cfg.LevelDB.BlockSize = 262144 // 256K
	cfg.LevelDB.CacheSize =	536870912 // 512MB
	cfg.LevelDB.WriteBufferSize = 268435456 // 256MB
	cfg.DBSyncCommit = 0

	if l, err = ledis.Open(cfg); err != nil {
		return
	}

	if db, err = l.Select(0); err != nil {
		return
	}

	s = &Storage{
		l: l,
		db: db,
		keys: make(map[string]bool),
	}

	// Preload available keys in storage
	members, _ := s.db.Scan(ledis.LIST, nil, 100, true, "")
	for i := range members {
		s.keys[string(members[i])] = true
	}
	return
}

func (s *Storage) Get(queue string, done <-chan struct{}) (message []byte, ok bool) {
	select {
 	case <-done:
 		return
	default:
 	}

	message, err := s.db.LPop([]byte(queue))
	if message == nil || err != nil {
		return
	}

	ok = true
	return
}

func (s *Storage) Put(queue string, message []byte) (err error) {
	_, err = s.db.RPush([]byte(queue), message)
	s.keys[queue] = true
	return
}

func (s *Storage) Flush(queue string) (messages [][]byte) {
	s.db.LClear([]byte(queue))
	return
}

func (s *Storage) GetSortedKeys() []string {
	keys := make([]string, 0)
	for k, _ := range s.keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (s *Storage) QueueSizes() map[string]int64 {
	var count int64
	info := make(map[string]int64)
	for _, key := range s.GetSortedKeys() {
		count, _ = s.db.LLen([]byte(key))
		info[key] = count
	}

	return info
}

func (s *Storage) Close() (err error) {
	s.l.Close()
	return
}
