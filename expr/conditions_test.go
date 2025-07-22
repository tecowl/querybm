package expr

import (
	"reflect"
	"testing"
)

func TestConditions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		conditions *Conditions
		wantString string
		wantValues []any
	}{
		{
			name:       "Empty conditions",
			conditions: NewConditions(" AND "),
			wantString: "",
			wantValues: []any{},
		},
		{
			name:       "Single condition",
			conditions: NewConditions(" AND ", Field("name", Eq("John"))),
			wantString: "name = ?",
			wantValues: []any{"John"},
		},
		{
			name: "Multiple conditions with AND",
			conditions: NewConditions(" AND ",
				Field("name", Eq("John")),
				Field("age", Gt(18)),
			),
			wantString: "name = ? AND age > ?",
			wantValues: []any{"John", 18},
		},
		{
			name: "Multiple conditions with OR",
			conditions: NewConditions(" OR ",
				Field("status", Eq("active")),
				Field("role", Eq("admin")),
			),
			wantString: "status = ? OR role = ?",
			wantValues: []any{"active", "admin"},
		},
		{
			name: "Three conditions with AND",
			conditions: NewConditions(" AND ",
				Field("name", Eq("John")),
				Field("age", Gte(18)),
				Field("status", NotEq("deleted")),
			),
			wantString: "name = ? AND age >= ? AND status <> ?",
			wantValues: []any{"John", 18, "deleted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.conditions.String(); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := tt.conditions.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestAnd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		conditions []ConditionExpr
		wantString string
		wantValues []any
	}{
		{
			name:       "Empty And",
			conditions: []ConditionExpr{},
			wantString: "",
			wantValues: []any{},
		},
		{
			name: "And with two conditions",
			conditions: []ConditionExpr{
				Field("name", Eq("John")),
				Field("active", Eq(true)),
			},
			wantString: "name = ? AND active = ?",
			wantValues: []any{"John", true},
		},
		{
			name: "And with In conditions",
			conditions: []ConditionExpr{
				Field("id", In(1, 2, 3)),
				Field("status", In("active", "pending")),
			},
			wantString: "id IN (?,?,?) AND status IN (?,?)",
			wantValues: []any{1, 2, 3, "active", "pending"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			condition := And(tt.conditions...)
			if got := condition.String(); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestOr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		conditions []ConditionExpr
		wantString string
		wantValues []any
	}{
		{
			name:       "Empty Or",
			conditions: []ConditionExpr{},
			wantString: "",
			wantValues: []any{},
		},
		{
			name: "Or with two conditions",
			conditions: []ConditionExpr{
				Field("role", Eq("admin")),
				Field("role", Eq("moderator")),
			},
			wantString: "role = ? OR role = ?",
			wantValues: []any{"admin", "moderator"},
		},
		{
			name: "Or with different field conditions",
			conditions: []ConditionExpr{
				Field("status", Eq("active")),
				Field("premium", Eq(true)),
			},
			wantString: "status = ? OR premium = ?",
			wantValues: []any{"active", true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			condition := Or(tt.conditions...)
			if got := condition.String(); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestNestedConditions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		condition  ConditionExpr
		wantString string
		wantValues []any
	}{
		{
			name: "Nested And inside Or",
			condition: Or(
				And(
					Field("name", Eq("John")),
					Field("age", Gte(18)),
				),
				And(
					Field("role", Eq("admin")),
					Field("verified", Eq(true)),
				),
			),
			wantString: "(name = ? AND age >= ?) OR (role = ? AND verified = ?)",
			wantValues: []any{"John", 18, "admin", true},
		},
		{
			name: "Nested Or inside And",
			condition: And(
				Or(
					Field("status", Eq("active")),
					Field("status", Eq("pending")),
				),
				Or(
					Field("role", Eq("user")),
					Field("role", Eq("guest")),
				),
			),
			wantString: "(status = ? OR status = ?) AND (role = ? OR role = ?)",
			wantValues: []any{"active", "pending", "user", "guest"},
		},
		{
			name: "Nested Or and And inside And",
			condition: And(
				Or(
					Field("status", Eq("active")),
					Field("status", Eq("pending")),
				),
				And(
					Field("role", Eq("user")),
					Field("role", Eq("guest")),
				),
			),
			wantString: "(status = ? OR status = ?) AND role = ? AND role = ?",
			wantValues: []any{"active", "pending", "user", "guest"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.condition.String(); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := tt.condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}
