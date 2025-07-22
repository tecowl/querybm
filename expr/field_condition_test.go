package expr

import (
	"reflect"
	"testing"
)

func TestFieldCondition(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := tt.field.String(); got != tt.wantString {
				t.Errorf("String() = %q, want %q", got, tt.wantString)
			}
			if got := tt.field.Values(); !reflect.DeepEqual(got, tt.wantValues) {
				t.Errorf("Values() = %+v, want %+v", got, tt.wantValues)
			}
		})
	}
}
