package querybm

import (
	"context"
	"database/sql"
)

type DB interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

var _ DB = (*sql.DB)(nil)
