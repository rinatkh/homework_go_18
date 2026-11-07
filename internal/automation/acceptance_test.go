//go:build acceptance

package automation

import (
	"context"
	documentv1 "github.com/rinatkh/homework_go_18/pkg/proto/document/v1"
	"google.golang.org/grpc"
	"testing"
	"time"
)

type clientStub struct {
	called bool
	ctx    context.Context
}

func (c *clientStub) GetDocument(ctx context.Context, r *documentv1.GetDocumentRequest, o ...grpc.CallOption) (*documentv1.DocumentResponse, error) {
	c.called = true
	c.ctx = ctx
	deadline, ok := ctx.Deadline()
	if !ok || time.Until(deadline) > 3*time.Second {
		panic("missing bounded RPC")
	}
	return &documentv1.DocumentResponse{Id: r.GetId()}, nil
}
func TestAcceptanceTypedGetPreservesCallerAndCleansUp(t *testing.T) {
	c := &clientStub{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := Get(ctx, c, 42)
	if err != nil || r.GetId() != 42 || !c.called {
		t.Fatalf("TODO-06 is incomplete: %v", err)
	}
	if c.ctx.Err() != context.Canceled {
		t.Fatal("RPC context was not canceled after return")
	}
}
