package watermillx

import (
	"fmt"
	"strings"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

// StructNameMarshaler uses the struct name (without the package) as the message name.
type StructNameMarshaler struct {
	marshaler cqrs.CommandEventMarshaler
}

// NewStructNameMarshaler returns a new StructNameMarshaler instance.
func NewStructNameMarshaler(marshaler cqrs.CommandEventMarshaler) *StructNameMarshaler {
	return &StructNameMarshaler{
		marshaler: marshaler,
	}
}

func (m *StructNameMarshaler) Marshal(v interface{}) (*message.Message, error) {
	msg, err := m.marshaler.Marshal(v)
	if err != nil {
		return nil, err
	}

	msg.Metadata.Set("name", m.Name(v))

	return msg, nil
}

func (m *StructNameMarshaler) Unmarshal(msg *message.Message, v interface{}) error {
	err := m.marshaler.Unmarshal(msg, v)
	if err != nil {
		return err
	}

	return nil
}

func (m *StructNameMarshaler) Name(v interface{}) string {
	segments := strings.Split(fmt.Sprintf("%T", v), ".")

	return segments[len(segments)-1]
}

func (m *StructNameMarshaler) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
