package hub

type (
	Subscription struct {
		Queues []string
		result chan Message
		done   chan struct{}
	}
	Message struct {
		Queue   string
		Message []byte
	}
)

func NewSubscription(queues []string) *Subscription {
	return &Subscription{
		Queues: queues,
		result: make(chan Message),
		done:   make(chan struct{}),
	}
}

func (s *Subscription) Need(queue string) bool {
	for _, q := range s.Queues {
		if q == queue {
			return true
		}
	}

	return false
}

func (s *Subscription) Send(msg Message) bool {
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

func (s *Subscription) Result() <-chan Message {
	return s.result
}

func (s *Subscription) Done() <-chan struct{} {
	return s.done
}

func (s *Subscription) Close() {
	close(s.done)
	close(s.result)
}
