package expr

import "strings"

// fieldInExpr represents an IN condition for checking if a field value is in a list of values.
type fieldInExpr struct {
	values []any
}

var _ FieldConditionBody = (*fieldInExpr)(nil)

// Build constructs the IN SQL clause for the given field.
// Returns an empty string if no values are provided.
func (c *fieldInExpr) Build(field string) string {
	if len(c.values) == 0 {
		return ""
	}
	return field + " IN (" + strings.Repeat("?,", len(c.values)-1) + "?)"
}

// Values returns all values for the IN clause.
func (c *fieldInExpr) Values() []any { return c.values }

// In creates a field condition for IN comparison.
// It checks if the field value is in the provided list of values.
func In(values ...any) FieldConditionBody { //nolint:ireturn
	if values == nil {
		values = []any{}
	}
	return &fieldInExpr{values: values}
}

// EqOrIn creates either an equality condition (if one value) or an IN condition (if multiple values).
// This is useful for optimizing queries when the number of values is dynamic.
func EqOrIn(values ...any) FieldConditionBody { //nolint:ireturn
	if len(values) == 1 {
		return Eq(values[0])
	}
	return In(values...)
}
