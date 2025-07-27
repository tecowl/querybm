package querybm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tecowl/querybm/statement"
)

// Query represents a SQL query builder with generic support for models, conditions, and sorting.
// It provides methods to build and execute SELECT queries with limitOffset support.
type Query[M any] struct {
	db          DB
	Table       string
	Fields      FieldMapper[M]
	Condition   Condition
	Sort        Sort
	LimitOffset LimitOffset
}

// New creates a new Query instance with the provided parameters.
// db: The database connection to use for executing queries.
// table: The name of the table to query.
// fields: The field mapper that defines how to map database rows to model instances. This is used for mapping in List method.
// c: The condition to apply to the query. This is used for List and Count methods.
// s: The sort item to apply to the query. This is used for ordering the results in List method.
// limitOffset: The limitOffset settings for the query. This is used to limit the number of results returned in List method.
func New[M any](db *sql.DB, table string, fields FieldMapper[M], c Condition, s Sort, limitOffset LimitOffset) *Query[M] {
	return &Query[M]{
		db:          newDBWrapper(db),
		Table:       table,
		Fields:      fields,
		Condition:   c,
		Sort:        s,
		LimitOffset: limitOffset,
	}
}

// Validate validates the query's condition, sort, and limitOffset components.
// It returns an error if any component's validation fails.
func (q *Query[M]) Validate() error {
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
	if v, ok := any(q.LimitOffset).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("limitOffset validation failed: %w", err)
		}
	}
	return nil
}

// BuildCountSelect builds a COUNT(*) query string with the current conditions.
// It returns the SQL query string and its arguments.
func (q *Query[M]) BuildCountSelect() (string, []any) {
	st := statement.New(q.Table, statement.NewSimpleFields("COUNT(*) AS count"))

	if q.Condition != nil {
		q.Condition.Build(st)
	}

	return st.Build()
}

// BuildRowsSelect builds a SELECT query string with all fields, conditions, sorting, and limitOffset.
// It returns the SQL query string and its arguments.
func (q *Query[M]) BuildRowsSelect() (string, []any) {
	st := statement.New(q.Table, q.Fields)
	if fb, ok := q.Fields.(Builder); ok {
		fb.Build(st)
	}
	if q.Condition != nil {
		q.Condition.Build(st)
	}
	if q.Sort != nil {
		q.Sort.Build(st)
	}
	if q.LimitOffset != nil {
		q.LimitOffset.Build(st)
	}

	return st.Build()
}

// RowsStatement prepares a SELECT statement for retrieving rows.
// It returns the prepared statement, query arguments, and any error that occurred.
func (q *Query[M]) RowsStatement(ctx context.Context) (Stmt, []any, error) { // nolint:ireturn
	queryStr, args := q.BuildRowsSelect()
	stmt, err := q.db.PrepareContext(ctx, queryStr)
	if err != nil {
		return nil, nil, err
	}
	return stmt, args, nil
}

// CountStatement prepares a COUNT statement for counting matching rows.
// It returns the prepared statement, query arguments, and any error that occurred.
func (q *Query[M]) CountStatement(ctx context.Context) (Stmt, []any, error) { // nolint:ireturn
	queryStr, args := q.BuildCountSelect()
	stmt, err := q.db.PrepareContext(ctx, queryStr)
	if err != nil {
		return nil, nil, err
	}
	return stmt, args, nil
}

// Count executes a COUNT query and returns the number of matching rows.
// It returns 0 if no rows match the conditions.
func (q *Query[M]) Count(ctx context.Context) (int64, error) {
	stmt, args, err := q.CountStatement(ctx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int64
	if err := stmt.QueryRowContext(ctx, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// FirstRow executes the query and returns the first row as *sql.Row.
// The caller is responsible for scanning the row.
func (q *Query[M]) FirstRow(ctx context.Context) (Row, error) { // nolint:ireturn
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
func (q *Query[M]) Rows(ctx context.Context) (Rows, error) { // nolint:ireturn
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
func (q *Query[M]) First(ctx context.Context) (*M, error) {
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
func (q *Query[M]) List(ctx context.Context) ([]*M, error) {
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
