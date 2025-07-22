package statement

import (
	"github.com/tecowl/querybm/expr"
)

// WhereBlock represents the WHERE clause of a SQL statement.
type WhereBlock struct {
	// Connector is the logical connector (AND/OR) used between conditions.
	Connector  string
	conditions []expr.ConditionExpr
}

// newWhere creates a new WhereBlock with the specified connector.
func newWhere(connector string) *WhereBlock {
	return &WhereBlock{Connector: connector, conditions: []expr.ConditionExpr{}}
}

// Add appends a condition expression to the WHERE clause.
func (b *WhereBlock) Add(condition expr.ConditionExpr) {
	b.conditions = append(b.conditions, condition)
}

// IsEmpty returns true if there are no conditions in the WHERE clause.
func (b *WhereBlock) IsEmpty() bool {
	return len(b.conditions) == 0
}

// Build constructs the WHERE clause string and returns it with placeholder values.
func (b *WhereBlock) Build() (string, []any) {
	conditions := expr.NewConditions(b.Connector, b.conditions...)
	return conditions.String(), conditions.Values()
}
