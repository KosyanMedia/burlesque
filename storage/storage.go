package storage

import (
  "github.com/siddontang/ledisdb/ledis"
  "github.com/siddontang/ledisdb/config"
)

type (
	Storage struct {
    l     *ledis.Ledis
		db    *ledis.DB
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
  cfg.LevelDB.BlockSize = 262144 // 256 KB
  cfg.LevelDB.CacheSize =  536870912 // 512 MB
  cfg.LevelDB.WriteBufferSize = 536870912   // 256 MB
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
	return
}

func (s *Storage) Flush(queue string) (messages [][]byte) {
  s.db.LClear([]byte(queue))
	return
}

func (s *Storage) QueueSizes() map[string]int64 {
  var count int64
	info := make(map[string]int64)
  members, _ := s.db.Scan(ledis.LIST, nil, 100, true, "")
  for i := range members {
    count, _ = s.db.LLen(members[i])
    info[string(members[i])] = count
  }

	return info
}

func (s *Storage) Close() (err error) {
	s.l.Close()
	return
}
