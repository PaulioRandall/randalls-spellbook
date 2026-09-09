package storm

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

var goKindToSqliteTypeMappings = map[reflect.Kind]string{
	reflect.String:  "TEXT",
	reflect.Int64:   "INTEGER",
	reflect.Float64: "REAL",
}

type Storm struct {
	path   string
	tables []Table
	db     *sql.DB
}

func New(path string) *Storm {
	return &Storm{
		path: path,
	}
}

// Open opens the database. If not an in-memory path then
// the missing directories in the directory path are
// created.
//
//	err := db.Open()
//	// YUDO: Handle error.
//	defer db.Close()
func (st *Storm) Open() error {
	e := st.mkdirs()
	if e != nil {
		return e
	}

	db, e := sql.Open("sqlite", st.path)
	if e != nil {
		return fmt.Errorf(
			"Unable to open SQLite database: %w",
			e,
		)
	}

	st.db = db
	return nil
}

func (st *Storm) mkdirs() error {
	if st.path == ":memory" {
		// SQlite in-memory database. There is no path!
		return nil
	}

	parent := filepath.Dir(st.path)
	e := os.MkdirAll(parent, os.ModePerm)
	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Unable to check or create directory path to SQLite database: %w",
		e,
	)
}

// IsOpen returns true if the database is open.
//
//	if db.IsOpen() {
//		// YUDO.
//	}
func (st *Storm) IsOpen() bool {
	return st.db != nil
}

// Close closes the database. Use with defer as usual.
//
//	err := db.Open()
//	// YUDO: Handle error.
//	defer db.Close()
func (st *Storm) Close() error {
	if !st.IsOpen() {
		return nil
	}

	defer func() {
		st.db = nil
	}()

	return st.db.Close()
}

// Create parses the passed model and creates a table
// of it in the database. Passing a model for a table that
// already exists does nothing unless the model was passed
// to Drop first.
//
// After passing a model to Create, calls to database
// interaction functions like Insert and Select can now be
// made using objects of the same struct type as model.
//
// While it's safe to call Create at anytime, it's
// recommended to create all tables upfront, straight after
// calling Open. The model must be a struct with at least
// one exported field or an error is returned. Only
// exported fields are parsed as part of the Table and
// field types are currently limited to int64, float64,
// and string; this will be expanded in future. The first exported
// field is designated the primary key, regardless of
// type. It's is recommended to use int64 for primary
// keys; some databases may require integers while others
// are less performant with non-int keys, e.g. SQLite
// works best with integers but the benefits aren't
// noticable unless you're storing and querying large
// datasets.
//
//	type Person struct {
//		Id int64
//		Name string
//		Height float64
//		ignored int64 // This field is ignored.
//	}
//
//	err := db.Create(Person{})
func (st *Storm) Create(model any) error {
	e := st.register(model)
	if e != nil {
		return fmt.Errorf(
			"Failed to create table: %w",
			e,
		)
	}

	typ := reflect.TypeOf(model)
	table, found := st.findTableFor(typ)
	if !found {
		return fmt.Errorf(
			"No table exists for struct '%s': %w",
			typ.Name(),
			ErrNoTableForType,
		)
	}

	e = st.createTable(table)
	if e != nil {
		return fmt.Errorf(
			"Failed to create table '%s': %w",
			table.GoName,
			e,
		)
	}

	return nil
}

func (st *Storm) register(model any) error {
	_, found := st.findTableFor(reflect.TypeOf(model))
	if found {
		return nil
	}

	table, e := Parse(model)

	if e != nil {
		return fmt.Errorf(
			"Failed to register struct/table: %w",
			e,
		)
	}

	st.tables = append(st.tables, table)
	return nil
}

func (st *Storm) createTable(tbl Table) error {
	query := nidoking.Given(`
		CREATE TABLE IF NOT EXISTS {{table}} (
			{{columns}},
		  PRIMARY KEY ({{id_column}})
		)
	`).
		Fmt("table", tbl.GoName).
		Join("columns", "", tbl.ColumnNames()...).
		Fmt("id_column", tbl.IdColumn().GoName).
		String()

	_, e := st.db.Exec(query)
	return e
}

