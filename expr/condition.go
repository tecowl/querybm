package expr

import (
	"fmt"
	"strings"
)

type ConditionExpr interface {
	String() string
	Values() []any
}

type ConditionBody interface {
	String() string
	Values() []any
}

type compareCondition struct {
	operator string
	value    any
}

var _ ConditionBody = (*compareCondition)(nil)

func newCompare(operator string, value any) *compareCondition {
	return &compareCondition{operator: operator, value: value}
}
func (c *compareCondition) String() string { return fmt.Sprintf("%s ?", c.operator) }
func (c *compareCondition) Values() []any  { return []any{c.value} }

func Eq(value any) ConditionBody    { return newCompare("=", value) }
func NotEq(value any) ConditionBody { return newCompare("<>", value) }
func Gt(value any) ConditionBody    { return newCompare(">", value) }
func Gte(value any) ConditionBody   { return newCompare(">=", value) }
func Lt(value any) ConditionBody    { return newCompare("<", value) }
func Lte(value any) ConditionBody   { return newCompare("<=", value) }

func Like(value string) ConditionBody           { return newCompare("LIKE", value) }
func LikeStartsWith(value string) ConditionBody { return Like(value + "%") }
func LikeEndsWith(value string) ConditionBody   { return Like("%" + value) }
func LikeContains(value string) ConditionBody   { return Like("%" + value + "%") }

type inCondition struct {
	values []any
}

var _ ConditionBody = (*inCondition)(nil)

func (c *inCondition) String() string {
	if len(c.values) == 0 {
		return ""
	}
	return "IN (" + strings.Repeat("?,", len(c.values)-1) + "?)"
}

func (c *inCondition) Values() []any { return c.values }

func In(values ...any) ConditionBody {
	if values == nil {
		values = []any{}
	}
	return &inCondition{values: values}
}
func EqOrIn(values ...any) ConditionBody {
	if len(values) == 1 {
		return Eq(values[0])
	}
	return In(values...)
}

type staticCondition struct {
	value string
}

var _ ConditionBody = (*staticCondition)(nil)

func (c *staticCondition) String() string { return c.value }
func (c *staticCondition) Values() []any  { return []any{} }

func IsNull() ConditionBody    { return &staticCondition{value: "IS NULL"} }
func IsNotNull() ConditionBody { return &staticCondition{value: "IS NOT NULL"} }

type FieldCondition struct {
	Name string
	Body ConditionBody
}

var _ ConditionExpr = (*FieldCondition)(nil)

func Field(name string, body ConditionBody) *FieldCondition {
	return &FieldCondition{Name: name, Body: body}
}

func (fc *FieldCondition) String() string {
	if fc.Body == nil {
		return ""
	}
	body := fc.Body.String()
	if body == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", fc.Name, body)
}
func (fc *FieldCondition) Values() []any {
	if fc.Body == nil {
		return []any{}
	}
	return fc.Body.Values()
}
