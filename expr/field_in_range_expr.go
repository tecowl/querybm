package expr

type inRangeExpr struct {
	start any
	end   any
}

var _ FieldConditionBody = (*inRangeExpr)(nil)

func (c *inRangeExpr) Build(field string) string {
	if c.start == nil && c.end == nil {
		return ""
	}
	if c.start == nil {
		return Lt(c.end).Build(field) // Not using Lte here because we want to exclude the end value
	}
	if c.end == nil {
		return Gte(c.start).Build(field)
	}
	r := And(Field(field, Gte(c.start)), Field(field, Lt(c.end)))
	return r.String()
}

func (c *inRangeExpr) Values() []any {
	var values []any
	if c.start != nil {
		values = append(values, c.start)
	}
	if c.end != nil {
		values = append(values, c.end)
	}
	return values
}

func InRange(start, end any) FieldConditionBody {
	return &inRangeExpr{start: start, end: end}
}
