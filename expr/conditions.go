package expr

import "strings"

type Conditions struct {
	items     []ConditionExpr
	connector string
}

var _ ConditionExpr = (*Conditions)(nil)

func NewConditions(connector string, items ...ConditionExpr) *Conditions {
	return &Conditions{items: items, connector: connector}
}
func And(conditions ...ConditionExpr) ConditionExpr { return NewConditions(" AND ", conditions...) }
func Or(conditions ...ConditionExpr) ConditionExpr  { return NewConditions(" OR ", conditions...) }

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
			sb.WriteString(c.connector)
		}
		subconditions, ok := item.(*Conditions)
		if ok && subconditions.connector != c.connector {
			sb.WriteString("(")
			sb.WriteString(subconditions.String())
			sb.WriteString(")")
		} else {
			sb.WriteString(item.String())
		}
	}
	return sb.String()
}
func (c *Conditions) Values() []any {
	values := []any{}
	for _, item := range c.items {
		values = append(values, item.Values()...)
	}
	return values
}
