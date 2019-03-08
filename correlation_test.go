package watermillx

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/subscriber"
)

func TestCorrelationIDPublisherDecorator(t *testing.T) {
	var fn CorrelationIDExtractor = func(msg *message.Message) (string, bool) {
		return "id", true
	}

	publisherStub := newPublisherStub()

	publisher, err := CorrelationIDPublisherDecorator(fn)(publisherStub)
	if err != nil {
		t.Fatal(err)
	}

	msg := message.NewMessage("uuid", []byte{1, 2, 3})

	const topic = "topic"

	err = publisher.Publish(topic, msg)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := middleware.MessageCorrelationID(publisherStub.messages[topic][0]), "id"; got != want {
		t.Errorf("message correlation ID does not match the expected\nactual:  %s\nexpected: %s", got, want)
	}
}

type subscriberStub struct {
	ch chan *message.Message
}

func (m subscriberStub) Subscribe(context.Context, string) (<-chan *message.Message, error) {
	return m.ch, nil
}

func (m subscriberStub) Close() error {
	close(m.ch)
	return nil
}

func TestMessageTransformSubscriberDecorator_transparent(t *testing.T) {
	var fn CorrelationIDInserter = func(msg *message.Message, id string) {
		msg.Metadata.Set("mycid", id)
	}
	sub := subscriberStub{make(chan *message.Message)}
	decorated, err := CorrelationIDSubscriberDecorator(fn)(sub)
	if err != nil {
		t.Fatal(err)
	}

	messages, err := decorated.Subscribe(context.Background(), "topic")
	if err != nil {
		t.Fatal(err)
	}

	richMessage := message.NewMessage("uuid", []byte("serious payloads"))
	richMessage.Metadata.Set("k1", "v1")
	richMessage.Metadata.Set("k2", "v2")
	middleware.SetCorrelationID("id", richMessage)

	go func() {
		sub.ch <- richMessage
	}()

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("expected to read all messages")
	}

	if got, want := received[0].Metadata.Get("mycid"), "id"; got != want {
		t.Errorf("unexpected cid\nactual:  %s\nexpected: %s", got, want)
	}
}
