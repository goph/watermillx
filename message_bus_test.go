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

	if got, want := publisher.messages[topic][0].UUID, idGenerator.id; got != want {
		t.Errorf("message id does not match the expected value\nactual:  %s\nexpected: %s", got, want)
	}

	if got, want := string(publisher.messages[topic][0].Payload), marshaler.message; got != want {
		t.Errorf("message payload does not match the expected value\nactual:  %s\nexpected: %s", got, want)
	}
}
