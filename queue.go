package main

type (
	queue struct {
		name    string
		counter *counter
	}
)

var (
	queues = make(map[string]*queue)
)

func (q *queue) push(msg message) bool {
	var err error

	q.counter.write(func(i uint) bool {
		key := newKey(q.name, i)
		err = storage.Set(key, msg)
		if err != nil {
			alert(err, "Failed to write %d bytes to record '%s'", len(msg), key)
		}

		return (err == nil)
	})

	return (err == nil)
}

func (q *queue) tryFetch(abort chan bool) (message, bool) {
	if q.counter.distance() > 0 {
		return q.fetch(abort)
	} else {
		return message{}, false
	}
}

func (q *queue) fetch(abort chan bool) (message, bool) {
	var i uint

	select {
	case i = <-q.counter.read:
	case <-abort:
		return message{}, false
	}

	k := newKey(q.name, i)
	msg, err := storage.Get(k)
	if err != nil {
		alert(err, "Failed to read record '%s'", k)
		return msg, false
	}

	err = storage.Remove(k)
	if err != nil {
		alert(err, "Failed to delete record '%s'", k)
		return msg, false
	}

	return msg, true
}

func getQueue(name string) *queue {
	if _, ok := queues[name]; !ok {
		registerQueue(name, 0, 0)
	}
	return queues[name]
}

func registerQueue(name string, wi, ri uint) {
	queues[name] = &queue{
		name:    name,
		counter: newCounter(wi, ri),
	}
}
