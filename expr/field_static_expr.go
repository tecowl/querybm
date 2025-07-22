package expr

import "fmt"

// fieldStaticExpr represents a static SQL expression that doesn't require placeholder values.
type fieldStaticExpr struct {
	value string
}

var _ FieldConditionBody = (*fieldStaticExpr)(nil)

// Build appends the static expression to the field name.
func (c *fieldStaticExpr) Build(field string) string { return fmt.Sprintf("%s %s", field, c.value) }
// Values returns an empty slice as static expressions don't have placeholder values.
func (c *fieldStaticExpr) Values() []any             { return []any{} }

// IsNull creates a field condition for checking if a field is NULL.
func IsNull() FieldConditionBody    { return &fieldStaticExpr{value: "IS NULL"} }     // nolint:ireturn

// IsNotNull creates a field condition for checking if a field is NOT NULL.
func IsNotNull() FieldConditionBody { return &fieldStaticExpr{value: "IS NOT NULL"} } // nolint:ireturn
