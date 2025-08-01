// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type BooksBookType string

const (
	BooksBookTypeMAGAZINE  BooksBookType = "MAGAZINE"
	BooksBookTypePAPERBACK BooksBookType = "PAPERBACK"
	BooksBookTypeHARDCOVER BooksBookType = "HARDCOVER"
)

func (e *BooksBookType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BooksBookType(s)
	case string:
		*e = BooksBookType(s)
	default:
		return fmt.Errorf("unsupported scan type for BooksBookType: %T", src)
	}
	return nil
}

type NullBooksBookType struct {
	BooksBookType BooksBookType
	Valid         bool // Valid is true if BooksBookType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBooksBookType) Scan(value interface{}) error {
	if value == nil {
		ns.BooksBookType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BooksBookType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBooksBookType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BooksBookType), nil
}

type Author struct {
	AuthorID int32
	Name     string
}

type Book struct {
	BookID    int32
	AuthorID  int32
	Isbn      string
	BookType  BooksBookType
	Title     string
	Yr        int32
	Available time.Time
	Tags      string
}
