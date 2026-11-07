package messaging

import (
	"context"
	"encoding/json"
	"github.com/rinatkh/homework_go_18/internal/document"
	kafka "github.com/segmentio/kafka-go"
)

type Publisher struct{ Writer *kafka.Writer }

func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{Writer: &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: topic, Balancer: &kafka.Hash{}, RequiredAcks: kafka.RequireAll}}
}

// TODO-02: serialize JSON, key=document_id, bounded context, synchronous ack.
func (p *Publisher) Publish(ctx context.Context, e document.Event) error {
	return document.ErrNotImplemented
}
func Decode(data []byte) (document.Event, error) {
	var e document.Event
	err := json.Unmarshal(data, &e)
	return e, err
}

// TODO-03: build reader with group, Fetch -> Handle -> Commit. Ignore other types.
// Stop on ctx cancellation. A transient handler/publish error must not commit offset.
func Consume(ctx context.Context, brokers []string, topic, group, eventType string, handle func(context.Context, document.Event) error) error {
	return document.ErrNotImplemented
}
