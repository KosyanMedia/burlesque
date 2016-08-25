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
	cfg.DBName = "rocksdb"
	cfg.DataDir = path

  cfg.LevelDB.Compression = true
  cfg.LevelDB.BlockSize = 536870912 // 512 MB
  cfg.LevelDB.CacheSize =  536870912 // 512 MB
  cfg.LevelDB.WriteBufferSize = 536870912 // 512 MB

  //  https://github.com/facebook/rocksdb/wiki/RocksDB-Tuning-Guide
  cfg.RocksDB.Compression = 1
  cfg.RocksDB.BlockSize = 2 * 1024 * 1024 // 512 KB
  cfg.RocksDB.CacheSize =  1024 * 1024 * 1024 // 1 GB
  cfg.RocksDB.WriteBufferSize = 256 * 1024 * 1024 // 256 MB
  cfg.RocksDB.MaxWriteBufferNum = 8
  cfg.RocksDB.MinWriteBufferNumberToMerge = 2
  cfg.RocksDB.EnableStatistics = false
  cfg.RocksDB.BackgroundThreads = 1
  cfg.RocksDB.HighPriorityBackgroundThreads = 4
  cfg.RocksDB.MaxBackgroundFlushes = 1
  cfg.RocksDB.AllowOsBuffer = true
  cfg.RocksDB.DisableWAL = true // !
  cfg.RocksDB.DisableAutoCompactions = false
  cfg.RocksDB.UseFsync = false
  cfg.RocksDB.MaxBackgroundCompactions = 4
  cfg.RocksDB.DisableDataSync = true // !
  cfg.RocksDB.TargetFileSizeBase = 64 * 1024 * 1024 // 512 MB
  cfg.RocksDB.TargetFileSizeMultiplier = 1
  cfg.RocksDB.NumLevels = 4
  cfg.RocksDB.Level0FileNumCompactionTrigger = 2
  cfg.RocksDB.Level0SlowdownWritesTrigger = 12
  cfg.RocksDB.Level0StopWritesTrigger = 16
  cfg.RocksDB.MaxBytesForLevelBase = 640 * 1024 * 1024 // 512 MB
  cfg.RocksDB.MaxBytesForLevelMultiplier = 10

  // cfg.RocksDB.Compression = 0
  // cfg.RocksDB.BlockSize = 536870912 // 512 MB
  // cfg.RocksDB.CacheSize =  536870912 // 512 MB
  // cfg.RocksDB.WriteBufferSize = 536870912 // 512 MB
  // cfg.RocksDB.MaxWriteBufferNum = 4
  // cfg.RocksDB.MinWriteBufferNumberToMerge = 2
  // cfg.RocksDB.EnableStatistics = false
  // cfg.RocksDB.BackgroundThreads = 1
  // cfg.RocksDB.HighPriorityBackgroundThreads = 4
  // cfg.RocksDB.MaxBackgroundCompactions = runtime.GOMAXPROCS(-1)
  // cfg.RocksDB.MaxBackgroundFlushes = 1
  // cfg.RocksDB.AllowOsBuffer = false
  // cfg.RocksDB.DisableWAL = false
  // cfg.RocksDB.DisableAutoCompactions = true
  // cfg.RocksDB.UseFsync = true
  // cfg.RocksDB.MaxBackgroundCompactions = 3
  // cfg.RocksDB.DisableDataSync = false
  // cfg.RocksDB.TargetFileSizeBase = 536870912 // 512 MB
  // cfg.RocksDB.NumLevels = 4
  // cfg.RocksDB.Level0FileNumCompactionTrigger = 2
  // cfg.RocksDB.Level0SlowdownWritesTrigger = 8
  // cfg.RocksDB.Level0StopWritesTrigger = 16
  // cfg.RocksDB.MaxBytesForLevelBase = 536870912 // 512 MB
  // cfg.RocksDB.MaxBytesForLevelMultiplier = 10
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
