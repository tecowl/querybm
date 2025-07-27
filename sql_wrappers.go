package querybm

import (
	"context"
	"database/sql"
)

type DB interface {
	PrepareContext(ctx context.Context, query string) (Stmt, error)
}

type DBWrapper struct {
	db *sql.DB
}

var _ DB = (*DBWrapper)(nil)

func newDBWrapper(db *sql.DB) *DBWrapper {
	return &DBWrapper{db: db}
}

func (w *DBWrapper) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	stmt, err := w.db.PrepareContext(ctx, query)
	return stmt, err
}

type Stmt interface {
	Close() error
	QueryRowContext(ctx context.Context, args ...any) *sql.Row
	QueryContext(ctx context.Context, args ...any) (*sql.Rows, error)
}

var _ Stmt = (*sql.Stmt)(nil)
