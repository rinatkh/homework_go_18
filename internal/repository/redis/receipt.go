package redis

import (
	"context"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rinatkh/homework_go_18/internal/document"
	"time"
)

type ReceiptStore struct{ Client *goredis.Client }

var _ document.ReceiptStore = (*ReceiptStore)(nil)

// TODO-01: use doc:receipt:<event_id>, Exists and SET with positive TTL. Pass ctx.
func (r *ReceiptStore) Exists(ctx context.Context, id string) (bool, error) {
	return false, document.ErrNotImplemented
}
func (r *ReceiptStore) Save(ctx context.Context, id string, ttl time.Duration) error {
	return document.ErrNotImplemented
}
