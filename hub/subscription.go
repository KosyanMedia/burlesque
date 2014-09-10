package hub

type (
	Subscription struct {
		queues []string
		result chan<- []byte
		done   chan struct{}
	}
)

func NewSubscription(queues []string, result chan<- []byte) *Subscription {
	return &Subscription{
		queues: queues,
		result: result,
		done:   make(chan struct{}),
	}
}

func (s *Subscription) Need(queue string) bool {
	for _, q := range s.queues {
		if q == queue {
			return true
		}
	}

	return false
}

func (s *Subscription) Send(msg []byte) bool {
	success := make(chan bool)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				success <- false
			}
		}()

		s.result <- msg
		success <- true
	}()

	return <-success
}

func (s *Subscription) Done() <-chan struct{} {
	return s.done
}

func (s *Subscription) Close() {
	close(s.result)
	close(s.done)
}
