package mail

type StubSender struct {
}

func NewStubSender() *StubSender {
	return &StubSender{}
}

func (s *StubSender) Send(input Payload, attempts int) error {
	return nil
}
