package expr

// FieldCondition represents a condition on a specific field with a condition body.
type FieldCondition struct {
	// Name is the field name to apply the condition to.
	Name string
	// Body is the condition body that defines the actual condition logic.
	Body FieldConditionBody
}

var (
	_ ConditionExpr       = (*FieldCondition)(nil)
	_ ConnectiveCondition = (*FieldCondition)(nil)
)

// Field creates a new field condition with the specified field name and condition body.
func Field(name string, body FieldConditionBody) ConditionExpr { //nolint:ireturn
	return &FieldCondition{Name: name, Body: body}
}

// String returns the SQL representation of the field condition.
func (fc *FieldCondition) String() string {
	if fc.Body == nil {
		return ""
	}
	return fc.Body.Build(fc.Name)
}

// Values returns the placeholder values for the field condition.
func (fc *FieldCondition) Values() []any {
	if fc.Body == nil {
		return []any{}
	}
	return fc.Body.Values()
}

// Connective returns the logical connective of the condition body if it implements ConnectiveCondition.
func (fc *FieldCondition) Connective() string {
	if connective, ok := fc.Body.(ConnectiveCondition); ok {
		return connective.Connective()
	}
	return ""
}

// FieldConditionBody defines the interface for field condition implementations.
type FieldConditionBody interface {
	// Build constructs the SQL condition string for the given field name.
	Build(field string) string
	// Values returns the placeholder values for the condition.
	Values() []any
}
