package ranges

type Range[T comparable] struct {
	Start *T
	End   *T
}

func New[T comparable](start, end *T) *Range[T] {
	return &Range[T]{
		Start: start,
		End:   end,
	}
}
