// Package expr provides expression types and builders for SQL WHERE clause conditions.
package expr

// ConditionExpr represents a SQL condition expression that can be converted to a string
// with placeholder values.
type ConditionExpr interface {
	// String returns the SQL condition string with placeholders.
	String() string
	// Values returns the values to be used with the placeholders in the SQL string.
	Values() []any
}
