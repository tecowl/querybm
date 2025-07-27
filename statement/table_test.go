package statement

import (
	"reflect"
	"testing"
)

func TestNewTableBlock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		tableName     string
		expectedTable string
		expectedAlias string
		expectedFrom  string
	}{
		{
			name:          "Create table block with simple name",
			tableName:     "users",
			expectedTable: "users",
			expectedAlias: "",
			expectedFrom:  "users",
		},
		{
			name:          "Create table block with schema",
			tableName:     "public.users",
			expectedTable: "public.users",
			expectedAlias: "",
			expectedFrom:  "public.users",
		},
		{
			name:          "Create table block with alias",
			tableName:     "users u",
			expectedTable: "users",
			expectedAlias: "u",
			expectedFrom:  "users u",
		},
		{
			name:          "Create table block with AS alias",
			tableName:     "users AS u",
			expectedTable: "users",
			expectedAlias: "u",
			expectedFrom:  "users AS u",
		},
		{
			name:          "Create table block with lower case AS alias",
			tableName:     "users as u",
			expectedTable: "users",
			expectedAlias: "u",
			expectedFrom:  "users AS u",
		},
		{
			name:          "Create table block with unexpected spaces 1",
			tableName:     "users `special users`",
			expectedTable: "users",
			expectedAlias: "`special users`",
			expectedFrom:  "users `special users`",
		},
		{
			name:          "Create table block with unexpected spaces 2",
			tableName:     "`special users` users",
			expectedTable: "`special",
			expectedAlias: "users` users",
			expectedFrom:  "`special users` users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tb := NewTableBlock(tt.tableName)
			if tb.tableName.Name != tt.expectedTable {
				t.Errorf("NewTableBlock() table name = %v, want %v", tb.tableName.Name, tt.expectedTable)
			}
			if tb.tableName.Alias != tt.expectedAlias {
				t.Errorf("NewTableBlock() alias = %v, want %v", tb.tableName.Alias, tt.expectedAlias)
			}
			r, args := tb.Build()
			if r != tt.expectedFrom {
				t.Errorf("NewTableBlock() content = %v, want %v", r, tt.expectedFrom)
			}
			if len(args) != 0 {
				t.Errorf("NewTableBlock() values = %v, want empty slice", args)
			}
		})
	}
}

func TestTableBlock_InnerJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		initial string
		joins   []struct {
			table     string
			condition string
			values    []any
		}
		wantContent string
		wantValues  []any
	}{
		{
			name:    "Single INNER JOIN",
			initial: "users u",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "profiles p",
					condition: "u.id = p.user_id",
					values:    nil,
				},
			},
			wantContent: "users u INNER JOIN profiles p ON u.id = p.user_id",
			wantValues:  []any{},
		},
		{
			name:    "INNER JOIN with parameter",
			initial: "orders o",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "customers c",
					condition: "o.customer_id = c.id AND c.status = ?",
					values:    []any{"active"},
				},
			},
			wantContent: "orders o INNER JOIN customers c ON o.customer_id = c.id AND c.status = ?",
			wantValues:  []any{"active"},
		},
		{
			name:    "INNER JOIN duplicated alias name",
			initial: "orders o",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "customers c",
					condition: "o.customer_id = c.id AND c.status = ?",
					values:    []any{"active"},
				},
				{
					table:     "customers c",
					condition: "o.customer_id = c.id AND c.status = ? AND c.value IS NULL", // This will be ignored
					values:    []any{"inactive"},                                           // This will be ignored
				},
			},
			wantContent: "orders o INNER JOIN customers c ON o.customer_id = c.id AND c.status = ?",
			wantValues:  []any{"active"},
		},
		{
			name:    "INNER JOIN duplicated table name",
			initial: "orders",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "customers",
					condition: "orders.customer_id = customers.id AND customers.status = ?",
					values:    []any{"active"},
				},
				{
					table:     "customers",
					condition: "orders.customer_id = customers.id AND customers.status = ? AND customers.value IS NULL", // This will be ignored
					values:    []any{"inactive"},                                                                        // This will be ignored
				},
			},
			wantContent: "orders INNER JOIN customers ON orders.customer_id = customers.id AND customers.status = ?",
			wantValues:  []any{"active"},
		},
		{
			name:    "Multiple INNER JOINs",
			initial: "orders o",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "customers c",
					condition: "o.customer_id = c.id",
					values:    nil,
				},
				{
					table:     "products p",
					condition: "o.product_id = p.id",
					values:    nil,
				},
			},
			wantContent: "orders o INNER JOIN customers c ON o.customer_id = c.id INNER JOIN products p ON o.product_id = p.id",
			wantValues:  []any{},
		},
		{
			name:    "Multiple INNER JOINs with parameters",
			initial: "posts p",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "users u",
					condition: "p.user_id = u.id AND u.status = ?",
					values:    []any{"active"},
				},
				{
					table:     "categories c",
					condition: "p.category_id = c.id AND c.visible = ?",
					values:    []any{true},
				},
			},
			wantContent: "posts p INNER JOIN users u ON p.user_id = u.id AND u.status = ? INNER JOIN categories c ON p.category_id = c.id AND c.visible = ?",
			wantValues:  []any{"active", true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tb := NewTableBlock(tt.initial)
			for _, join := range tt.joins {
				tb.InnerJoin(join.table, join.condition, join.values...)
			}
			r, args := tb.Build()
			if r != tt.wantContent {
				t.Errorf("InnerJoin() content = %v, want %v", r, tt.wantContent)
			}
			if len(args) != len(tt.wantValues) {
				if !reflect.DeepEqual(args, tt.wantValues) {
					t.Errorf("InnerJoin() values = %v, want %v", args, tt.wantValues)
				}
			}
		})
	}
}

