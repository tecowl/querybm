package expr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFieldInRangeExpr(t *testing.T) {
	field := "field1"
	tests := []struct {
		name       string
		condition  FieldConditionBody
		wantString string
		wantValues []any
	}{
		{
			name:       "InRange with start and end",
			condition:  InRange(10, 20),
			wantString: fmt.Sprintf("%s >= ? AND %s < ?", field, field),
			wantValues: []any{10, 20},
		},
		{
			name:       "InRange with only start",
			condition:  InRange(10, nil),
			wantString: field + " >= ?",
			wantValues: []any{10},
		},
		{
			name:       "InRange with only end",
			condition:  InRange(nil, 20),
			wantString: field + " < ?",
			wantValues: []any{20},
		},
		{
			name:       "InRange with no values",
			condition:  InRange(nil, nil),
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
				// reflect.DeepEqual returns false for empty slices, so we handle that case
				if len(got) == 0 && len(tt.wantValues) == 0 {
					return
				}
				t.Errorf("Values() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}
