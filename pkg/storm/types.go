package storm

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	strZero     string
	int64Zero   int64
	float64Zero float64

	int64ArrZero []int64

	strType     = reflect.TypeOf(strZero)
	int64Type   = reflect.TypeOf(int64Zero)
	float64Type = reflect.TypeOf(float64Zero)

	int64ArrType = reflect.TypeOf(int64ArrZero)

	validKinds = []reflect.Kind{
		reflect.String,
		reflect.Int64,
		reflect.Float64,
	}
)

// Table holds metadata about a SQL table. It is
// constructed from and maps to a model struct.
type Table struct {
	// GoType is the struct's type.
	GoType reflect.Type

	// GoName is the struct's name.
	GoName string

	// Columns holds information about the table's columns,
	// i.e. all exported fields in the struct.
	Columns []Column
}

// NumColumn returns the number of columns.
func (tbl *Table) NumColumn() int {
	return len(tbl.Columns)
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

// Column holds metadata about a Table column, i.e. an
// exported field of a model.
type Column struct {
	// GoName is the field name of the Go struct field.
	GoName string

	// GoType is the reflect.Type of the Go StructField.
	GoType reflect.Type

	// GoIndex is the index of the Go StructField within its
	// struct.
	GoIndex int
}

// String returns the human readable string representation
// of a Column.
func (col *Column) String() string {
	return fmt.Sprintf(
		"[%d] %s: %s",
		col.GoIndex,
		col.GoName,
		col.GoType.Name(),
	)
}

// Zero returns the zero value of the columns GoType.
func (col *Column) Zero() any {
	return reflect.Zero(col.GoType).Interface()
}
