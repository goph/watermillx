package watermillx

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// CorrelationIDExtractor extracts a correlation ID from a message (if any).
type CorrelationIDExtractor func(msg *message.Message) (string, bool)

// ContextCorrelationIDExtractor extracts a correlation ID from a context (if any).
type ContextCorrelationIDExtractor func(ctx context.Context) (string, bool)

// ContextCorrelationIDExtractorFunc wraps a context correlation ID extractor.
func ContextCorrelationIDExtractorFunc(fn ContextCorrelationIDExtractor) CorrelationIDExtractor {
	return func(msg *message.Message) (string, bool) {
		ctx := msg.Context()

		return fn(ctx)
	}
}

// CorrelationIDPublisherDecorator creates a publisher decorator that extracts correlation ID
// from each message context that passes through the publisher.
func CorrelationIDPublisherDecorator(extractor CorrelationIDExtractor) message.PublisherDecorator {
	return func(pub message.Publisher) (message.Publisher, error) {
		return &cidPublisherDecorator{
			publisher: pub,
			extractor: extractor,
		}, nil
	}
}

type cidPublisherDecorator struct {
	publisher message.Publisher
	extractor CorrelationIDExtractor
}

func (d *cidPublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	for _, msg := range messages {
		cid, ok := d.extractor(msg)
		if ok {
			middleware.SetCorrelationID(cid, msg)
		}
	}

	return d.publisher.Publish(topic, messages...)
}

func (d *cidPublisherDecorator) Close() error {
	return d.publisher.Close()
}

// CorrelationIDInserter inserts a correlation ID into a message.
type CorrelationIDInserter func(msg *message.Message, id string)

// ContextCorrelationIDInserter inserts a correlation ID into a context.
type ContextCorrelationIDInserter func(ctx context.Context, id string) context.Context

// ContextCorrelationIDInserterFunc wraps a context correlation ID inserter.
func ContextCorrelationIDInserterFunc(fn ContextCorrelationIDInserter) CorrelationIDInserter {
	return func(msg *message.Message, id string) {
		msg.SetContext(fn(msg.Context(), id))
	}
}

// CorrelationIDSubscriberDecorator creates a subscriber decorator that inserts a correlation ID
// into each message context that passes through the subscriber.
func CorrelationIDSubscriberDecorator(inserter CorrelationIDInserter) message.SubscriberDecorator {
	return func(sub message.Subscriber) (message.Subscriber, error) {
		d := &cidSubscriberDecorator{}

		var err error
		d.subscriber, err = message.MessageTransformSubscriberDecorator(func(msg *message.Message) {
			cid := middleware.MessageCorrelationID(msg)
			if cid != "" {
				inserter(msg, cid)
			}
		})(sub)

		return d, err
	}
}

type cidSubscriberDecorator struct {
	subscriber message.Subscriber
}

func (d *cidSubscriberDecorator) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return d.subscriber.Subscribe(ctx, topic)
}

func (d *cidSubscriberDecorator) Close() error {
	return d.subscriber.Close()
}
