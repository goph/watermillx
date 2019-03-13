package watermillx

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

// MessageBus sends a message to its recipient(s) through a message publisher.
type MessageBus struct {
	publisher   message.Publisher
	marshaler   MessageMarshaler
	idgenerator IDGenerator
}

// NewMessageBus returns a new MessageBus instance.
func NewMessageBus(publisher message.Publisher, marshaler MessageMarshaler, idgenerator IDGenerator) *MessageBus {
	return &MessageBus{
		publisher:   publisher,
		marshaler:   marshaler,
		idgenerator: idgenerator,
	}
}

// MessageMarshaler marshals a message into a format that can be sent using the messaging protocol of Watermill.
type MessageMarshaler interface {
	// Marshal marshals a message into a format that can be sent using the messaging protocol of Watermill.
	Marshal(msg interface{}) ([]byte, error)
}

// IDGenerator generates a new ID.
type IDGenerator interface {
	// Generate generates a new ID.
	Generate() (string, error)
}

// Publish publishes a message.
func (b *MessageBus) Publish(ctx context.Context, topic string, m interface{}) error {
	payload, err := b.marshaler.Marshal(m)
	if err != nil {
		return errors.WithMessage(err, "cannot marshal message payload")
	}

	msgID, err := b.idgenerator.Generate()
	if err != nil {
		return errors.WithMessage(err, "cannot generate message ID")
	}

	msg := message.NewMessage(msgID, payload)
	msg.SetContext(ctx)

	err = b.publisher.Publish(topic, msg)
	if err != nil {
		return errors.WithMessage(err, "failed to publish message")
	}

	return nil
}
