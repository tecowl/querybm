package expr

import (
	"reflect"
	"testing"
)

func TestCompareCondition(t *testing.T) {
	field := "field1"
	tests := []struct {
		name       string
		condition  ConditionBody
		wantString string
		wantValues []any
	}{
		{
			name:       "Eq",
			condition:  Eq("test"),
			wantString: field + " = ?",
			wantValues: []any{"test"},
		},
		{
			name:       "NotEq",
			condition:  NotEq(123),
			wantString: field + " <> ?",
			wantValues: []any{123},
		},
		{
			name:       "Gt",
			condition:  Gt(10),
			wantString: field + " > ?",
			wantValues: []any{10},
		},
		{
			name:       "Gte",
			condition:  Gte(10.5),
			wantString: field + " >= ?",
			wantValues: []any{10.5},
		},
		{
			name:       "Lt",
			condition:  Lt(100),
			wantString: field + " < ?",
			wantValues: []any{100},
		},
		{
			name:       "Lte",
			condition:  Lte(99.9),
			wantString: field + " <= ?",
			wantValues: []any{99.9},
		},
		{
			name:       "Like",
			condition:  Like("test%"),
			wantString: field + " LIKE ?",
			wantValues: []any{"test%"},
		},
		{
			name:       "LikeStartsWith",
			condition:  LikeStartsWith("prefix"),
			wantString: field + " LIKE ?",
			wantValues: []any{"prefix%"},
		},
		{
			name:       "LikeEndsWith",
			condition:  LikeEndsWith("suffix"),
			wantString: field + " LIKE ?",
			wantValues: []any{"%suffix"},
		},
		{
			name:       "LikeContains",
			condition:  LikeContains("middle"),
			wantString: field + " LIKE ?",
			wantValues: []any{"%middle%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.Build(field); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := tt.condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestInCondition(t *testing.T) {
	field := "field2"
	tests := []struct {
		name       string
		condition  ConditionBody
		wantString string
		wantValues []any
	}{
		{
			name:       "In with single value",
			condition:  In(1),
			wantString: field + " IN (?)",
			wantValues: []any{1},
		},
		{
			name:       "In with multiple values",
			condition:  In(1, 2, 3),
			wantString: field + " IN (?,?,?)",
			wantValues: []any{1, 2, 3},
		},
		{
			name:       "In with string values",
			condition:  In("a", "b", "c"),
			wantString: field + " IN (?,?,?)",
			wantValues: []any{"a", "b", "c"},
		},
		{
			name:       "In with empty values",
			condition:  In(),
			wantString: "",
			wantValues: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.Build(field); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := tt.condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestEqOrIn(t *testing.T) {
	field := "field3"
	tests := []struct {
		name       string
		values     []any
		wantString string
		wantValues []any
	}{
		{
			name:       "Single value uses Eq",
			values:     []any{"test"},
			wantString: field + " = ?",
			wantValues: []any{"test"},
		},
		{
			name:       "Multiple values uses In",
			values:     []any{1, 2, 3},
			wantString: field + " IN (?,?,?)",
			wantValues: []any{1, 2, 3},
		},
		{
			name:       "Empty values",
			values:     []any{},
			wantString: "",
			wantValues: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := EqOrIn(tt.values...)
			if got := condition.Build(field); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestStaticCondition(t *testing.T) {
	field := "field4"
	tests := []struct {
		name       string
		condition  ConditionBody
		wantString string
		wantValues []any
	}{
		{
			name:       "IsNull",
			condition:  IsNull(),
			wantString: field + " IS NULL",
			wantValues: []any{},
		},
		{
			name:       "IsNotNull",
			condition:  IsNotNull(),
			wantString: field + " IS NOT NULL",
			wantValues: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.Build(field); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := tt.condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}

func TestFieldCondition(t *testing.T) {
	tests := []struct {
		name       string
		field      ConditionExpr
		wantString string
		wantValues []any
	}{
		{
			name:       "Field with Eq condition",
			field:      Field("name", Eq("John")),
			wantString: "name = ?",
			wantValues: []any{"John"},
		},
		{
			name:       "Field with In condition",
			field:      Field("id", In(1, 2, 3)),
			wantString: "id IN (?,?,?)",
			wantValues: []any{1, 2, 3},
		},
		{
			name:       "Field with IsNull condition",
			field:      Field("deleted_at", IsNull()),
			wantString: "deleted_at IS NULL",
			wantValues: []any{},
		},
		{
			name:       "Field with nil body",
			field:      Field("test", nil),
			wantString: "",
			wantValues: []any{},
		},
		{
			name:       "Field with empty In condition",
			field:      Field("test", In()),
			wantString: "",
			wantValues: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.String(); got != tt.wantString {
				t.Errorf("String() = %q, want %q", got, tt.wantString)
			}
			if got := tt.field.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %+v, want %+v", got, tt.wantValues)
			}
		})
	}
}
