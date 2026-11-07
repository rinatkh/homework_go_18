//go:build acceptance

package grpcapi

import (
	"context"
	"github.com/rinatkh/homework_go_18/internal/document"
	documentv1 "github.com/rinatkh/homework_go_18/pkg/proto/document/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

type repoStub struct{}

func (repoStub) Create(context.Context, string) (document.Document, error) {
	return document.Document{}, nil
}
func (repoStub) Get(ctx context.Context, id int64) (document.Document, error) {
	if id == 9 {
		return document.Document{}, document.ErrNotFound
	}
	return document.Document{ID: id, Title: "existing", Status: "ready"}, nil
}
func (repoStub) MarkReady(context.Context, int64) error { return nil }
func TestAcceptanceGetDocumentMapping(t *testing.T) {
	s := &Server{UseCase: document.NewUseCase(repoStub{}, nil)}
	for _, tc := range []struct {
		name string
		id   int64
		code codes.Code
	}{{"positive", 1, codes.OK}, {"zero", 0, codes.InvalidArgument}, {"negative", -1, codes.InvalidArgument}, {"absent", 9, codes.NotFound}} {
		t.Run(tc.name, func(t *testing.T) {
			r, err := s.GetDocument(context.Background(), &documentv1.GetDocumentRequest{Id: tc.id})
			if status.Code(err) != tc.code {
				t.Fatalf("TODO-05: got %v want %v", status.Code(err), tc.code)
			}
			if tc.code == codes.OK && (r.GetId() != 1 || r.GetTitle() != "existing") {
				t.Fatal("DTO mismatch")
			}
		})
	}
}