// Insert inserts the object into the database. The
// object's type must match a registered type or an error
// is returned.
//
//	object := Model{
//		Id: 123,
//		Name: "Alice",
//	}
//
//	err := db.Insert(object)
func (st *Storm) Insert(object any) error {
	value := reflect.ValueOf(object)
	e := st.insertValue(value)

	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to insert object into database: %w",
		e,
	)
}

func (st *Storm) insertValue(
	value reflect.Value,
) error {
	table, found := st.findTableFor(value.Type())

	if !found {
		return fmt.Errorf(
			"No table exists for struct '%s': %w",
			value.Type().Name(),
			ErrNoTableForType,
		)
	}

	fieldValues := st.listOrderedFieldValues(
		table.Columns,
		value,
	)
	return st.execInsert(table, fieldValues)
}

func (st *Storm) findTableFor(
	typ reflect.Type,
) (Table, bool) {
	for _, table := range st.tables {
		if table.GoType == typ {
			return table, true
		}
	}
	return Table{}, false
}

func (st *Storm) findTableForOrError(
	typ reflect.Type,
) (Table, error) {
	table, found := st.findTableFor(typ)

	if found {
		return table, nil
	}

	return Table{}, fmt.Errorf(
		"No table exists for struct '%s': %w",
		typ.Name(),
		ErrNoTableForType,
	)
}

func (st *Storm) listOrderedFieldValues(
	columns []Column,
	value reflect.Value,
) []any {
	result := make([]any, len(columns))

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		field := value.FieldByName(col.GoName)
		result[i] = field.Interface()
	}

	return result
}

func (st *Storm) execInsert(
	tbl Table,
	fieldValues []any,
) error {
	query := nidoking.Given(`
		INSERT INTO {{table}} (
		  {{columns}}
		)
		VALUES (
			{{q_marks}}
		)
	`).
		Fmt("table", tbl.GoName).
		Join("columns", ",", tbl.ColumnNames()...).
		Repeat("q_marks", ",", "?", tbl.ColumnCount()).
		String()

	_, e := st.db.Exec(query, fieldValues...)
	if e != nil {
		return fmt.Errorf(
			"Failed to insert into table '%s': %w",
			tbl.GoName,
			e,
		)
	}

	return nil
}

// Update updates the object within the database. The
// object's type must match a registered type or an error
// is returned. All fields are updated except the ID
// field, which is used to determine which record to
// update.
//
//	object := Model{
//		Id: 123,
//		Name: "Alice",
//	}
//	err := db.Insert(object)
//	// YUDO: Handle error.
//
//	object.Name = "Bob"
//	err = db.Update(object)
func (st *Storm) Update(object any) error {
	value := reflect.ValueOf(object)
	e := st.updateValue(value)

	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to update object in database: %w",
		e,
	)
}

func (st *Storm) updateValue(
	value reflect.Value,
) error {
	table, e := st.findTableForOrError(value.Type())
	if e != nil {
		return e
	}

	fieldValues := st.listOrderedFieldValues(
		table.Columns,
		value,
	)

	// Move ID value to end (for the WHERE clause)
	fieldValues = append(fieldValues[1:], fieldValues[0])
	return st.execUpdate(table, fieldValues)
}

func (st *Storm) execUpdate(
	tbl Table,
	fieldValues []any,
) error {
	query := nidoking.Given(`
		UPDATE
			{{table}}
		SET
			{{columns}} = ?
		WHERE
			{{id_column}} = ?
	`).
		Fmt("table", tbl.GoName).
		Join("columns", ",", tbl.ColumnNames()[1:]...).
		Fmt("id_column", tbl.IdColumn().GoName).
		String()

	_, e := st.db.Exec(query, fieldValues...)
	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to update row in table '%s': %w",
		tbl.GoName,
		e,
	)
}

