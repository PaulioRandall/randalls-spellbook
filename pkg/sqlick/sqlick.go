package sqlick

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrNoSqlType is used when there is no mapping from a
	// Go type to a Sql type.
	ErrNoSqlType = errors.New("No mapping for Go to SQL type")
)

// SqlGenerator is an interface implemented by specific
// SQL dialects for generating SQL query strings.
// Implementations assume that the passed Table is valid.
//
// The exact types and nature of the query will vary
// slightly based upon SQL dialect and implementation.
type SqlGenerator interface {
	// TableCreate returns a CREATE TABLE SQL query for the
	// given table.
	TableCreate(Table) (string, error)

	// TableInsert returns an INSERTS INTO SQL query for the
	// given Table.
	TableInsert(Table) (string, error)

	// TableSelect returns a SELECT SQL query for all rows in
	// the given Table.
	TableSelect(Table) (string, error)

	// TableSelectById returns a SELECT SQL query for the row
	// with a certain ID in the given Table.
	TableSelectById(Table) (string, error)

	// TableUpdateById returns an UPDATE SQL query for the
	// row with a certain ID in the given table.
	TableUpdateById(Table) (string, error)

	// TableDeleteById returns a DELETE FROM SQL query for
	// the row with a certain ID in the given Table.
	TableDeleteById(Table) (string, error)
}

// Table represents the base information required to
// construct a SQL database table. It is derived from a
// Go struct via the Parse functions.
type Table struct {
	// GoType is the reflect.Type of the Go struct from which
	// this table was constructed.
	GoType reflect.Type

	// GoName is the name of the Go struct.
	GoName string

	// Columns is all the parsed public fields in the Go
	// struct.
	Columns []Column
}

// IdColumn returns the ID column which is always the first
// column. A valid Table must have at least one column so
// is proper usage this function will never panic with
// out of range message.
func (tbl *Table) IdColumn() Column {
	return tbl.Columns[0]
}

// NonIdColumns returns the all columns except the ID
// column. The ID column is always the first column. If the
// Table only contains an ID column then the returned slice
// will be empty.
func (tbl *Table) NonIdColumns() []Column {
	return tbl.Columns[1:]
}

// String returns the human readable string representation
// of a Table.
func (tbl *Table) String() string {
	sb := strings.Builder{}

	sb.WriteString(tbl.GoName)
	for _, col := range tbl.Columns {
		sb.WriteRune('\n')
		sb.WriteString("  ")
		sb.WriteString(col.String())
	}

	return sb.String()
}

// Column represents the base information required to
// construct a SQL database table column. It is derived
// from a Go struct's field via the Parse function.
type Column struct {
	// GoName is the field name of the Go struct field.
	GoName string

	// GoType is the field name of the type of the Go struct
	// field in the form used by Go programmers.
	GoType string
}

// String returns the human readable string representation
// of a Column.
func (col *Column) String() string {
	return col.GoName + ": " + col.GoType
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

type queryBuilder struct {
	strings.Builder
}

func (qb *queryBuilder) WriteFmt(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	qb.WriteString(s)
}

func genList[T any](
	list []T,
	itemToStr func(int, T) (string, error),
) (string, error) {
	qb := queryBuilder{}

	for i, item := range list {
		if i != 0 {
			qb.WriteString(",\n")
		}

		s, e := itemToStr(i, item)
		if e != nil {
			return "", e
		}

		qb.WriteString(s)
	}

	return qb.String(), nil
}
