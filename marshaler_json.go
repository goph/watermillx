package watermillx

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// JSONMessageMarshaler marshals/unmarshals messages.
type JSONMessageMarshaler struct{}

// NewJSONMessageMarshaler returns a new JSONMessageMarshaler instance.
func NewJSONMessageMarshaler() *JSONMessageMarshaler {
	return &JSONMessageMarshaler{}
}

// Marshal serializes a message.
func (*JSONMessageMarshaler) Marshal(msg interface{}) ([]byte, error) {
	payload, err := json.Marshal(msg)

	return payload, errors.Wrap(err, "failed to marshal message")
}
