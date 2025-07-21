package expr

import "strings"

type fieldInExpr struct {
	values []any
}

var _ FieldConditionBody = (*fieldInExpr)(nil)

func (c *fieldInExpr) Build(field string) string {
	if len(c.values) == 0 {
		return ""
	}
	return field + " IN (" + strings.Repeat("?,", len(c.values)-1) + "?)"
}

func (c *fieldInExpr) Values() []any { return c.values }

func In(values ...any) FieldConditionBody {
	if values == nil {
		values = []any{}
	}
	return &fieldInExpr{values: values}
}
func EqOrIn(values ...any) FieldConditionBody {
	if len(values) == 1 {
		return Eq(values[0])
	}
	return In(values...)
}
