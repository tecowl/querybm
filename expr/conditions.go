package expr

import "strings"

// Conditions represents a collection of condition expressions joined by a logical connective (AND/OR).
type Conditions struct {
	items      []ConditionExpr
	connective string
}

var (
	_ ConditionExpr       = (*Conditions)(nil)
	_ ConnectiveCondition = (*Conditions)(nil)
)

// NewConditions creates a new Conditions instance with the specified connector and condition expressions.
func NewConditions(connector string, items ...ConditionExpr) *Conditions {
	return &Conditions{items: items, connective: connector}
}

// And creates a new condition expression that combines the given conditions with AND logic.
func And(conditions ...ConditionExpr) ConditionExpr { return NewConditions(" AND ", conditions...) } //nolint:ireturn

// Or creates a new condition expression that combines the given conditions with OR logic.
func Or(conditions ...ConditionExpr) ConditionExpr  { return NewConditions(" OR ", conditions...) }  //nolint:ireturn

// String returns the SQL representation of the conditions with appropriate parentheses
// when mixing different connectives.
func (c *Conditions) String() string {
	switch len(c.items) {
	case 0:
		return ""
	case 1:
		return c.items[0].String()
	}
	var sb strings.Builder
	for i, item := range c.items {
		if i > 0 {
			sb.WriteString(c.connective)
		}
		if HasDifferentConnective(item, c.connective) {
			sb.WriteString("(")
			sb.WriteString(item.String())
			sb.WriteString(")")
		} else {
			sb.WriteString(item.String())
		}
	}
	return sb.String()
}

// Values returns all placeholder values from all contained condition expressions.
func (c *Conditions) Values() []any {
	values := []any{}
	for _, item := range c.items {
		values = append(values, item.Values()...)
	}
	return values
}

// Connective returns the logical connective used to join the conditions (AND/OR).
func (c *Conditions) Connective() string {
	return c.connective
}
