package ranges

import (
	"time"

	"github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/statement"
)

type Range[T comparable] struct {
	Start *T
	End   *T

	useBetween bool // Indicates if the range should be used with BETWEEN SQL syntax
}

func New[T comparable](start, end *T) *Range[T] {
	return &Range[T]{
		Start:      start,
		End:        end,
		useBetween: false,
	}
}

func (r *Range[T]) UseBetween() *Range[T] {
	r.useBetween = true
	return r
}

func NewTimeRange(start, end time.Time) *Range[time.Time] {
	var vStart *time.Time
	if !start.IsZero() {
		vStart = &start
	}
	var vEnd *time.Time
	if !end.IsZero() {
		vEnd = &end
	}
	return &Range[time.Time]{Start: vStart, End: vEnd}
}

func NewIntRange(start, end int) *Range[int] {
	var vStart *int
	if start != 0 {
		vStart = &start
	}
	var vEnd *int
	if end != 0 {
		vEnd = &end
	}
	return &Range[int]{Start: vStart, End: vEnd}
}
func NewInt32Range(start, end int32) *Range[int32] {
	var vStart *int32
	if start != 0 {
		vStart = &start
	}
	var vEnd *int32
	if end != 0 {
		vEnd = &end
	}
	return &Range[int32]{Start: vStart, End: vEnd}
}
func NewInt64Range(start, end int64) *Range[int64] {
	var vStart *int64
	if start != 0 {
		vStart = &start
	}
	var vEnd *int64
	if end != 0 {
		vEnd = &end
	}
	return &Range[int64]{Start: vStart, End: vEnd}
}

func NewUintRange(start, end uint) *Range[uint] {
	var vStart *uint
	if start != 0 {
		vStart = &start
	}
	var vEnd *uint
	if end != 0 {
		vEnd = &end
	}
	return &Range[uint]{Start: vStart, End: vEnd}
}
func NewUint32Range(start, end uint32) *Range[uint32] {
	var vStart *uint32
	if start != 0 {
		vStart = &start
	}
	var vEnd *uint32
	if end != 0 {
		vEnd = &end
	}
	return &Range[uint32]{Start: vStart, End: vEnd}
}
func NewUint64Range(start, end uint64) *Range[uint64] {
	var vStart *uint64
	if start != 0 {
		vStart = &start
	}
	var vEnd *uint64
	if end != 0 {
		vEnd = &end
	}
	return &Range[uint64]{Start: vStart, End: vEnd}
}

func (r *Range[T]) Build(field string, st *statement.Statement) {
	if r.useBetween {
		st.Where.Add(expr.Field(field, expr.Between(r.Start, r.End)))
	} else {
		st.Where.Add(expr.Field(field, expr.InRange(r.Start, r.End)))
	}
}
