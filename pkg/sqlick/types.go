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

	// ErrNotStruct is returned when attempting to use a
	// model type with a non-struct kind.
	ErrNotStruct = errors.New("Model must be a struct")

	// ErrNotPublic is returned when a model type is not
	// public (exported).
	ErrNotPublic = errors.New("Model struct must be public")

	// ErrBadFieldKind is returned when a model's type
	// contains an unsupported kind for one of its exported
	// fields.
	ErrBadFieldKind = errors.New(
		"Model struct has unsupported field kind",
	)

	// ErrMissFields is returned when a model's type has no/
	// public (exported) fields. Every table must have at
	// least one column.
	ErrMissFields = errors.New(
		"Model struct must have at least one exported field",
	)

	// ErrNoSqlType is returned when there is no mapping from
	// a Go type to an database specific Sql type, e.g, for
	// SQLite 'int => INTEGER' but there is no mapping for
	// 'error => _'.
	ErrNoSqlType = errors.New(
		"No mapping for Go to SQL type",
	)

	// ErrTooFewRows is returned when passing a non-positive
	// number to SqlGenerator functions that can generate
	// queries working with multiple records, e.g. INSERT.
	ErrTooFewRows = errors.New(
		"Not enough rows to generate insert",
	)

	// ErrNoTableForType is returned when an object is passed
	// to a function which does not have a registered table
	// for its type.
	ErrNoTableForType = errors.New(
		"No matching table for object type",
	)
)

// Sqlick is an interface to the database with intentions
// for SQL dialect specific implementations.
//
// All receiving functions require the database to be open
// except for Open, IsOpen, and Close.
type Sqlick interface {
	// Open opens the database. The path will usually be
	// provided during object construction.
	//
	//		err := db.Open()
	//		// YUDO: Handle error.
	//		defer db.Close()
	Open() error

	// IsOpen returns true if the database is open.
	//
	//		if db.IsOpen() {
	//			// YUDO: Whatever you like.
	//		}
	IsOpen() bool

	// Close closes the database. As standard, recommended to
	// use with defer.
	//
	//		err := db.Open()
	//		// YUDO: Handle error.
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
	//		// YUDO: db.Register of types first.
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
	//		// YUDO: Handle error.
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
	//		// Select a specific record by ID.
	//		object, err := Select(Model{}, id)
	//
	// Any other configuration will produce an error.
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
