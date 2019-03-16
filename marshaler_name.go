package watermillx

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

// NameMarshaler retrieves the name from a message implementing the following interface:
//		type namedMessage interface {
//			Name() string
//		}
type NameMarshaler struct {
	marshaler cqrs.CommandEventMarshaler
}

// NewNameMarshaler returns a new NameMarshaler instance.
func NewNameMarshaler(marshaler cqrs.CommandEventMarshaler) *NameMarshaler {
	return &NameMarshaler{
		marshaler: marshaler,
	}
}

type namedMessage interface {
	Name() string
}

func (m *NameMarshaler) Marshal(v interface{}) (*message.Message, error) {
	msg, err := m.marshaler.Marshal(v)
	if err != nil {
		return nil, err
	}

	msg.Metadata.Set("name", m.Name(v))

	return msg, nil
}

func (m *NameMarshaler) Unmarshal(msg *message.Message, v interface{}) error {
	err := m.marshaler.Unmarshal(msg, v)
	if err != nil {
		return err
	}

	return nil
}

func (m *NameMarshaler) Name(v interface{}) string {
	if v, ok := v.(namedMessage); ok {
		return v.Name()
	}

	return m.marshaler.Name(v)
}

func (m *NameMarshaler) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
