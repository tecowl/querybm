package querybm

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
)

type MockDB struct {
	PrepareContextFunc func(ctx context.Context, query string) (Stmt, error)
}

func (m *MockDB) PrepareContext(ctx context.Context, query string) (Stmt, error) { // nolint:ireturn
	if m.PrepareContextFunc != nil {
		return m.PrepareContextFunc(ctx, query)
	}
	return nil, nil // nolint:nilnil
}

type MockStmt struct {
	close           func() error
	queryRowContext func(ctx context.Context, args ...any) Row
	queryContext    func(ctx context.Context, args ...any) (Rows, error)
}

var _ Stmt = (*MockStmt)(nil)

func (m *MockStmt) Close() error {
	if m.close != nil {
		return m.close()
	}
	return nil
}

func (m *MockStmt) QueryContext(ctx context.Context, args ...any) (Rows, error) { // nolint:ireturn
	if m.queryContext != nil {
		return m.queryContext(ctx, args...)
	}
	return nil, nil // nolint:nilnil
}

func (m *MockStmt) QueryRowContext(ctx context.Context, args ...any) Row { // nolint:ireturn
	if m.queryRowContext != nil {
		return m.queryRowContext(ctx, args...)
	}
	return nil
}

type MockRow struct {
	err  error
	scan func(dest ...any) error
}

func (m *MockRow) Err() error {
	return m.err
}

func (m *MockRow) Scan(dest ...any) error {
	if m.scan != nil {
		return m.scan(dest...)
	}
	return nil // nolint:nilnil
}

func TestQueryPrepareContextError(t *testing.T) {
	t.Parallel()

	db := &MockDB{
		PrepareContextFunc: func(context.Context, string) (Stmt, error) {
			return nil, fmt.Errorf("runtime error") // nolint:err113,perfsprint
		},
	}

	q := &Query[any]{
		db:     db,
		Table:  "users",
		Fields: NewFields[any]([]string{"id", "name"}, nil),
	}

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		_, err := q.List(t.Context())
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("First", func(t *testing.T) {
		t.Parallel()
		_, err := q.First(t.Context())
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestStmtQueryRowContextError(t *testing.T) {
	t.Parallel()

	scanErrors := []error{
		sql.ErrNoRows,
		fmt.Errorf("scan error"), // nolint:err113,perfsprint
	}

	for _, scanErr := range scanErrors {
		t.Run(scanErr.Error(), func(t *testing.T) {
			t.Parallel()

			row := &MockRow{
				err:  fmt.Errorf("runtime row error"), // nolint:err113,perfsprint
				scan: func(...any) error { return scanErr },
			}
			stmt := &MockStmt{
				queryRowContext: func(context.Context, ...any) Row { return row },
			}
			db := &MockDB{
				PrepareContextFunc: func(context.Context, string) (Stmt, error) { return stmt, nil },
			}

			q := &Query[any]{
				db:     db,
				Table:  "users",
				Fields: NewFields[any]([]string{"id", "name"}, nil),
			}

			t.Run("Count", func(t *testing.T) {
				t.Parallel()
				if _, err := q.Count(t.Context()); err == nil {
					t.Error("Expected error, got nil")
				}
			})
			t.Run("First", func(t *testing.T) {
				t.Parallel()
				if _, err := q.First(t.Context()); err == nil {
					t.Error("Expected error, got nil")
				}
			})
		})
	}
}
