package expr

import (
	"reflect"
	"testing"
)

func TestFieldStaticExpr(t *testing.T) {
	t.Parallel()
	field := "field4"
	tests := []struct {
		name       string
		condition  FieldConditionBody
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
