package expr

type FieldCondition struct {
	Name string
	Body FieldConditionBody
}

var _ ConditionExpr = (*FieldCondition)(nil)
var _ ConnectiveCondition = (*FieldCondition)(nil)

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

func (fc *FieldCondition) Connective() string {
	if connective, ok := fc.Body.(ConnectiveCondition); ok {
		return connective.Connective()
	}
	return ""
}

type FieldConditionBody interface {
	Build(field string) string
	Values() []any
}
