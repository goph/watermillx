package watermillx

import (
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
)

type InmemoryPublisher struct {
	messages map[string][]*message.Message

	mu sync.Mutex
}

// NewInmemoryPublisher returns a publisher that keeps the published messages in memory without forwarding them.
// It is useful for testing purposes.
func NewInmemoryPublisher() *InmemoryPublisher {
	return &InmemoryPublisher{
		messages: make(map[string][]*message.Message),
	}
}

// Publish stores the messages in memory.
func (p *InmemoryPublisher) Publish(topic string, messages ...*message.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.messages[topic] = append(p.messages[topic], messages...)

	return nil
}

// CLose does nothing.
func (*InmemoryPublisher) Close() error {
	return nil
}

// Messages returns the messages stored in memory.
func (p *InmemoryPublisher) Messages() map[string][]*message.Message {
	return p.messages
}

// Messages returns the messages stored in memory for a specific topic.
func (p *InmemoryPublisher) MessagesForTopic(topic string) []*message.Message {
	messages, ok := p.messages[topic]
	if !ok {
		return []*message.Message{}
	}

	return messages
}
