package expr

// fieldBetweenExpr represents a BETWEEN condition for a field.
type fieldBetweenExpr struct {
	start any
	end   any
}

var _ FieldConditionBody = (*fieldBetweenExpr)(nil)

// Build constructs the BETWEEN SQL clause for the given field.
func (c *fieldBetweenExpr) Build(field string) string {
	return field + " BETWEEN ? AND ?"
}

// Values returns the start and end values for the BETWEEN clause.
func (c *fieldBetweenExpr) Values() []any {
	return []any{c.start, c.end}
}

// Between creates a field condition body for a BETWEEN expression.
// It checks if the field value is between start and end (inclusive).
func Between(start, end any) FieldConditionBody { //nolint:ireturn
	return &fieldBetweenExpr{start: start, end: end}
}
