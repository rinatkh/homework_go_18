package document

import (
	"context"
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("TODO: infrastructure exercise is not implemented")
var ErrNotFound = errors.New("document not found")
var ErrEmptyTitle = errors.New("empty title")

type Document struct {
	ID     int64
	Title  string
	Status string
}
type Event struct {
	ID         string `json:"event_id"`
	Type       string `json:"event_type"`
	DocumentID int64  `json:"document_id"`
	Title      string `json:"title,omitempty"`
	Words      int    `json:"words,omitempty"`
}
type Repository interface {
	Create(context.Context, string) (Document, error)
	Get(context.Context, int64) (Document, error)
	MarkReady(context.Context, int64) error
}
type Publisher interface {
	Publish(context.Context, Event) error
}

// Receipt — временная отметка об успешно опубликованном результате, не business DB.
type ReceiptStore interface {
	Exists(context.Context, string) (bool, error)
	Save(context.Context, string, time.Duration) error
}