// SelectAll returns all records for the table associated
// with the passed model. The model's type must be
// registered (Register()) and table created
// (CreateTables()) for the select to return without
// error.
//
//	slice, err := SelectAll(Model{})
func (st *Storm) SelectAll(
	model any,
) (any, error) {
	typ := reflect.TypeOf(model)
	result, e := st.selectAllOfType(typ)

	if e == nil {
		return result, nil
	}

	return nil, fmt.Errorf(
		"Failed to select all from database: %w",
		e,
	)
}

func (st *Storm) selectAllOfType(
	typ reflect.Type,
) (any, error) {
	table, e := st.findTableForOrError(typ)
	if e != nil {
		return nil, e
	}

	fieldValues := st.listOrderedFieldValues(
		table.Columns,
		reflect.Zero(typ),
	)

	return st.querySelectAll(table, fieldValues)
}

func (st *Storm) querySelectAll(
	tbl Table,
	fieldValues []any,
) (any, error) {
	query := nidoking.Given(`
		SELECT
			{{columns}}
		FROM
			{{table}}
	`).
		Join("columns", ",", tbl.ColumnNames()...).
		Fmt("table", tbl.GoName).
		String()

	rows, e := st.db.Query(query, fieldValues...)
	if e != nil {
		return nil, fmt.Errorf(
			"Failed to select all from table '%s': %w",
			tbl.GoName,
			e,
		)
	}

	result, e := st.scanSelectedRows(tbl, rows)
	if e != nil {
		return nil, e
	}

	return result, nil
}

func (st *Storm) scanSelectedRows(
	table Table,
	rows *sql.Rows,
) (any, error) {
	values, valuePtrs := createValueContainers(table)

	sliceValue := reflect.MakeSlice(
		reflect.SliceOf(table.GoType),
		0,
		0,
	)

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, fmt.Errorf(
				"Failed to scan row %d: %w",
				i,
				e,
			)
		}

		object := createNewObjectWithValues(
			table,
			values,
		)

		sliceValue = reflect.Append(
			sliceValue,
			reflect.ValueOf(object),
		)
	}

	e := rows.Err()
	if e != nil {
		return nil, fmt.Errorf(
			"Error occurred after scanning database rows: %w",
			e,
		)
	}

	return sliceValue.Interface(), nil
}

func createValueContainers(
	tbl Table,
) ([]any, []any) {
	colCount := tbl.ColumnCount()
	values := make([]any, colCount, colCount)
	valuePtrs := make([]any, colCount, colCount)

	for i, col := range tbl.Columns {
		values[i] = col.Zero()
		valuePtrs[i] = &values[i]
	}

	return values, valuePtrs
}

func createNewObjectWithValues(
	table Table,
	values []any,
) any {
	objectValue := reflect.New(table.GoType).Elem()

	for i, col := range table.Columns {
		field := objectValue.Field(col.GoIndex)
		value := reflect.ValueOf(values[i])
		field.Set(value)
	}

	return objectValue.Interface()
}

func toSliceOfType(
	typ reflect.Type,
	values []any,
) any {
	sliceValue := reflect.MakeSlice(
		reflect.SliceOf(typ),
		len(values),
		len(values),
	)

	for i, v := range values {
		itemValue := reflect.ValueOf(v)
		sliceValue.Index(i).Set(itemValue)
	}

	return sliceValue.Interface()
}

// SelectById returns the record with the given id from
// the table associated with the passed model. If no
// record is found then an error is returned. The
// model's type must be registered (Register()) and table
// created (CreateTables()) for the select to return
// without error.
//
//	slice, err := SelectById(Model{}, 123)
func (st *Storm) SelectById(
	model any,
	id any,
) (any, error) {
	typ := reflect.TypeOf(model)
	result, e := st.selectByIdOfType(typ, id)
	if e == nil {
		return result, nil
	}

	return nil, fmt.Errorf(
		"Failed to select by ID from database: %w",
		e,
	)
}

func (st *Storm) selectByIdOfType(
	typ reflect.Type,
	id any,
) (any, error) {
	tbl, e := st.findTableForOrError(typ)
	if e != nil {
		return nil, e
	}

	e = st.validateModelIdType(tbl, id)
	if e != nil {
		return nil, e
	}

	return st.querySelectById(tbl, id)
}

