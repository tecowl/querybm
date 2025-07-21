package expr

type fieldBetweenExpr struct {
	start any
	end   any
}

var _ FieldConditionBody = (*fieldBetweenExpr)(nil)

func (c *fieldBetweenExpr) Build(field string) string {
	if c.start == nil && c.end == nil {
		return ""
	}
	if c.start != nil && c.end != nil {
		return field + " BETWEEN ? AND ?"
	}
	if c.start != nil {
		return field + " >= ?"
	}
	return field + " <= ?"
}
func (c *fieldBetweenExpr) Values() []any {
	if c.start == nil && c.end == nil {
		return []any{}
	}
	if c.start != nil && c.end != nil {
		return []any{c.start, c.end}
	}
	if c.start != nil {
		return []any{c.start}
	}
	return []any{c.end}
}

func Between(start, end any) FieldConditionBody {
	return &fieldBetweenExpr{start: start, end: end}
}
