package watermillx

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

type publisherStub struct {
	messages map[string][]*message.Message
}

func newPublisherStub() *publisherStub {
	return &publisherStub{
		messages: make(map[string][]*message.Message),
	}
}

func (p *publisherStub) Publish(topic string, messages ...*message.Message) error {
	p.messages[topic] = append(p.messages[topic], messages...)

	return nil
}

func (p *publisherStub) Close() error {
	return nil
}
