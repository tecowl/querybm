package bookswithenum

import (
	"time"

	"github.com/tecowl/querybm"
	"github.com/tecowl/querybm/helpers/ranges"
)

type Condition struct {
	YrRange        *ranges.Range[int32]
	AvailableRange *ranges.Range[time.Time]
}

var _ querybm.Condition = (*Condition)(nil)

func (c *Condition) Build(s *querybm.Statement) {
	if c.YrRange != nil {
		c.YrRange.Build("yr", s)
	}
	if c.AvailableRange != nil {
		c.AvailableRange.Build("available", s)
	}
}
