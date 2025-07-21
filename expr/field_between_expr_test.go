package expr

import (
	"reflect"
	"testing"
)

func TestFieldBetweenExpr(t *testing.T) {
	field := "field1"
	tests := []struct {
		name       string
		condition  FieldConditionBody
		wantString string
		wantValues []any
	}{
		{
			name:       "Between with start and end",
			condition:  Between(10, 20),
			wantString: field + " BETWEEN ? AND ?",
			wantValues: []any{10, 20},
		},
		{
			name:       "Between with only start",
			condition:  Between(10, nil),
			wantString: field + " >= ?",
			wantValues: []any{10},
		},
		{
			name:       "Between with only end",
			condition:  Between(nil, 20),
			wantString: field + " <= ?",
			wantValues: []any{20},
		},
		{
			name:       "Between with no values",
			condition:  Between(nil, nil),
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
