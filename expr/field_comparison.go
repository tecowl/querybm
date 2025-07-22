package expr

import "fmt"

// fieldComparison represents a comparison operation between a field and a value.
type fieldComparison struct {
	operator string
	value    any
}

var _ FieldConditionBody = (*fieldComparison)(nil)

// newCompare creates a new field comparison with the specified operator and value.
func newCompare(operator string, value any) *fieldComparison {
	return &fieldComparison{operator: operator, value: value}
}

// Build constructs the comparison SQL clause for the given field.
func (c *fieldComparison) Build(field string) string {
	return fmt.Sprintf("%s %s ?", field, c.operator)
}

// Values returns the comparison value as a slice.
func (c *fieldComparison) Values() []any { return []any{c.value} }

// Eq creates a field condition for equality comparison (=).
func Eq(value any) FieldConditionBody { return newCompare("=", value) } //nolint:ireturn

// NotEq creates a field condition for inequality comparison (<>).
func NotEq(value any) FieldConditionBody { return newCompare("<>", value) } //nolint:ireturn

// Gt creates a field condition for greater than comparison (>).
func Gt(value any) FieldConditionBody { return newCompare(">", value) } //nolint:ireturn

// Gte creates a field condition for greater than or equal comparison (>=).
func Gte(value any) FieldConditionBody { return newCompare(">=", value) } //nolint:ireturn

// Lt creates a field condition for less than comparison (<).
func Lt(value any) FieldConditionBody { return newCompare("<", value) } //nolint:ireturn

// Lte creates a field condition for less than or equal comparison (<=).
func Lte(value any) FieldConditionBody { return newCompare("<=", value) } //nolint:ireturn

// Like creates a field condition for LIKE comparison.
func Like(value string) FieldConditionBody { return newCompare("LIKE", value) } //nolint:ireturn

// LikeStartsWith creates a field condition for prefix matching (value%).
func LikeStartsWith(value string) FieldConditionBody { return Like(value + "%") } //nolint:ireturn

// LikeEndsWith creates a field condition for suffix matching (%value).
func LikeEndsWith(value string) FieldConditionBody { return Like("%" + value) } //nolint:ireturn

// LikeContains creates a field condition for substring matching (%value%).
func LikeContains(value string) FieldConditionBody { return Like("%" + value + "%") } //nolint:ireturn
