package watermillx

import (
	"context"
	"testing"
)

type messageMarshalerStub struct {
	message string
}

func (m *messageMarshalerStub) Marshal(msg interface{}) ([]byte, error) {
	return []byte(m.message), nil
}

type idGeneratorStub struct {
	id string
}

func (g *idGeneratorStub) Generate() (string, error) {
	return g.id, nil
}

func TestMessageBus_Publish(t *testing.T) {
	publisher := newPublisherStub()
	marshaler := &messageMarshalerStub{
		message: "message",
	}
	idGenerator := &idGeneratorStub{
		id: "id",
	}
	messageBus := NewMessageBus(publisher, marshaler, idGenerator)

	const topic = "topic"

	err := messageBus.Publish(context.Background(), topic, "message")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(publisher.messages[topic][0].Payload), "message"; got != want {
		t.Errorf("unexpected message\nactual:  %s\nexpected: %s", got, want)
	}
}
