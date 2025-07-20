package bookswithenum

import (
	"fmt"

	"github.com/tecowl/querybm/helpers"
)

type BookType string

const (
	Magazine  BookType = "MAGAZINE"
	Paperback BookType = "PAPERBACK"
	Hardcover BookType = "HARDCOVER"
)

var BookTypeAll = []BookType{
	Magazine,
	Paperback,
	Hardcover,
}

func (bt *BookType) Scan(srcRaw interface{}) error {
	var src string
	switch s := srcRaw.(type) {
	case []byte:
		src = string(s)
	case string:
		src = s
	default:
		return fmt.Errorf("unsupported scan type for BookType: %T", srcRaw)
	}
	if !helpers.SliceContains(BookTypeAll, BookType(src)) {
		return fmt.Errorf("invalid BookType: %s", src)
	}
	*bt = BookType(src)
	return nil
}
