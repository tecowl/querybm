package querybm

import (
	"reflect"
	"testing"

	"github.com/tecowl/querybm/statement"
)

func TestNoop(t *testing.T) {
	t.Parallel()
	fields := statement.NewSimpleFields("id", "name")
	st := statement.New("users", fields)
	origSelect, origVars := st.Build()

	Noop.Build(st)

	afterSelect, afterVars := st.Build()
	if origSelect != afterSelect {
		t.Errorf("Expected select statement to remain unchanged, got %s", afterSelect)
	}
	if !reflect.DeepEqual(origVars, afterVars) {
		t.Errorf("Expected variables to remain unchanged, got %v", afterVars)
	}
}
