package querybm

import (
	"testing"

	"github.com/tecowl/querybm/statement"
)

// mockCondition is a test implementation of Condition interface
type mockCondition struct {
	buildCalled bool
	statement   *statement.Statement
}

func (m *mockCondition) Build(s *statement.Statement) {
	m.buildCalled = true
	m.statement = s
}

func TestConditionInterface(t *testing.T) {
	// Verify that mockCondition implements Condition interface
	var _ Condition = &mockCondition{}
	
	// Test that Build method is called correctly
	mock := &mockCondition{}
	stmt := statement.NewStatement("test_table", statement.NewSimpleFields("id", "name"))
	
	// Call Build method through interface
	var cond Condition = mock
	cond.Build(stmt)
	
	// Verify Build was called
	if !mock.buildCalled {
		t.Error("Build() was not called")
	}
	
	// Verify the statement was passed correctly
	if mock.statement != stmt {
		t.Error("Build() did not receive the correct statement")
	}
}

// TestMultipleConditions tests that multiple conditions can be applied
func TestMultipleConditions(t *testing.T) {
	stmt := statement.NewStatement("users", statement.NewSimpleFields("id", "name", "email"))
	
	conditions := []Condition{
		&mockCondition{},
		&mockCondition{},
		&mockCondition{},
	}
	
	// Apply all conditions
	for _, cond := range conditions {
		cond.Build(stmt)
	}
	
	// Verify all conditions were called
	for i, cond := range conditions {
		mock := cond.(*mockCondition)
		if !mock.buildCalled {
			t.Errorf("Condition %d: Build() was not called", i)
		}
		if mock.statement != stmt {
			t.Errorf("Condition %d: Build() did not receive the correct statement", i)
		}
	}
}