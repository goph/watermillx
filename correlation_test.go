package watermillx

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/subscriber"
)

func TestCorrelationIDPublisherDecorator(t *testing.T) {
	var fn CorrelationIDExtractor = func(msg *message.Message) (string, bool) {
		return "id", true
	}

	pubsub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})

	publisher, err := CorrelationIDPublisherDecorator(fn)(pubsub)
	if err != nil {
		t.Fatal(err)
	}

	const topic = "topic"

	messages, err := pubsub.Subscribe(context.Background(), topic)
	if err != nil {
		t.Fatal(err)
	}

	msg := message.NewMessage("uuid", []byte{1, 2, 3})

	err = publisher.Publish(topic, msg)
	if err != nil {
		t.Fatal(err)
	}

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("no message received")
	}

	if got, want := middleware.MessageCorrelationID(received[0]), "id"; got != want {
		t.Errorf("message correlation ID does not match the expected\nactual:   %s\nexpected: %s", got, want)
	}
}

func TestCorrelationIDSubscriberDecorator(t *testing.T) {
	var fn CorrelationIDInserter = func(msg *message.Message, id string) {
		msg.Metadata.Set("mycid", id)
	}

	pubsub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})

	sub, err := CorrelationIDSubscriberDecorator(fn)(pubsub)
	if err != nil {
		t.Fatal(err)
	}

	const topic = "topic"

	messages, err := sub.Subscribe(context.Background(), topic)
	if err != nil {
		t.Fatal(err)
	}

	msg := message.NewMessage("uuid", []byte{1, 2, 3})
	middleware.SetCorrelationID("id", msg)

	err = pubsub.Publish(topic, msg)
	if err != nil {
		t.Fatal(err)
	}

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("no message received")
	}

	if got, want := received[0].Metadata.Get("mycid"), "id"; got != want {
		t.Errorf("message correlation ID does not match the expected\nactual:   %s\nexpected: %s", got, want)
	}
}
