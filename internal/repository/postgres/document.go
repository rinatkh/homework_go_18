package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinatkh/homework_go_18/internal/document"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool} }
func (r *Repository) Create(ctx context.Context, title string) (document.Document, error) {
	var d document.Document
	err := r.pool.QueryRow(ctx, `INSERT INTO documents(title) VALUES($1) RETURNING id,title,status`, title).Scan(&d.ID, &d.Title, &d.Status)
	return d, err
}
func (r *Repository) Get(ctx context.Context, id int64) (document.Document, error) {
	var d document.Document
	err := r.pool.QueryRow(ctx, `SELECT id,title,status FROM documents WHERE id=$1`, id).Scan(&d.ID, &d.Title, &d.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return document.Document{}, document.ErrNotFound
	}
	return d, err
}
func (r *Repository) MarkReady(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `UPDATE documents SET status='ready' WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return document.ErrNotFound
	}
	return nil
}
