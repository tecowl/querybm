package statement

// Fields is an interface for types that can provide a list of field names for SELECT clauses.
type Fields interface {
	// Fields returns the list of field names to be selected.
	Fields() []string
}

// SimpleFields is a basic implementation of Fields using a string slice.
type SimpleFields []string

var _ Fields = SimpleFields{}

// NewSimpleFields creates a new SimpleFields instance with the provided field names.
func NewSimpleFields(fields ...string) SimpleFields {
	return fields
}

// Fields returns the field names as a string slice.
func (f SimpleFields) Fields() []string {
	return f
}
