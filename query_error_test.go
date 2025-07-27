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
