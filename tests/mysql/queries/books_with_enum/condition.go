package bookswithenum

import (
	"mysql-test/models"
	"time"

	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/helpers"
)

type Condition struct {
	IsbnPrefix     string
	BookTypes      []models.BooksBookType
	Title          string
	YrRange        helpers.Range[int32]
	AvailableRange helpers.Range[time.Time]
}

func (c *Condition) Build(s *querybm.Statement) {
	if c.IsbnPrefix != "" {
		s.Where.Add(Field("isbn", LikeStartsWith(c.IsbnPrefix)))
	}
	if len(c.BookTypes) > 0 && !helpers.SliceAll(c.BookTypes, helpers.SliceBind(c.BookTypes, helpers.SliceContains)) {
		s.Where.Add(Field("book_type", EqOrIn(helpers.GeneralizeSlice(c.BookTypes))))
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
