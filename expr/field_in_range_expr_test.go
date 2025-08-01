package expr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFieldInRangeExpr(t *testing.T) {
	t.Parallel()
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

func TestFieldInRangeExprWithConnectives(t *testing.T) {
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
				Field("amount", InRange(10, 20)),
				And(
					Field("role", Eq("admin")),
					Field("verified", Eq(true)),
				),
			),
			wantString: "(amount >= ? AND amount < ?) OR (role = ? AND verified = ?)",
			wantValues: []any{10, 20, "admin", true},
		},
		{
			name: "Nested Or inside And",
			condition: And(
				Field("available", InRange("2023-01-01", "2023-12-31")),
				Or(
					Field("role", Eq("user")),
					Field("role", Eq("guest")),
				),
			),
			wantString: "available >= ? AND available < ? AND (role = ? OR role = ?)",
			wantValues: []any{"2023-01-01", "2023-12-31", "user", "guest"},
		},
		{
			name: "Nested Or and And inside And",
			condition: And(
				Field("available", InRange("2023-01-01", "2023-12-31")),
				And(
					Field("role", Eq("user")),
					Field("role", Eq("guest")),
				),
			),
			wantString: "available >= ? AND available < ? AND role = ? AND role = ?",
			wantValues: []any{"2023-01-01", "2023-12-31", "user", "guest"},
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
