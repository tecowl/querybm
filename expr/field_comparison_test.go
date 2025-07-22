package expr

import (
	"reflect"
	"testing"
)

func TestFieldComparison(t *testing.T) {
	t.Parallel()
	field := "field1"
	tests := []struct {
		name       string
		condition  FieldConditionBody
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
			t.Parallel()
			if got := tt.condition.Build(field); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
			if got := tt.condition.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}
