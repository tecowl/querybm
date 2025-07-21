package expr

import (
	"fmt"
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

type fieldStaticExpr struct {
	value string
}

var _ FieldConditionBody = (*fieldStaticExpr)(nil)

func (c *fieldStaticExpr) Build(field string) string { return fmt.Sprintf("%s %s", field, c.value) }
func (c *fieldStaticExpr) Values() []any             { return []any{} }

func IsNull() FieldConditionBody    { return &fieldStaticExpr{value: "IS NULL"} }
func IsNotNull() FieldConditionBody { return &fieldStaticExpr{value: "IS NOT NULL"} }
