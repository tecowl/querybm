package helpers

type Range[T comparable] struct {
	Start *T
	End   *T
}

func NewRange[T comparable](start, end *T) Range[T] {
	return Range[T]{
		Start: start,
		End:   end,
	}
}
