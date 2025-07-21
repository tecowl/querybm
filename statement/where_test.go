package statement

import (
	"reflect"
	"testing"

	"github.com/tecowl/querybm/expr"
)

func TestNewWhere(t *testing.T) {
	tests := []struct {
		name      string
		connector string
	}{
		{
			name:      "Create WHERE with AND connector",
			connector: " AND ",
		},
		{
			name:      "Create WHERE with OR connector",
			connector: " OR ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newWhere(tt.connector)
			if w.Connector != tt.connector {
				t.Errorf("newWhere() Connector = %v, want %v", w.Connector, tt.connector)
			}
			if len(w.conditions) != 0 {
				t.Errorf("newWhere() conditions = %v, want empty slice", w.conditions)
			}
		})
	}
}

func TestWhereBlock_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *WhereBlock
		want  bool
	}{
		{
			name: "Empty WHERE block",
			setup: func() *WhereBlock {
				return newWhere(" AND ")
			},
			want: true,
		},
		{
			name: "WHERE block with one condition",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("status", expr.Eq("active")))
				return w
			},
			want: false,
		},
		{
			name: "WHERE block with multiple conditions",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("status", expr.Eq("active")))
				w.Add(expr.Field("age", expr.Gte(18)))
				return w
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.setup()
			if got := w.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereBlock_Add(t *testing.T) {
	tests := []struct {
		name       string
		connector  string
		conditions []expr.ConditionExpr
		wantCount  int
	}{
		{
			name:      "Add single condition",
			connector: " AND ",
			conditions: []expr.ConditionExpr{
				expr.Field("name", expr.Eq("John")),
			},
			wantCount: 1,
		},
		{
			name:      "Add multiple conditions",
			connector: " AND ",
			conditions: []expr.ConditionExpr{
				expr.Field("name", expr.Eq("John")),
				expr.Field("age", expr.Gte(18)),
				expr.Field("status", expr.In("active", "pending")),
			},
			wantCount: 3,
		},
		{
			name:      "Add nested conditions",
			connector: " AND ",
			conditions: []expr.ConditionExpr{
				expr.Field("status", expr.Eq("active")),
				expr.Or(
					expr.Field("role", expr.Eq("admin")),
					expr.Field("role", expr.Eq("moderator")),
				),
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newWhere(tt.connector)
			for _, cond := range tt.conditions {
				w.Add(cond)
			}
			if len(w.conditions) != tt.wantCount {
				t.Errorf("Add() condition count = %v, want %v", len(w.conditions), tt.wantCount)
			}
		})
	}
}

func TestWhereBlock_Build(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *WhereBlock
		wantSQL    string
		wantValues []any
	}{
		{
			name: "Empty WHERE block",
			setup: func() *WhereBlock {
				return newWhere(" AND ")
			},
			wantSQL:    "",
			wantValues: []any{},
		},
		{
			name: "Single condition",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("status", expr.Eq("active")))
				return w
			},
			wantSQL:    "status = ?",
			wantValues: []any{"active"},
		},
		{
			name: "Multiple conditions with AND",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("name", expr.Eq("John")))
				w.Add(expr.Field("age", expr.Gte(18)))
				return w
			},
			wantSQL:    "name = ? AND age >= ?",
			wantValues: []any{"John", 18},
		},
		{
			name: "Multiple conditions with OR",
			setup: func() *WhereBlock {
				w := newWhere(" OR ")
				w.Add(expr.Field("status", expr.Eq("active")))
				w.Add(expr.Field("status", expr.Eq("pending")))
				return w
			},
			wantSQL:    "status = ? OR status = ?",
			wantValues: []any{"active", "pending"},
		},
		{
			name: "Complex conditions with IN",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("id", expr.In(1, 2, 3)))
				w.Add(expr.Field("status", expr.NotEq("deleted")))
				return w
			},
			wantSQL:    "id IN (?,?,?) AND status <> ?",
			wantValues: []any{1, 2, 3, "deleted"},
		},
		{
			name: "Conditions with LIKE",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("email", expr.LikeContains("@example")))
				w.Add(expr.Field("name", expr.LikeStartsWith("John")))
				return w
			},
			wantSQL:    "email LIKE ? AND name LIKE ?",
			wantValues: []any{"%@example%", "John%"},
		},
		{
			name: "Conditions with NULL checks",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("deleted_at", expr.IsNull()))
				w.Add(expr.Field("verified_at", expr.IsNotNull()))
				return w
			},
			wantSQL:    "deleted_at IS NULL AND verified_at IS NOT NULL",
			wantValues: []any{},
		},
		{
			name: "Nested conditions",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Field("active", expr.Eq(true)))
				w.Add(expr.Or(
					expr.Field("role", expr.Eq("admin")),
					expr.Field("role", expr.Eq("moderator")),
				))
				return w
			},
			wantSQL:    "active = ? AND (role = ? OR role = ?)",
			wantValues: []any{true, "admin", "moderator"},
		},
		{
			name: "Complex nested conditions",
			setup: func() *WhereBlock {
				w := newWhere(" AND ")
				w.Add(expr.Or(
					expr.And(
						expr.Field("status", expr.Eq("active")),
						expr.Field("verified", expr.Eq(true)),
					),
					expr.And(
						expr.Field("status", expr.Eq("pending")),
						expr.Field("admin_approved", expr.Eq(true)),
					),
				))
				w.Add(expr.Field("deleted_at", expr.IsNull()))
				return w
			},
			wantSQL:    "((status = ? AND verified = ?) OR (status = ? AND admin_approved = ?)) AND deleted_at IS NULL",
			wantValues: []any{"active", true, "pending", true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.setup()
			gotSQL, gotValues := w.Build()
			if gotSQL != tt.wantSQL {
				t.Errorf("Build() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Build() values = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}
