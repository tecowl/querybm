package expr

type ConditionExpr interface {
	String() string
	Values() []any
}
