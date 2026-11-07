//go:build acceptance

package worker

import (
	"context"
	"github.com/rinatkh/homework_go_18/internal/document"
	"testing"
	"time"
)

type publisherStub struct {
	events []document.Event
	err    error
}

func (p *publisherStub) Publish(ctx context.Context, e document.Event) error {
	if p.err != nil {
		return p.err
	}
	p.events = append(p.events, e)
	return nil
}

type receiptStub struct {
	saved map[string]bool
	ttl   time.Duration
}

func (r *receiptStub) Exists(ctx context.Context, id string) (bool, error) { return r.saved[id], nil }
func (r *receiptStub) Save(ctx context.Context, id string, ttl time.Duration) error {
	r.saved[id] = true
	r.ttl = ttl
	return nil
}
func TestAcceptanceProcessorPublishesAndDeduplicates(t *testing.T) {
	for _, title := range []string{"one", "one two", "  one\t two\nthree  ", "слова на русском"} {
		t.Run(title, func(t *testing.T) {
			p := &publisherStub{}
			r := &receiptStub{saved: map[string]bool{}}
			processor := &Processor{Publisher: p, Receipts: r}
			event := document.Event{ID: "event-1", Type: "document.created", DocumentID: 3, Title: title}
			if err := processor.Handle(context.Background(), event); err != nil {
				t.Fatalf("TODO-04 is incomplete: %v", err)
			}
			if len(p.events) != 1 || p.events[0].Type != "document.processed" || p.events[0].DocumentID != 3 || p.events[0].Words <= 0 || !r.saved[event.ID] || r.ttl != 10*time.Minute {
				t.Fatal("result/receipt contract")
			}
			if err := processor.Handle(context.Background(), event); err != nil || len(p.events) != 1 {
				t.Fatal("duplicate receipt did not suppress work", err)
			}
		})
	}
}
