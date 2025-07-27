package statement

import (
	"reflect"
	"testing"

	"github.com/tecowl/querybm/expr"
)

func TestNewStatement(t *testing.T) {
	t.Parallel()
	fields := NewSimpleFields("id", "name", "email")
	s := New("users", fields)

	if s.Table == nil {
		t.Errorf("NewStatement() Table should not be nil")
	}
	{
		r, args := s.Table.Build()
		if r != "users" {
			t.Errorf("NewStatement() Table content = %v, want %v", r, "users")
		}
		if len(args) != 0 {
			t.Errorf("NewStatement() Table values = %v, want empty slice", args)
		}
	}

	if !reflect.DeepEqual(s.Fields, fields) {
		t.Errorf("NewStatement() fields = %v, want %v", s.Fields, fields)
	}
	if s.Where == nil {
		t.Errorf("NewStatement() Where should not be nil")
	}
	if s.Sort == nil {
		t.Errorf("NewStatement() Sort should not be nil")
	}
	if s.LimitOffset == nil {
		t.Errorf("NewStatement() LimitOffset should not be nil")
	}
}

func TestStatement_Build_SimpleSelect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		statement  *Statement
		wantSQL    string
		wantValues []any
	}{
		{
			name: "Simple SELECT without conditions",
			statement: New("users",
				NewSimpleFields("id", "name", "email"),
			),
			wantSQL:    "SELECT id, name, email FROM users",
			wantValues: []any{},
		},
		{
			name: "SELECT with single field",
			statement: New("products",
				NewSimpleFields("count(*)"),
			),
			wantSQL:    "SELECT count(*) FROM products",
			wantValues: []any{},
		},
		{
			name: "SELECT all fields",
			statement: New("orders",
				NewSimpleFields("*"),
			),
			wantSQL:    "SELECT * FROM orders",
			wantValues: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotSQL, gotValues := tt.statement.Build()
			if gotSQL != tt.wantSQL {
				t.Errorf("Build() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Build() values = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestStatement_Build_WithWhere(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func() *Statement
		wantSQL    string
		wantValues []any
	}{
		{
			name: "SELECT with single WHERE condition",
			setup: func() *Statement {
				s := New("users", NewSimpleFields("id", "name"))
				s.Where.Add(expr.Field("status", expr.Eq("active")))
				return s
			},
			wantSQL:    "SELECT id, name FROM users WHERE status = ?",
			wantValues: []any{"active"},
		},
		{
			name: "SELECT with multiple WHERE conditions",
			setup: func() *Statement {
				s := New("users", NewSimpleFields("*"))
				s.Where.Add(expr.Field("age", expr.Gte(18)))
				s.Where.Add(expr.Field("status", expr.Eq("active")))
				return s
			},
			wantSQL:    "SELECT * FROM users WHERE age >= ? AND status = ?",
			wantValues: []any{18, "active"},
		},
		{
			name: "SELECT with IN condition",
			setup: func() *Statement {
				s := New("products", NewSimpleFields("id", "name", "price"))
				s.Where.Add(expr.Field("category_id", expr.In(1, 2, 3)))
				return s
			},
			wantSQL:    "SELECT id, name, price FROM products WHERE category_id IN (?,?,?)",
			wantValues: []any{1, 2, 3},
		},
		{
			name: "SELECT with LIKE condition",
			setup: func() *Statement {
				s := New("users", NewSimpleFields("id", "email"))
				s.Where.Add(expr.Field("email", expr.LikeContains("@example")))
				return s
			},
			wantSQL:    "SELECT id, email FROM users WHERE email LIKE ?",
			wantValues: []any{"%@example%"},
		},
		{
			name: "SELECT with IS NULL condition",
			setup: func() *Statement {
				s := New("users", NewSimpleFields("id", "name"))
				s.Where.Add(expr.Field("deleted_at", expr.IsNull()))
				return s
			},
			wantSQL:    "SELECT id, name FROM users WHERE deleted_at IS NULL",
			wantValues: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.setup()
			gotSQL, gotValues := s.Build()
			if gotSQL != tt.wantSQL {
				t.Errorf("Build() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Build() values = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestStatement_Build_WithSort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func() *Statement
		wantSQL    string
		wantValues []any
	}{
		{
			name: "SELECT with single ORDER BY",
			setup: func() *Statement {
				s := New("users", NewSimpleFields("id", "name"))
				s.Sort.Add("created_at DESC")
				return s
			},
			wantSQL:    "SELECT id, name FROM users ORDER BY created_at DESC",
			wantValues: []any{},
		},
		{
			name: "SELECT with multiple ORDER BY",
			setup: func() *Statement {
				s := New("products", NewSimpleFields("*"))
				s.Sort.Add("category_id")
				s.Sort.Add("price DESC")
				return s
			},
			wantSQL:    "SELECT * FROM products ORDER BY category_id, price DESC",
			wantValues: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.setup()
			gotSQL, gotValues := s.Build()
			if gotSQL != tt.wantSQL {
				t.Errorf("Build() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Build() values = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestStatement_Build_WithLimitOffset(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func() *Statement
		wantSQL    string
		wantValues []any
	}{
		{
			name: "SELECT with LIMIT",
			setup: func() *Statement {
				s := New("users", NewSimpleFields("id", "name"))
				s.LimitOffset.Add("LIMIT ?", 10)
				return s
			},
			wantSQL:    "SELECT id, name FROM users LIMIT ?",
			wantValues: []any{10},
		},
		{
			name: "SELECT with LIMIT and OFFSET",
			setup: func() *Statement {
				s := New("products", NewSimpleFields("*"))
				s.LimitOffset.Add("LIMIT ?", 20)
				s.LimitOffset.Add("OFFSET ?", 40)
				return s
			},
			wantSQL:    "SELECT * FROM products LIMIT ? OFFSET ?",
			wantValues: []any{20, 40},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.setup()
			gotSQL, gotValues := s.Build()
			if gotSQL != tt.wantSQL {
				t.Errorf("Build() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Build() values = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestStatement_Build_Complex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func() *Statement
		wantSQL    string
		wantValues []any
	}{
		{
			name: "SELECT with WHERE, ORDER BY, and LIMIT",
			setup: func() *Statement {
				s := New("users", NewSimpleFields("id", "name", "email"))
				s.Where.Add(expr.Field("status", expr.Eq("active")))
				s.Where.Add(expr.Field("age", expr.Gte(18)))
				s.Sort.Add("created_at DESC")
				s.LimitOffset.Add("LIMIT ?", 10)
				return s
			},
			wantSQL:    "SELECT id, name, email FROM users WHERE status = ? AND age >= ? ORDER BY created_at DESC LIMIT ?",
			wantValues: []any{"active", 18, 10},
		},
		{
			name: "Complex query with multiple conditions and sorting",
			setup: func() *Statement {
				s := New("products", NewSimpleFields("id", "name", "price", "category_id"))
				s.Where.Add(expr.Field("price", expr.Lt(1000)))
				s.Where.Add(expr.Field("category_id", expr.In(1, 2, 3)))
				s.Where.Add(expr.Field("deleted_at", expr.IsNull()))
				s.Sort.Add("category_id")
				s.Sort.Add("price ASC")
				s.LimitOffset.Add("LIMIT ?", 20)
				s.LimitOffset.Add("OFFSET ?", 100)
				return s
			},
			wantSQL:    "SELECT id, name, price, category_id FROM products WHERE price < ? AND category_id IN (?,?,?) AND deleted_at IS NULL ORDER BY category_id, price ASC LIMIT ? OFFSET ?",
			wantValues: []any{1000, 1, 2, 3, 20, 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.setup()
			gotSQL, gotValues := s.Build()
			if gotSQL != tt.wantSQL {
				t.Errorf("Build() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Build() values = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}
