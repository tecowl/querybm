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

func Eq(value any) FieldConditionBody    { return newCompare("=", value) }  //nolint:ireturn
func NotEq(value any) FieldConditionBody { return newCompare("<>", value) } //nolint:ireturn
func Gt(value any) FieldConditionBody    { return newCompare(">", value) }  //nolint:ireturn
func Gte(value any) FieldConditionBody   { return newCompare(">=", value) } //nolint:ireturn
func Lt(value any) FieldConditionBody    { return newCompare("<", value) }  //nolint:ireturn
func Lte(value any) FieldConditionBody   { return newCompare("<=", value) } //nolint:ireturn

func Like(value string) FieldConditionBody           { return newCompare("LIKE", value) } //nolint:ireturn
func LikeStartsWith(value string) FieldConditionBody { return Like(value + "%") }         //nolint:ireturn
func LikeEndsWith(value string) FieldConditionBody   { return Like("%" + value) }         //nolint:ireturn
func LikeContains(value string) FieldConditionBody   { return Like("%" + value + "%") }   //nolint:ireturn
