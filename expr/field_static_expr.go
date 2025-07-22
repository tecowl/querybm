package expr

import "fmt"

type fieldStaticExpr struct {
	value string
}

var _ FieldConditionBody = (*fieldStaticExpr)(nil)

func (c *fieldStaticExpr) Build(field string) string { return fmt.Sprintf("%s %s", field, c.value) }
func (c *fieldStaticExpr) Values() []any             { return []any{} }

func IsNull() *fieldStaticExpr    { return &fieldStaticExpr{value: "IS NULL"} }
func IsNotNull() *fieldStaticExpr { return &fieldStaticExpr{value: "IS NOT NULL"} }
