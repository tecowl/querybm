package querybm

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/statement"
)

// Define static errors for testing.
var (
	errConditionError = errors.New("condition error")
	errSortError      = errors.New("sort error")
)

type TestModel struct {
	ID   int
	Name string
}

type TestCondition struct {
	whereAdded bool
}

func (t *TestCondition) Build(s *statement.Statement) {
	t.whereAdded = true
	s.Where.Add(expr.Field("status", expr.Eq("active")))
}

type TestSort struct {
	sortAdded bool
}

func (t *TestSort) Build(s *statement.Statement) {
	t.sortAdded = true
	s.Sort.Add("created_at DESC")
}

type ValidatableCondition struct {
	TestCondition
	validateCalled bool
	validateErr    error
}

func (v *ValidatableCondition) Validate() error {
	v.validateCalled = true
	return v.validateErr
}

type ValidatableSort struct {
	TestSort
	validateCalled bool
	validateErr    error
}

func (v *ValidatableSort) Validate() error {
	v.validateCalled = true
	return v.validateErr
}

func TestNew(t *testing.T) {
	t.Parallel()
	db := &sql.DB{}
	condition := &TestCondition{}
	sort := &TestSort{}
	fields := NewFields([]string{"id", "name"}, func(s Scanner, m *TestModel) error {
		return s.Scan(&m.ID, &m.Name)
	})
	limitOffset := NewLimitOffset(10, 0)

	query := New(db, "users", fields, condition, sort, limitOffset)

	if query.Table != "users" {
		t.Errorf("New() Table = %v, want %v", query.Table, "users")
	}
	if query.Condition != condition {
		t.Error("New() Condition not set correctly")
	}
	if query.Sort != sort {
		t.Error("New() Sort not set correctly")
	}
	if query.LimitOffset != limitOffset {
		t.Error("New() LimitOffset not set correctly")
	}
}

func TestQuery_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupQuery    func() *Query[TestModel]
		wantErr       bool
		wantErrString string
	}{
		{
			name: "All validations pass",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &ValidatableCondition{}
				sort := &ValidatableSort{}
				fields := NewFields[TestModel]([]string{"id"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantErr: false,
		},
		{
			name: "Condition validation fails",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &ValidatableCondition{validateErr: errors.New("invalid condition")} // nolint:err113
				sort := &ValidatableSort{}
				fields := NewFields[TestModel]([]string{"id"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantErr:       true,
			wantErrString: "condition validation failed:",
		},
		{
			name: "Sort validation fails",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &ValidatableCondition{}
				sort := &ValidatableSort{validateErr: errors.New("invalid sort")} // nolint:err113
				fields := NewFields[TestModel]([]string{"id"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantErr:       true,
			wantErrString: "sort validation failed:",
		},
		{
			name: "Non-validatable condition and sort",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &TestCondition{}
				sort := &TestSort{}
				fields := NewFields[TestModel]([]string{"id"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := tt.setupQuery()
			err := q.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.wantErrString != "" && err != nil {
				if !contains(err.Error(), tt.wantErrString) {
					t.Errorf("Validate() error = %v, want to contain %v", err, tt.wantErrString)
				}
			}
		})
	}
}

func TestQuery_BuildCountSelect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setupQuery func() *Query[TestModel]
		wantSQL    string
		wantValues []any
	}{
		{
			name: "Count without conditions",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &TestCondition{}
				sort := &TestSort{}
				fields := NewFields[TestModel]([]string{"id", "name"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantSQL:    "SELECT COUNT(*) AS count FROM users WHERE status = ?",
			wantValues: []any{"active"},
		},
		{
			name: "Count with table name",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &TestCondition{}
				sort := &TestSort{}
				fields := NewFields[TestModel]([]string{"id", "name"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "products", fields, condition, sort, limitOffset)
			},
			wantSQL:    "SELECT COUNT(*) AS count FROM products WHERE status = ?",
			wantValues: []any{"active"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := tt.setupQuery()
			gotSQL, gotValues := q.BuildCountSelect()
			if gotSQL != tt.wantSQL {
				t.Errorf("BuildCountSelect() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("BuildCountSelect() values = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestQuery_BuildRowsSelect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setupQuery func() *Query[TestModel]
		wantSQL    string
		wantValues []any
	}{
		{
			name: "Rows with all components",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &TestCondition{}
				sort := &TestSort{}
				fields := NewFields[TestModel]([]string{"id", "name", "email"}, nil)
				limitOffset := NewLimitOffset(20, 40)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantSQL:    "SELECT id, name, email FROM users WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
			wantValues: []any{"active", int64(20), int64(40)},
		},
		{
			name: "Rows without limitOffset offset",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &TestCondition{}
				sort := &TestSort{}
				fields := NewFields[TestModel]([]string{"id", "name"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantSQL:    "SELECT id, name FROM users WHERE status = ? ORDER BY created_at DESC LIMIT ?",
			wantValues: []any{"active", int64(10)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := tt.setupQuery()
			gotSQL, gotValues := q.BuildRowsSelect()
			if gotSQL != tt.wantSQL {
				t.Errorf("BuildRowsSelect() SQL = %v, want %v", gotSQL, tt.wantSQL)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("BuildRowsSelect() values = %v, want %v", gotValues, tt.wantValues)
			}
			// Verify that condition and sort were applied
			if tc, ok := q.Condition.(*TestCondition); ok {
				if !tc.whereAdded {
					t.Error("BuildRowsSelect() did not apply condition")
				}
			} else {
				t.Error("Condition is not of type TestCondition, cannot verify whereAdded")
			}
			if ts, ok := q.Sort.(*TestSort); ok {
				if !ts.sortAdded {
					t.Error("BuildRowsSelect() did not apply sort")
				}
			} else {
				t.Error("Sort is not of type TestSort, cannot verify sortAdded")
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && findSubstring(s, substr) != -1))
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestQuery_Validate_WithValidatableComponents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		setupQuery  func() *Query[TestModel]
		wantError   bool
		errorString string
	}{
		{
			name: "condition validation fails",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &ValidatableCondition{validateErr: errConditionError}
				sort := &TestSort{}
				fields := NewFields[TestModel]([]string{"id", "name"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantError:   true,
			errorString: "condition validation failed: condition error",
		},
		{
			name: "sort validation fails",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &TestCondition{}
				sort := &ValidatableSort{validateErr: errSortError}
				fields := NewFields[TestModel]([]string{"id", "name"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantError:   true,
			errorString: "sort validation failed: sort error",
		},
		{
			name: "all validations pass",
			setupQuery: func() *Query[TestModel] {
				db := &sql.DB{}
				condition := &ValidatableCondition{validateErr: nil}
				sort := &ValidatableSort{validateErr: nil}
				fields := NewFields[TestModel]([]string{"id", "name"}, nil)
				limitOffset := NewLimitOffset(10, 0)
				return New(db, "users", fields, condition, sort, limitOffset)
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := tt.setupQuery()
			err := q.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
			if err != nil && tt.errorString != "" && !contains(err.Error(), tt.errorString) {
				t.Errorf("Validate() error = %v, want error containing %v", err, tt.errorString)
			}
		})
	}
}