func (st *Storm) validateModelIdType(
	tbl Table,
	id any,
) error {
	idTyp := reflect.TypeOf(id)

	if tbl.IdColumn().GoType.Kind() != idTyp.Kind() {
		return ErrBadIdType
	}

	return nil
}

func (st *Storm) querySelectById(
	tbl Table,
	id any,
) (any, error) {
	query := nidoking.Given(`
		SELECT
			{{columns}}
		FROM
			{{table}}
		WHERE
			{{id_column}} = ?
	`).
		Join("columns", ",", tbl.ColumnNames()...).
		Fmt("table", tbl.GoName).
		Fmt("id_column", tbl.IdColumn().GoName).
		String()

	rows, e := st.db.Query(query, id)
	if e != nil {
		return nil, fmt.Errorf(
			"Failed to select by ID from table '%s': %w",
			tbl.GoName,
			e,
		)
	}

	result, e := st.scanSelectedRows(tbl, rows)
	if e != nil {
		return nil, e
	}

	object, ok := getFirstItemIfArray(result)
	if !ok {
		return nil, fmt.Errorf(
			"Failed to select by ID '%v' from table '%s': %w",
			id,
			tbl.GoName,
			ErrRecordNotFound,
		)
	}

	return object, nil
}

func getFirstItemIfArray(v any) (any, bool) {
	rv := reflect.ValueOf(v)

	isArray := rv.Kind() == reflect.Array
	isSlice := rv.Kind() == reflect.Slice

	if !isArray && !isSlice {
		return nil, false
	}

	if rv.Len() > 0 {
		return rv.Index(0).Interface(), true
	}

	return nil, false
}

func printQuery(query string, values []any) {
	for _, v := range values {
		vs := fmt.Sprintf("%v", v)
		query = strings.Replace(query, "?", vs, 1)
	}
	println(query)
}

// DeleteById removes the record with the given id from
// the table associated with the passed model. If no
// record is found then nothing happens.
//
//	e := DeleteById(Model{}, 123)
func (st *Storm) DeleteById(
	model any,
	id any,
) error {
	typ := reflect.TypeOf(model)
	e := st.deleteByIdOfType(typ, id)
	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to delete by ID from database: %w",
		e,
	)
}

func (st *Storm) deleteByIdOfType(
	typ reflect.Type,
	id any,
) error {
	tbl, e := st.findTableForOrError(typ)
	if e != nil {
		return e
	}

	e = st.validateModelIdType(tbl, id)
	if e != nil {
		return e
	}

	return st.execDeleteById(tbl, id)
}

func (st *Storm) execDeleteById(
	tbl Table,
	id any,
) error {
	query := nidoking.Given(`
		DELETE FROM
			{{table}}
		WHERE
			{{id_column}} = ?
	`).
		Fmt("table", tbl.GoName).
		Fmt("id_column", tbl.IdColumn().GoName).
		String()

	_, e := st.db.Exec(query, id)
	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to delete by ID '%v' from table '%s': %w",
		id,
		tbl.GoName,
		e,
	)
}

// Drop removes a table from the database, deleting all
// records in the process. Passing a model for a table that
// doesn't exists does nothing.
//
//	type Person struct {
//		Id int64
//		Name string
//		Height float64
//		ignored int64 // This field is ignored.
//	}
//
//	err := db.Create(Person{})
//	// YUDO: Handle error.
//
//	err := db.Drop(Person{})
func (st *Storm) Drop(model any) error {
	typ := reflect.TypeOf(model)
	tbl, found := st.findTableFor(typ)
	if !found {
		return nil
	}

	query := nidoking.Given(`
		DROP TABLE IF EXISTS {{table}}
	`).
		Fmt("table", tbl.GoName).
		String()

	_, e := st.db.Exec(query)
	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to drop table '%s': %w",
		tbl.GoName,
		e,
	)
}

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
