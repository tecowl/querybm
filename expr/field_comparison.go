package expr

import "fmt"

type fieldComparison struct {
	operator string
	value    any
}

var _ FieldConditionBody = (*fieldComparison)(nil)

func newCompare(operator string, value any) *fieldComparison {
	return &fieldComparison{operator: operator, value: value}
}
func (c *fieldComparison) Build(field string) string {
	return fmt.Sprintf("%s %s ?", field, c.operator)
}
func (c *fieldComparison) Values() []any { return []any{c.value} }

func Eq(value any) FieldConditionBody    { return newCompare("=", value) }
func NotEq(value any) FieldConditionBody { return newCompare("<>", value) }
func Gt(value any) FieldConditionBody    { return newCompare(">", value) }
func Gte(value any) FieldConditionBody   { return newCompare(">=", value) }
func Lt(value any) FieldConditionBody    { return newCompare("<", value) }
func Lte(value any) FieldConditionBody   { return newCompare("<=", value) }

func Like(value string) FieldConditionBody           { return newCompare("LIKE", value) }
func LikeStartsWith(value string) FieldConditionBody { return Like(value + "%") }
func LikeEndsWith(value string) FieldConditionBody   { return Like("%" + value) }
func LikeContains(value string) FieldConditionBody   { return Like("%" + value + "%") }
