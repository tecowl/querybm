// Package ranges provides generic range types and utilities for building range-based SQL conditions.
package ranges

import (
	"time"

	"github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/statement"
)

// Range represents a range of values with optional start and end boundaries.
type Range[T comparable] struct {
	// Start is the beginning of the range (inclusive).
	Start *T
	// End is the end of the range (exclusive by default, inclusive with UseBetween).
	End   *T

	useBetween bool // Indicates if the range should be used with BETWEEN SQL syntax
}

// New creates a new Range with the specified start and end values.
// Either start or end can be nil to create an open-ended range.
func New[T comparable](start, end *T) *Range[T] {
	return &Range[T]{
		Start:      start,
		End:        end,
		useBetween: false,
	}
}

// UseBetween configures the range to use BETWEEN syntax (inclusive end) instead of >= and < syntax.
func (r *Range[T]) UseBetween() *Range[T] {
	r.useBetween = true
	return r
}

// NewTimeRange creates a Range for time.Time values.
// Zero time values are treated as nil boundaries.
func NewTimeRange(start, end time.Time) *Range[time.Time] {
	var vStart *time.Time
	if !start.IsZero() {
		vStart = &start
	}
	var vEnd *time.Time
	if !end.IsZero() {
		vEnd = &end
	}
	return New[time.Time](vStart, vEnd)
}

// NewIntRange creates a Range for int values.
// Zero values are treated as nil boundaries.
func NewIntRange(start, end int) *Range[int] {
	var vStart *int
	if start != 0 {
		vStart = &start
	}
	var vEnd *int
	if end != 0 {
		vEnd = &end
	}
	return New[int](vStart, vEnd)
}

// NewInt32Range creates a Range for int32 values.
// Zero values are treated as nil boundaries.
func NewInt32Range(start, end int32) *Range[int32] {
	var vStart *int32
	if start != 0 {
		vStart = &start
	}
	var vEnd *int32
	if end != 0 {
		vEnd = &end
	}
	return New[int32](vStart, vEnd)
}

// NewInt64Range creates a Range for int64 values.
// Zero values are treated as nil boundaries.
func NewInt64Range(start, end int64) *Range[int64] {
	var vStart *int64
	if start != 0 {
		vStart = &start
	}
	var vEnd *int64
	if end != 0 {
		vEnd = &end
	}
	return New[int64](vStart, vEnd)
}

// NewUintRange creates a Range for uint values.
// Zero values are treated as nil boundaries.
func NewUintRange(start, end uint) *Range[uint] {
	var vStart *uint
	if start != 0 {
		vStart = &start
	}
	var vEnd *uint
	if end != 0 {
		vEnd = &end
	}
	return New[uint](vStart, vEnd)
}

// NewUint32Range creates a Range for uint32 values.
// Zero values are treated as nil boundaries.
func NewUint32Range(start, end uint32) *Range[uint32] {
	var vStart *uint32
	if start != 0 {
		vStart = &start
	}
	var vEnd *uint32
	if end != 0 {
		vEnd = &end
	}
	return New[uint32](vStart, vEnd)
}

// NewUint64Range creates a Range for uint64 values.
// Zero values are treated as nil boundaries.
func NewUint64Range(start, end uint64) *Range[uint64] {
	var vStart *uint64
	if start != 0 {
		vStart = &start
	}
	var vEnd *uint64
	if end != 0 {
		vEnd = &end
	}
	return New[uint64](vStart, vEnd)
}

// Build adds the appropriate range conditions to the statement's WHERE clause.
// It handles different cases: both boundaries, only start, only end, or no boundaries.
func (r *Range[T]) Build(field string, st *statement.Statement) {
	if r.Start == nil && r.End == nil {
		return
	}
	if r.Start != nil && r.End != nil {
		if r.useBetween {
			st.Where.Add(expr.Field(field, expr.Between(r.Start, r.End)))
		} else {
			st.Where.Add(expr.Field(field, expr.InRange(r.Start, r.End)))
		}
		return
	}
	if r.End != nil {
		if r.useBetween {
			st.Where.Add(expr.Field(field, expr.Lte(*r.End)))
		} else {
			st.Where.Add(expr.Field(field, expr.Lt(*r.End)))
		}
		return
	}
	st.Where.Add(expr.Field(field, expr.Gte(*r.Start)))
}
