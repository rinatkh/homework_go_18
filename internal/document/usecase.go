package document

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

type UseCase struct {
	repo      Repository
	publisher Publisher
}

func NewUseCase(repo Repository, publisher Publisher) *UseCase { return &UseCase{repo, publisher} }
func (u *UseCase) Create(ctx context.Context, title string) (Document, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Document{}, ErrEmptyTitle
	}
	value, err := u.repo.Create(ctx, title)
	if err != nil {
		return Document{}, err
	}
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return Document{}, err
	}
	if err := u.publisher.Publish(ctx, Event{ID: hex.EncodeToString(id[:]), Type: "document.created", DocumentID: value.ID, Title: value.Title}); err != nil {
		return Document{}, err
	}
	return value, nil
}
func (u *UseCase) Get(ctx context.Context, id int64) (Document, error) { return u.repo.Get(ctx, id) }
