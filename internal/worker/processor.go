package worker

import (
	"context"
	"github.com/rinatkh/homework_go_18/internal/document"
)

type Processor struct {
	Publisher document.Publisher
	Receipts  document.ReceiptStore
}

// TODO-04: document.created -> word analysis -> document.processed.
// Check receipt; after successful result publish save a 10min receipt. Failures retry.
// Duplicate publication is still possible across crashes/TTL; owner DB update is idempotent.
func (p *Processor) Handle(ctx context.Context, event document.Event) error {
	return document.ErrNotImplemented
}