func TestTableBlock_LeftOuterJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		initial string
		joins   []struct {
			table     string
			condition string
			values    []any
		}
		wantContent string
		wantValues  []any
	}{
		{
			name:    "Single LEFT OUTER JOIN",
			initial: "users u",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "profiles p",
					condition: "u.id = p.user_id",
					values:    nil,
				},
			},
			wantContent: "users u LEFT OUTER JOIN profiles p ON u.id = p.user_id",
			wantValues:  []any{},
		},
		{
			name:    "LEFT OUTER JOIN with parameter",
			initial: "posts p",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "comments c",
					condition: "p.id = c.post_id AND c.deleted_at IS NULL AND c.status = ?",
					values:    []any{"published"},
				},
			},
			wantContent: "posts p LEFT OUTER JOIN comments c ON p.id = c.post_id AND c.deleted_at IS NULL AND c.status = ?",
			wantValues:  []any{"published"},
		},
		{
			name:    "Multiple LEFT OUTER JOINs",
			initial: "products p",
			joins: []struct {
				table     string
				condition string
				values    []any
			}{
				{
					table:     "reviews r",
					condition: "p.id = r.product_id",
					values:    nil,
				},
				{
					table:     "categories c",
					condition: "p.category_id = c.id",
					values:    nil,
				},
			},
			wantContent: "products p LEFT OUTER JOIN reviews r ON p.id = r.product_id LEFT OUTER JOIN categories c ON p.category_id = c.id",
			wantValues:  []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tb := NewTableBlock(tt.initial)
			for _, join := range tt.joins {
				tb.LeftOuterJoin(join.table, join.condition, join.values...)
			}
			r, args := tb.Build()
			if r != tt.wantContent {
				t.Errorf("LeftOuterJoin() content = %v, want %v", r, tt.wantContent)
			}
			if len(args) != len(tt.wantValues) {
				if !reflect.DeepEqual(args, tt.wantValues) {
					t.Errorf("LeftOuterJoin() values = %v, want %v", args, tt.wantValues)
				}
			}
		})
	}
}

func TestTableBlock_MixedJoins(t *testing.T) {
	t.Parallel()
	tb := NewTableBlock("orders o")

	// Add INNER JOIN
	tb.InnerJoin("customers c", "o.customer_id = c.id")

	// Add LEFT OUTER JOIN
	tb.LeftOuterJoin("order_discounts d", "o.id = d.order_id AND d.active = ?", true)

	// Add another INNER JOIN
	tb.InnerJoin("products p", "o.product_id = p.id AND p.available = ?", true)

	wantContent := "orders o INNER JOIN customers c ON o.customer_id = c.id LEFT OUTER JOIN order_discounts d ON o.id = d.order_id AND d.active = ? INNER JOIN products p ON o.product_id = p.id AND p.available = ?"
	wantValues := []any{true, true}

	r, args := tb.Build()
	if r != wantContent {
		t.Errorf("Mixed joins content = %v, want %v", r, wantContent)
	}
	if !reflect.DeepEqual(args, wantValues) {
		t.Errorf("Mixed joins values = %v, want %v", args, wantValues)
	}
}
