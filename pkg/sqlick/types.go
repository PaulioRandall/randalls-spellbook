package sqlick

import (
	"errors"
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

	// ErrNoSqlType is used when there is no mapping from a
	// Go type to a Sql type.
	ErrNoSqlType = errors.New(
		"No mapping for Go to SQL type",
	)

	// ErrTooFewRows is used when attempting to generate
	// INSERT SQL but the number of rows for insert is less
	// than one.
	ErrTooFewRows = errors.New(
		"Not enough rows to generate insert",
	)

	// ErrNoTableForType is used when an object is passed
	// to a Sqlick function who's type does not map to a
	// table. This means the type was not added via
	// Sqlick.Register before use.
	ErrNoTableForType = errors.New(
		"No matching table for object type",
	)
)

type Sqlick interface {
	// Open opens the database. The path will usually be
	// provided during object construction.
	//
	//		err := db.Open()
	//		// YUDO Handle error.
	//		defer db.Close()
	Open() error

	// IsOpen returns true if the database is open.
	//
	//		if db.IsOpen() {
	//			// YUDO Whatever you like.
	//		}
	IsOpen() bool

	// Close closes the database. As standard, recommended to
	// use with defer.
	//
	//		err := db.Open()
	//		// YUDO Handle error.
	//		defer db.Close()
	Close() error

	// Register accepts a struct object (model), parses it
	// into a Table, and registers the Table. If no error is
	// returned, calls to database interaction functions like
	// CreateTables, Insert, and Select may be made using
	// objects of the same struct type as model.
	//
	// The model must be a struct with at least one exported
	// field or an error is returned. Only exported fields
	// are parsed as part of the Table and field types are
	// currently limited to int64, float64, and string;
	// this will be expanded in future. The first exported
	// field is designated the primary key, regardless of
	// type. It's is recommended to use int64 for primary
	// keys; some databases may require integers while others
	// are less performant with non-int keys, e.g. SQLite
	// works best with integers but the benefits aren't
	// noticable unless you're storing and querying large
	// datasets.
	//
	//		type Model struct {
	//			Id int64
	//			Name string
	//			volume float64 // This field is ignored.
	//		}
	//
	//		err := db.Register(Model{})
	Register(model any) error

	// CreateTables creates the registered Tables within the
	// database where a Table has not yet been created. While
	// it's safe to call Register and then CreateTables at
	// any time, it's recommended to do it all upfront,
	// straight after calling Open.
	//
	//		// YUDO db.Register of types first.
	//		err := db.CreateTables()
	CreateTables() error

	// Insert inserts the object into the database. The
	// object's type must match a registered type or an error
	// is returned.
	//
	//		object := Model{
	//			Id: 123,
	//			Name: "Alice",
	//		}
	//
	//		err := db.Insert(object)
	Insert(object any) error

	// Update updates the object within the database. The
	// object's type must match a registered type or an error
	// is returned. All fields are updated except the ID
	// field, which is used to determine which record to
	// update.
	//
	//		object := Model{
	//			Id: 123,
	//			Name: "Alice",
	//		}
	//		err := db.Insert(object)
	//		// YUDO Handle error.
	//
	//		object.Name = "Bob"
	//		err = db.Update(object)
	Update(object any) error

	// SelectAll returns all records for the table associated
	// with the passed model. The model's type must be
	// registered (Register()) and table created
	// (CreateTables()) for the select to return without
	// error.
	//
	//		slice, err := SelectAll(Model{})
	SelectAll(model any) (any, error)

	// SelectById returns the record with the given id from
	// the table associated with the passed model. The
	// model's type must be registered (Register()) and table
	// created (CreateTables()) for the select to return
	// without error.
	//
	//		slice, err := SelectById(Model{}, 123)
	//SelectById(model any, id any) (any, error)

	// SelectByIdInto queries for the record with the given
	// id from the table associated with the type of the
	// passed object. The row values are placed in the
	// object. The model's type must be registered
	// (Register()) and table created (CreateTables()) for
	// the select to return without error.
	//
	//		err := SelectByIdInto(&object, 123)
	//SelectByIdInto(object *any, id any) (any, error)

	// Select queries the database for one or many records.
	// It calls one of the other select functions based on
	// the arguments. The model's type must be registered
	// (Register()) and table created (CreateTables()) for
	// the select to return without error.
	//
	// The model parameter determines the return type. If
	// it's an array then zero or multiple records may be
	// returned. If it's a struct then either a single
	// record or an empty record is returned.
	//
	// If the id parameter is nil then all records are
	// returned else the record with the specified ID will be
	// returned. Attempting to pass a struct as the model
	// without an id will result in an error.
	//
	//		// Select all records.
	//		slice, err := Select([]Model{}, nil)
	//		sliceOfModel := slice.([]Model)
	//
	//		// Select all records and append them to the passed
	//		// slice.
	//		sliceOfModel := []Model{}
	//		ptrToSlice, err := Select(&sliceOfModel, nil)
	//
	//		// Select a specific record by ID and return it as
	//		// a slice.
	//		slice, err := Select([]Model{}, id)
	//		sliceOfModel := slice.([]Model)
	//
	//		// Select a specific record by ID and append it to
	//		// the passed slice.
	//		sliceOfModel := []Model{}
	//		ptrToSlice, err := Select(&sliceOfModel, id)
	//
	//		// Select a specific record by ID.
	//		object, err := Select(Model{}, id)
	//
	//		// Select a specific record by ID and put it's
	//		// values into the passed object.
	//		object := Model{}
	//		ptrToObject, err := Select(&object, id)
	//
	// Any other configuration produces will produce an
	// error.
	//Select(model any, id any) (any, error)
}

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
	TableInsert(Table, int) (string, error)

	// TableSelectAll returns a SELECT SQL query for all rows
	// in the given Table.
	TableSelectAll(Table) (string, error)

	// TableSelectById returns a SELECT SQL query for the row
	// with a certain ID in the given Table.
	TableSelectById(Table) (string, error)

	// TableUpdateById returns an UPDATE SQL query for the
	// row with a certain ID in the given table.
	TableUpdateById(Table) (string, error)

	// TableDeleteAll returns a DELETE FROM SQL query for
	// all rows in the given Table.
	TableDeleteAll(Table) (string, error)

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

// Column represents the base information required to
// construct a SQL database table column. It is derived
// from a Go struct's field via the Parse function.
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
