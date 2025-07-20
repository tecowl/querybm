package bookswithenum

import (
	"mysql-test/models"
	"time"

	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/helpers/ranges"
	"github.com/tecowl/querybm/helpers/slices"
)

type Condition struct {
	IsbnPrefix     string
	BookTypes      []models.BooksBookType
	Title          string
	YrRange        ranges.Range[int32]
	AvailableRange ranges.Range[time.Time]
}

var _ querybm.Condition = (*Condition)(nil)

func (c *Condition) Build(s *querybm.Statement) {
	if c.IsbnPrefix != "" {
		s.Where.Add(Field("isbn", LikeStartsWith(c.IsbnPrefix)))
	}
	if len(c.BookTypes) > 0 && !slices.All(BookTypeAll, slices.Bind(c.BookTypes, slices.Contains)) {
		s.Where.Add(Field("book_type", EqOrIn(slices.Generalize(c.BookTypes)...)))
	}
	if c.Title != "" {
		s.Where.Add(Field("title", LikeContains(c.Title)))
	}
	if c.YrRange.Start != nil {
		s.Where.Add(Field("yr", Gte(*c.YrRange.Start)))
	}
	if c.YrRange.End != nil {
		s.Where.Add(Field("yr", Lt(*c.YrRange.End)))
	}
	if c.AvailableRange.Start != nil {
		s.Where.Add(Field("available", Gte(*c.AvailableRange.Start)))
	}
	if c.AvailableRange.End != nil {
		s.Where.Add(Field("available", Lt(*c.AvailableRange.End)))
	}
}

var BookTypeAll = []models.BooksBookType{
	models.BooksBookTypeMAGAZINE,
	models.BooksBookTypePAPERBACK,
	models.BooksBookTypeHARDCOVER,
}
