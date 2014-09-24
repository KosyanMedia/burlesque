package hub

type (
	Subscription struct {
		Queues []string
		result chan Result
		done   chan struct{}
	}
	Result struct {
		Queue   string
		Message []byte
	}
)

func NewSubscription(queues []string) *Subscription {
	return &Subscription{
		Queues: queues,
		result: make(chan Result),
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

func (s *Subscription) Send(res Result) bool {
	success := make(chan bool)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				success <- false
			}
		}()

		s.result <- res
		success <- true
	}()

	return <-success
}

func (s *Subscription) Result() <-chan Result {
	return s.result
}

func (s *Subscription) Done() <-chan struct{} {
	return s.done
}

func (s *Subscription) Close() {
	close(s.done)
	close(s.result)
}
