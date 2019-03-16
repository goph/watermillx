package watermillx

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func TestStructNameMarshaler_Marshal(t *testing.T) {
	message := NamedMessage{}

	marshaler := NewStructNameMarshaler(cqrs.JSONMarshaler{})

	msg, err := marshaler.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := marshaler.NameFromMessage(msg), "NamedMessage"; got != want {
		t.Errorf("message name does not match the expected one\nactual:   %s\nexpected: %s", got, want)
	}
}

func TestStructNameMarshaler_Unmarshal(t *testing.T) {
	type message struct {
		Key string
	}

	m := message{
		Key: "value",
	}

	jsonMarshaler := cqrs.JSONMarshaler{}
	marshaler := NewStructNameMarshaler(jsonMarshaler)

	msg, err := jsonMarshaler.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	var m2 message

	err = marshaler.Unmarshal(msg, &m2)
	if err != nil {
		t.Fatal(err)
	}

	if m != m2 {
		t.Errorf("unmarshaled message does not match the original one\nactual:   %+v\nexpected: %+v", m, m2)
	}
}

func TestStructNameMarshaler_Name(t *testing.T) {
	message := NamedMessage{}

	marshaler := NewStructNameMarshaler(cqrs.JSONMarshaler{})

	if got, want := marshaler.Name(message), "NamedMessage"; got != want {
		t.Errorf("message name does not match the expected one\nactual:   %s\nexpected: %s", got, want)
	}
}
