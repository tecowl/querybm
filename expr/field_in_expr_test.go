package expr

import (
	"reflect"
	"testing"
)

func TestFieldInExpr(t *testing.T) {
	field := "field2"
	tests := []struct {
		name       string
		condition  FieldConditionBody
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
