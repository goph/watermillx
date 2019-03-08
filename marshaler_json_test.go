package watermillx

import (
	"testing"
)

type messageStub struct {
	Key string
}

func TestJSONMessageMarshaler_Marshal(t *testing.T) {
	marshaler := NewJSONMessageMarshaler()

	msg := messageStub{
		Key: "value",
	}

	payload, err := marshaler.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(payload), `{"Key":"value"}`; got != want {
		t.Errorf("marshaling message failed\nactual:  %s\nexpected: %s", got, want)
	}
}
