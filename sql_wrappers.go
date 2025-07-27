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

func (w *DBWrapper) PrepareContext(ctx context.Context, query string) (Stmt, error) { // nolint:ireturn
	stmt, err := w.db.PrepareContext(ctx, query)
	return newStmtWrapper(stmt), err
}

type Stmt interface {
	Close() error
	QueryRowContext(ctx context.Context, args ...any) Row
	QueryContext(ctx context.Context, args ...any) (Rows, error)
}

type StmtWrapper struct {
	*sql.Stmt
}

var _ Stmt = (*StmtWrapper)(nil)

func newStmtWrapper(stmt *sql.Stmt) *StmtWrapper {
	return &StmtWrapper{Stmt: stmt}
}

func (s *StmtWrapper) QueryRowContext(ctx context.Context, args ...any) Row { // nolint:ireturn
	row := s.Stmt.QueryRowContext(ctx, args...)
	return row
}

func (s *StmtWrapper) QueryContext(ctx context.Context, args ...any) (Rows, error) { // nolint:ireturn
	rows, err := s.Stmt.QueryContext(ctx, args...) // nolint:rowserrcheck
	return rows, err
}

type Row interface {
	Err() error
	Scan(dest ...any) error
}

var _ Row = (*sql.Row)(nil)

type Rows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...any) error
}

var _ Rows = (*sql.Rows)(nil)
