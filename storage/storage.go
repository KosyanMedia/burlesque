package storage

import (
  "encoding/binary"
  "sync"
  "github.com/jmhodges/levigo"
  "errors"
)

type (
	Storage struct {
		db    *levigo.DB
	}

  Queue struct {
  	sync.RWMutex
  	head    uint64
  	tail    uint64
  }

)


var (
  queues = make(map[string]*Queue)
  db_ro = levigo.NewReadOptions()
  db_wo = levigo.NewWriteOptions()
)

func getQueue(queue_key string) *Queue {
  var ok bool
  if _, ok = queues[queue_key]; ok == false {
    queues[queue_key] = &Queue{
      head: 0,
      tail: 0,
    }
  }

  return queues[queue_key]
}

func idToKey(id uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, id)
	return key
}

func (s *Storage) GetItemById(q *Queue, id uint64) (item []byte, err error) {
	// Check if empty or out of bounds.
	if q.Length() == 0 {
		err = errors.New("Empty queue")
    return
	} else if id <= q.head || id > q.tail {
		err = errors.New("Queue is out of bounds")
    return
	}

	// Get item from database.
	if item, err = s.db.Get(db_ro, idToKey(id)); err != nil {
		return nil, err
	}

	return
}

func (q *Queue) Length() uint64 {
	return q.tail - q.head
}

func New(path string) (s *Storage, err error) {
  var (
    db *levigo.DB
  )

  opts := levigo.NewOptions()
  opts.SetCache(levigo.NewLRUCache(512 * 1024 * 1024))
  opts.SetCreateIfMissing(true)
  opts.SetBlockSize(256 * 1024)
  opts.SetBlockRestartInterval(8)
  opts.SetMaxOpenFiles(128)
  opts.SetInfoLog(nil)
  opts.SetWriteBufferSize(512 * 1024 * 1024)
  opts.SetParanoidChecks(false)
  opts.SetCompression(levigo.SnappyCompression)
  opts.SetFilterPolicy(levigo.NewBloomFilter(32))

  if db, err = levigo.Open(path, opts); err != nil {
    return
  }

  db_ro.SetFillCache(false)
  db_ro.SetVerifyChecksums(false)
  db_wo.SetSync(false)

	s = &Storage{
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

  q := getQueue(queue)
  q.Lock()
  defer q.Unlock()

  message, err := s.GetItemById(q, q.head + 1)
	if err != nil {
		return
	}

  if err := s.db.Delete(db_wo, idToKey(q.head + 1)); err != nil {
    return
  }

  q.head++

	ok = true
	return
}

func (s *Storage) Put(queue string, message []byte) (err error) {
  q := getQueue(queue)
  q.Lock()
  defer q.Unlock()

  if err = s.db.Put(db_wo, idToKey(q.tail + 1), message); err != nil {
    return
  }

  q.tail++
  return
}

func (s *Storage) Flush(queue string) (messages [][]byte) {
  // s.db.LClear([]byte(queue))
	return
}

func (s *Storage) QueueSizes() map[string]uint64 {
	info := make(map[string]uint64)
  for k, q := range queues {
    info[k] = q.Length()
  }

	return info
}

func (s *Storage) Close() (err error) {
	s.db.Close()
	return
}
