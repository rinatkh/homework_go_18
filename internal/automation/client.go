package automation

import (
	"context"
	"github.com/rinatkh/homework_go_18/internal/document"
	documentv1 "github.com/rinatkh/homework_go_18/pkg/proto/document/v1"
)

// TODO-06: typed GetDocument call with a 3s deadline, caller cancellation preserved.
func Get(ctx context.Context, client documentv1.DocumentServiceClient, id int64) (*documentv1.DocumentResponse, error) {
	return nil, document.ErrNotImplemented
}
