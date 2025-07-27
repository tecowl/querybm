package querybm

import (
	"context"
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

type MockRows struct {
	close func() error
	err   func() error
	next  func() bool
	scan  func(dest ...any) error
}

var _ Rows = (*MockRows)(nil)

// Close implements Rows.
func (m *MockRows) Close() error {
	if m.close != nil {
		return m.close()
	}
	return nil
}

// Err implements Rows.
func (m *MockRows) Err() error {
	if m.err != nil {
		return m.err()
	}
	return nil
}

// Next implements Rows.
func (m *MockRows) Next() bool {
	if m.next != nil {
		return m.next()
	}
	return false
}

// Scan implements Rows.
func (m *MockRows) Scan(dest ...any) error {
	if m.scan != nil {
		return m.scan(dest...)
	}
	return nil
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

	row := &MockRow{
		err:  fmt.Errorf("runtime row error"),                        // nolint:err113,perfsprint
		scan: func(...any) error { return fmt.Errorf("scan error") }, // nolint:err113,perfsprint
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
}

func TestStmtQueryContextError(t *testing.T) {
	t.Parallel()

	stmt := &MockStmt{
		queryContext: func(context.Context, ...any) (Rows, error) {
			return nil, fmt.Errorf("runtime rows error") // nolint:err113,perfsprint
		},
	}
	db := &MockDB{
		PrepareContextFunc: func(context.Context, string) (Stmt, error) { return stmt, nil },
	}

	q := &Query[any]{
		db:     db,
		Table:  "users",
		Fields: NewFields[any]([]string{"id", "name"}, nil),
	}
	t.Run("List", func(t *testing.T) {
		t.Parallel()
		if _, err := q.List(t.Context()); err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestStmtQueryContextRowsError(t *testing.T) {
	t.Parallel()

	rows := &MockRows{
		err: func() error {
			return fmt.Errorf("runtime rows error") // nolint:err113,perfsprint
		},
	}
	stmt := &MockStmt{
		queryContext: func(context.Context, ...any) (Rows, error) {
			return rows, nil
		},
	}
	db := &MockDB{
		PrepareContextFunc: func(context.Context, string) (Stmt, error) { return stmt, nil },
	}

	q := &Query[any]{
		db:     db,
		Table:  "users",
		Fields: NewFields[any]([]string{"id", "name"}, nil),
	}
	t.Run("List", func(t *testing.T) {
		t.Parallel()
		if _, err := q.List(t.Context()); err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

type errorLimitOffset struct {
	err error
}

var (
	_ LimitOffset = (*errorLimitOffset)(nil)
	_ Validatable = (*errorLimitOffset)(nil)
)

func (m *errorLimitOffset) Build(*Statement) {
}

func (m *errorLimitOffset) Validate() error {
	return m.err
}

func TestQueryValidateError(t *testing.T) {
	t.Parallel()

	q := &Query[any]{
		db:          &MockDB{},
		Table:       "users",
		Fields:      NewFields[any]([]string{"id", "name"}, nil),
		LimitOffset: &errorLimitOffset{err: fmt.Errorf("limit offset error")}, // nolint:err113,perfsprint
	}

	t.Run("LimitOffset", func(t *testing.T) {
		t.Parallel()
		if err := q.Validate(); err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
