package querybm

import (
	"biz/queries/querybm/statement"
	"context"
	"database/sql"
	"fmt"
)

type Query[M any, C Condition, S Sort] struct {
	db         *sql.DB
	Table      string
	Fields     FieldMapper[M]
	Condition  C
	Sort       S
	Pagination *Pagination
}

func New[M any, C Condition, S Sort](db *sql.DB, c C, s S, table string, fields FieldMapper[M], pagination *Pagination) *Query[M, C, S] {
	return &Query[M, C, S]{
		db:         db,
		Table:      table,
		Fields:     fields,
		Condition:  c,
		Sort:       s,
		Pagination: pagination,
	}
}

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

func (q *Query[M, C, S]) BuildCountSelect() (string, []any) {
	st := statement.NewStatement(q.Table, statement.NewSimpleFields("COUNT(*) AS count"))

	q.Condition.Build(st)

	return st.Build()
}

func (q *Query[M, C, S]) BuildRowsSelect() (string, []any) {
	st := statement.NewStatement(q.Table, q.Fields)
	q.Condition.Build(st)
	q.Sort.Build(st)
	q.Pagination.Build(st)

	return st.Build()
}

func (q *Query[M, C, S]) RowsStatement(ctx context.Context) (*sql.Stmt, []any, error) {
	queryStr, args := q.BuildRowsSelect()
	stmt, err := q.db.PrepareContext(ctx, queryStr)
	if err != nil {
		return nil, nil, err
	}
	return stmt, args, nil
}

func (q *Query[M, C, S]) CountStatement(ctx context.Context) (*sql.Stmt, []any, error) {
	queryStr, args := q.BuildCountSelect()
	stmt, err := q.db.PrepareContext(ctx, queryStr)
	if err != nil {
		return nil, nil, err
	}
	return stmt, args, nil
}

func (q *Query[M, C, S]) Count(ctx context.Context) (int64, error) {
	stmt, args, err := q.CountStatement(ctx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int64
	if err := stmt.QueryRowContext(ctx, args...).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (q *Query[M, C, S]) FirstRow(ctx context.Context) (*sql.Row, error) {
	stmt, args, err := q.RowsStatement(ctx)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return row, nil
}

func (q *Query[M, C, S]) Rows(ctx context.Context) (*sql.Rows, error) {
	stmt, args, err := q.RowsStatement(ctx)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (q *Query[M, C, S]) First(ctx context.Context) (*M, error) {
	row, err := q.FirstRow(ctx)
	if err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	org := new(M)
	if err := q.Fields.Mapper()(row, org); err != nil {
		return nil, err
	}
	return org, nil
}

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
