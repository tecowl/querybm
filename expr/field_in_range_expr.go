package expr

type inRangeExpr struct {
	start any
	end   any
}

var _ FieldConditionBody = (*inRangeExpr)(nil)
var _ ConnectiveCondition = (*inRangeExpr)(nil)

func (c *inRangeExpr) Build(field string) string {
	r := And(Field(field, Gte(c.start)), Field(field, Lt(c.end)))
	return r.String()
}

func (c *inRangeExpr) Values() []any {
	return []any{c.start, c.end}
}

func (c *inRangeExpr) Connective() string {
	return " AND "
}

func InRange(start, end any) FieldConditionBody { //nolint:ireturn
	return &inRangeExpr{start: start, end: end}
}
