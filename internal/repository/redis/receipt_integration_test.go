//go:build integration

package redis

import (
	"context"
	"fmt"
	goredis "github.com/redis/go-redis/v9"
	"os"
	"testing"
	"time"
)

func TestAcceptanceReceiptLifecycleAndTTL(t *testing.T) {
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Fatal("TEST_REDIS_URL required")
	}
	options, err := goredis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	client := goredis.NewClient(options)
	defer client.Close()
	store := &ReceiptStore{Client: client}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	id := fmt.Sprintf("acceptance-%d", time.Now().UnixNano())
	defer client.Del(context.Background(), "doc:receipt:"+id)
	exists, err := store.Exists(ctx, id)
	if err != nil || exists {
		t.Fatalf("TODO-01: missing receipt: %v", err)
	}
	if err = store.Save(ctx, id, 120*time.Millisecond); err != nil {
		t.Fatalf("TODO-01: Save: %v", err)
	}
	exists, err = store.Exists(ctx, id)
	if err != nil || !exists {
		t.Fatal("receipt not saved", err)
	}
	ttl, err := client.PTTL(ctx, "doc:receipt:"+id).Result()
	if err != nil || ttl <= 0 || ttl > 120*time.Millisecond {
		t.Fatal(ttl, err)
	}
	for ctx.Err() == nil {
		exists, err = store.Exists(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if !exists {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("receipt did not expire")
}
