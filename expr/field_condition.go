package expr

import (
	"fmt"
	"strings"
)

type FieldCondition struct {
	Name string
	Body FieldConditionBody
}

var _ ConditionExpr = (*FieldCondition)(nil)

func Field(name string, body FieldConditionBody) ConditionExpr {
	return &FieldCondition{Name: name, Body: body}
}

func (fc *FieldCondition) String() string {
	if fc.Body == nil {
		return ""
	}
	return fc.Body.Build(fc.Name)
}
func (fc *FieldCondition) Values() []any {
	if fc.Body == nil {
		return []any{}
	}
	return fc.Body.Values()
}

type FieldConditionBody interface {
	Build(field string) string
	Values() []any
}

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

type inCondition struct {
	values []any
}

var _ FieldConditionBody = (*inCondition)(nil)

func (c *inCondition) Build(field string) string {
	if len(c.values) == 0 {
		return ""
	}
	return field + " IN (" + strings.Repeat("?,", len(c.values)-1) + "?)"
}

func (c *inCondition) Values() []any { return c.values }

func In(values ...any) FieldConditionBody {
	if values == nil {
		values = []any{}
	}
	return &inCondition{values: values}
}
func EqOrIn(values ...any) FieldConditionBody {
	if len(values) == 1 {
		return Eq(values[0])
	}
	return In(values...)
}

type staticCondition struct {
	value string
}

var _ FieldConditionBody = (*staticCondition)(nil)

func (c *staticCondition) Build(field string) string { return fmt.Sprintf("%s %s", field, c.value) }
func (c *staticCondition) Values() []any             { return []any{} }

func IsNull() FieldConditionBody    { return &staticCondition{value: "IS NULL"} }
func IsNotNull() FieldConditionBody { return &staticCondition{value: "IS NOT NULL"} }
