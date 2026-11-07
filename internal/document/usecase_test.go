package document

import (
	"context"
	"testing"
)

type fakeRepo struct{ calls int }

func (r *fakeRepo) Create(context.Context, string) (Document, error) {
	r.calls++
	return Document{}, nil
}
func (r *fakeRepo) Get(context.Context, int64) (Document, error) {
	r.calls++
	return Document{}, ErrNotFound
}
func (r *fakeRepo) MarkReady(context.Context, int64) error { r.calls++; return nil }

type fakePublisher struct{ calls int }

func (p *fakePublisher) Publish(context.Context, Event) error { p.calls++; return nil }
func TestEmptyTitleDoesNotTouchDependencies(t *testing.T) {
	repo := &fakeRepo{}
	pub := &fakePublisher{}
	_, err := NewUseCase(repo, pub).Create(context.Background(), "  ")
	if err != ErrEmptyTitle || repo.calls != 0 || pub.calls != 0 {
		t.Fatal("validation must precede adapters")
	}
}

// TODO-08: add successful Create, repository/publisher errors and worker fake tests.
