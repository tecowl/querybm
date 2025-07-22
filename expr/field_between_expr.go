package expr

type fieldBetweenExpr struct {
	start any
	end   any
}

var _ FieldConditionBody = (*fieldBetweenExpr)(nil)

func (c *fieldBetweenExpr) Build(field string) string {
	return field + " BETWEEN ? AND ?"
}

func (c *fieldBetweenExpr) Values() []any {
	return []any{c.start, c.end}
}

func Between(start, end any) FieldConditionBody { //nolint:ireturn
	return &fieldBetweenExpr{start: start, end: end}
}
