package expr

import "strings"

type Conditions struct {
	items      []ConditionExpr
	connective string
}

var _ ConditionExpr = (*Conditions)(nil)
var _ ConnectiveCondition = (*Conditions)(nil)

func NewConditions(connector string, items ...ConditionExpr) *Conditions {
	return &Conditions{items: items, connective: connector}
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
func (c *Conditions) Values() []any {
	values := []any{}
	for _, item := range c.items {
		values = append(values, item.Values()...)
	}
	return values
}

func (c *Conditions) Connective() string {
	return c.connective
}
