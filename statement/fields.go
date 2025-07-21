package statement

type Fields interface {
	Fields() []string
}

type SimpleFields []string

var _ Fields = SimpleFields{}

func NewSimpleFields(fields ...string) SimpleFields {
	return fields
}

func (f SimpleFields) Fields() []string {
	return f
}
