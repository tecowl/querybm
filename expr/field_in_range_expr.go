package expr

// inRangeExpr represents a range condition checking if a field value is within [start, end).
type inRangeExpr struct {
	start any
	end   any
}

var (
	_ FieldConditionBody  = (*inRangeExpr)(nil)
	_ ConnectiveCondition = (*inRangeExpr)(nil)
)

// Build constructs the range condition as field >= start AND field < end.
func (c *inRangeExpr) Build(field string) string {
	r := And(Field(field, Gte(c.start)), Field(field, Lt(c.end)))
	return r.String()
}

// Values returns the start and end values for the range condition.
func (c *inRangeExpr) Values() []any {
	return []any{c.start, c.end}
}

// Connective returns " AND " as the range condition uses AND logic internally.
func (c *inRangeExpr) Connective() string {
	return " AND "
}

// InRange creates a field condition for checking if a value is in the range [start, end).
// The condition is equivalent to: field >= start AND field < end.
func InRange(start, end any) FieldConditionBody { //nolint:ireturn
	return &inRangeExpr{start: start, end: end}
}
