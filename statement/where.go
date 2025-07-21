package statement

import (
	"github.com/tecowl/querybm/expr"
)

type WhereBlock struct {
	Connector  string
	conditions []expr.ConditionExpr
}

func newWhere(connector string) *WhereBlock {
	return &WhereBlock{Connector: connector, conditions: []expr.ConditionExpr{}}
}

func (b *WhereBlock) Add(condition expr.ConditionExpr) {
	b.conditions = append(b.conditions, condition)
}

func (b *WhereBlock) IsEmpty() bool {
	return len(b.conditions) == 0
}

func (b *WhereBlock) Build() (string, []any) {
	conditions := expr.NewConditions(b.Connector, b.conditions...)
	return conditions.String(), conditions.Values()
}
