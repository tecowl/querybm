package querybm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tecowl/querybm/statement"
)

// Query represents a SQL query builder with generic support for models, conditions, and sorting.
// It provides methods to build and execute SELECT queries with pagination support.
type Query[M any, C Condition, S Sort] struct {
	db         *sql.DB
	Table      string
	Fields     FieldMapper[M]
	Condition  C
	Sort       S
	Pagination *Pagination
}

// New creates a new Query instance with the specified database connection, condition, sort, table name, field mapper, and pagination.
func New[M any, C Condition, S Sort](db *sql.DB, table string, fields FieldMapper[M], c C, s S, pagination *Pagination) *Query[M, C, S] {
	return &Query[M, C, S]{
		db:         db,
		Table:      table,
		Fields:     fields,
		Condition:  c,
		Sort:       s,
		Pagination: pagination,
	}
}

// Validate validates the query's condition, sort, and pagination components.
// It returns an error if any component's validation fails.
func (q *Query[M, C, S]) Validate() error {
	if v, ok := any(q.Condition).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("condition validation failed: %w", err)
		}
	}
	if v, ok := any(q.Sort).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("sort validation failed: %w", err)
		}
	}
	if err := q.Pagination.Validate(); err != nil {
		return fmt.Errorf("pagination validation failed: %w", err)
	}
	return nil
}

// BuildCountSelect builds a COUNT(*) query string with the current conditions.
// It returns the SQL query string and its arguments.
func (q *Query[M, C, S]) BuildCountSelect() (string, []any) {
	st := statement.New(q.Table, statement.NewSimpleFields("COUNT(*) AS count"))

	q.Condition.Build(st)

	return st.Build()
}

// BuildRowsSelect builds a SELECT query string with all fields, conditions, sorting, and pagination.
// It returns the SQL query string and its arguments.
func (q *Query[M, C, S]) BuildRowsSelect() (string, []any) {
	st := statement.New(q.Table, q.Fields)
	q.Condition.Build(st)
	q.Sort.Build(st)
	q.Pagination.Build(st)

	return st.Build()
}

// RowsStatement prepares a SELECT statement for retrieving rows.
// It returns the prepared statement, query arguments, and any error that occurred.
func (q *Query[M, C, S]) RowsStatement(ctx context.Context) (*sql.Stmt, []any, error) {
	queryStr, args := q.BuildRowsSelect()
	stmt, err := q.db.PrepareContext(ctx, queryStr)
	if err != nil {
		return nil, nil, err
	}
	return stmt, args, nil
}

// CountStatement prepares a COUNT statement for counting matching rows.
// It returns the prepared statement, query arguments, and any error that occurred.
func (q *Query[M, C, S]) CountStatement(ctx context.Context) (*sql.Stmt, []any, error) {
	queryStr, args := q.BuildCountSelect()
	stmt, err := q.db.PrepareContext(ctx, queryStr)
	if err != nil {
		return nil, nil, err
	}
	return stmt, args, nil
}

// Count executes a COUNT query and returns the number of matching rows.
// It returns 0 if no rows match the conditions.
func (q *Query[M, C, S]) Count(ctx context.Context) (int64, error) {
	stmt, args, err := q.CountStatement(ctx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int64
	if err := stmt.QueryRowContext(ctx, args...).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

// FirstRow executes the query and returns the first row as *sql.Row.
// The caller is responsible for scanning the row.
func (q *Query[M, C, S]) FirstRow(ctx context.Context) (*sql.Row, error) {
	stmt, args, err := q.RowsStatement(ctx)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	if err := row.Err(); err != nil {
		return nil, err
	}
	return row, nil
}

// Rows executes the query and returns the result set as *sql.Rows.
// It prepares the statement, executes it, and returns the rows.
// The caller is responsible for closing the rows.
func (q *Query[M, C, S]) Rows(ctx context.Context) (*sql.Rows, error) {
	stmt, args, err := q.RowsStatement(ctx)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...) // nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// First executes the query and returns the first matching model instance.
// It returns nil if no rows match the conditions.
func (q *Query[M, C, S]) First(ctx context.Context) (*M, error) {
	row, err := q.FirstRow(ctx)
	if err != nil {
		return nil, err
	}
	org := new(M)
	if err := q.Fields.Mapper()(row, org); err != nil {
		return nil, err
	}
	return org, nil
}

// List executes the query and returns all matching model instances as a slice.
// It returns an empty slice if no rows match the conditions.
func (q *Query[M, C, S]) List(ctx context.Context) ([]*M, error) {
	rows, err := q.Rows(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*M
	for rows.Next() {
		org := new(M)
		if err := q.Fields.Mapper()(rows, org); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orgs, nil
}
