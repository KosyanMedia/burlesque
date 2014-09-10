package hub

type (
	Subscription struct {
		Queue  string
		result chan<- []byte
		done   chan struct{}
	}
)

func NewSubscription(queue string, result chan<- []byte) *Subscription {
	return &Subscription{
		Queue:  queue,
		result: result,
		done:   make(chan struct{}),
	}
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
