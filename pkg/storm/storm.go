package storm

import (
	"database/sql"
	"errors"
	"reflect"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

var (
	// ErrNotOpen occurs when trying to perform an operation
	// before opening the database.
	ErrNotOpen = curse.Proto(
		"Database not open",
	)

	// ErrTableRequest occurs within an error chain when any
	// error occurs involving a specific table/struct, except
	// for 'not found' errors.
	ErrTableRequest = curse.Proto(
		"Request failed for table: %s",
	)

	// ErrScanningRow is returned when an error occurs
	// scanning database results.
	ErrScanningRow = curse.Proto(
		"Scanning row: %d",
	)

	// ErrObjectNotFound is returned when a search for a
	// specific object/row failed.
	ErrObjectNotFound = curse.Proto(
		"Object not found: %s with ID %v",
	)

	// ErrDatabaseFile is returned when an error occurs with
	// or while opening or closing the database.
	ErrDatabaseFile = curse.Proto(
		"Database IO error: %s",
	)

	// ErrNoSuchTable is returned when an object is passed
	// to a function which does not have a registered table
	// for its type.
	ErrNoSuchTable = curse.Proto(
		"No matching table for object type: %s",
	)

	// ErrBadIdType is returned when an ID passed to a
	// function, e.g. SelectById, is not of the same type as
	// the ID field of the associated model type. This may
	// be returned even for compatible types like int when
	// int64 is expected.
	ErrBadIdType = curse.Proto(
		"ID type mismatch for '%s', got %s, want %s",
	)
)

// Storm is the core type and the interface to the SQLite
// database.
type Storm struct {
	path   string
	tables []Table
	db     *sql.DB
}

// New returns a new [Storm] for the database represented
// by path.
func New(path string) *Storm {
	return &Storm{
		path: path,
	}
}

// Open opens the database. If not an 'in-memory' path then
// the missing directories in the directory path are
// created.
//
//	err := db.Open()
//	// YUDO: Handle error.
//	defer db.Close()
func (st *Storm) Open() error {
	if st.IsOpen() {
		return nil
	}

	e := makeParentDirs(st.path)
	if e != nil {
		return e
	}

	db, e := sql.Open("sqlite", st.path)
	if e != nil {
		return curse.Err("Unable to open SQLite database").
			Wraps(e).
			WrapIn(ErrDatabaseFile).
			Fmt(st.path)
	}

	st.db = db
	return nil
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

	e := st.db.Close()
	if e == nil {
		return nil
	}

	return curse.Err("Unable to close SQLite database").
		Wraps(e).
		WrapIn(ErrDatabaseFile).
		Fmt(st.path)
}

// Create parses the passed models and creates tables
// for each within the database, if a table doesn't already
// exist.
//
// After passing a model to Create, calls to database
// interaction functions like [Storm.Insert] and
// [Storm.Select] can now be made using objects of the same
// type as model.
//
// While it's safe to call Create at anytime, it's
// recommended to create all tables upfront, straight after
// calling [Storm.Open]. The model must be a struct with at
// least one exported field or an error is returned. Only
// exported fields are parsed as part of the [Table] and
// field types are currently limited to int64, float64,
// and string; this will be expanded in future. The first
// exported field is designated the primary key, regardless
// of type. It's is recommended to use int64 for primary
// keys but not essential; SQLite works best with integers
// but the benefits aren't noticable in most use cases.
//
//	type Player struct {
//		Id int64
//		Name string
//		RoleId int64
//		role *Role // This field is ignored.
//	}
//
//	type Role struct {
//		Id int64
//		Name string
//		Strength int64
//		Stamina int64
//		Intellect int64
//		Health int64
//	}
//
//	err := db.Create(Person{}, Role{})
func (st *Storm) Create(models ...any) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, m := range models {
		e := st.createTable(m)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(m)).Wraps(e)
		}
	}

	return nil
}

func (st *Storm) createTable(model any) error {
	table, e := st.registerTable(model)
	if e != nil {
		return e
	}

	query := nidoking.Given(`
		CREATE TABLE IF NOT EXISTS {{table.GoName}} (
			{{col.GoName}} {{col.SqlType}} NOT NULL,
		  PRIMARY KEY ({{id_col.GoName}})
		)
	`).
		InlineMap("table", table).
		ListMap("col", "", table.Columns...).
		InlineMap("id_col", table.IdColumn()).
		String()

	_, e = st.db.Exec(query)
	return e
}

func (st *Storm) registerTable(model any) (Table, error) {
	table, e := st.findTableForModel(model)
	if e == nil {
		return table, nil
	}

	if !errors.Is(e, ErrNoSuchTable) {
		return Table{}, e
	}

	table, e = Parse(model)
	if e != nil {
		return Table{}, e
	}

	st.tables = append(st.tables, table)
	return table, nil
}

// Insert inserts the set of objects into the database. The
// type of each object must match a type registered via
// the [Storm.Create] function or an error is returned.
//
//	alice := Player{
//		Id: 69,
//		Name: "Alice",
//		RoleId: 5,
//	}
//
//	bob := Player{
//		Id: 42,
//		Name: "Bob",
//		RoleId: 3,
//	}
//
//	err := db.Insert(alice, bob)
func (st *Storm) Insert[T any](objects ...T) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, o := range objects {
		e := st.insertObject(o)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(o)).Wraps(e)
		}
	}

	return nil
}

func (st *Storm) insertObject(object any) error {
	table, e := st.findTableForModel(object)
	if e != nil {
		return e
	}

	query := nidoking.Given(`
		INSERT INTO {{table.GoName}} (
		  {{col.GoName}}
		)
		VALUES (
			{{q_marks}}
		)
	`).
		InlineMap("table", table).
		ListMap("col", ",", table.Columns...).
		ListRepeat("q_marks", ",", "?", table.ColumnCount()).
		String()

	values := extractColumnValues(table.Columns, object)
	_, e = st.db.Exec(query, values...)
	return e
}

// Update updates the objects within the database. Each
// object's type must match a registered type or an error
// is returned. All fields are updated except the ID
// field, which is used to determine which record to
// update.
//
//	alice := Player{
//		Id: 69,
//		Name: "Alice",
//		RoleId: 5,
//	}
//	err := db.Insert(alice)
//	// YUDO: Handle error.
//
//	alice.Name = "Alicia"
//	err = db.Update(alice)
func (st *Storm) Update[T any](objects ...T) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, obj := range objects {
		e := st.updateObject(obj)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(obj)).Wraps(e)
		}
	}

	return nil
}

func (st *Storm) updateObject(object any) error {
	table, e := st.findTableForModel(object)
	if e != nil {
		return e
	}

	query := nidoking.Given(`
			UPDATE
				{{table.GoName}}
			SET
				{{col.GoName}} = ?
			WHERE
				{{id_col.GoName}} = ?
		`).
		InlineMap("table", table).
		InlineMap("id_col", table.IdColumn()).
		ListMap("col", ",", table.Columns[1:]...).
		String()

	values := extractColumnValues(table.Columns, object)
	// Move ID value to end (for the WHERE clause)
	values = append(values[1:], values[0])

	_, e = st.db.Exec(query, values...)
	return e
}

// Select returns all records for the table associated
// with the passed model. The model's type must match a
// registered type or an error is returned.
//
//	slice, err := Select(Model{})
func (st *Storm) Select[T any](model T) ([]T, error) {
	if !st.IsOpen() {
		return nil, ErrNotOpen
	}

	err := ErrTableRequest.Fmt(typeName(model))

	table, e := st.findTableForModel(model)
	if e != nil {
		return nil, err.Wraps(e)
	}

	query := nidoking.Given(`
		SELECT
			{{col.GoName}}
		FROM
			{{table.GoName}}
	`).
		InlineMap("table", table).
		ListMap("col", ",", table.Columns...).
		String()

	rows, e := st.db.Query(query)
	if e != nil {
		return nil, err.Wraps(e)
	}

	result, e := st.scanSelectedRows[T](table, rows)
	if e != nil {
		return nil, err.Wraps(e)
	}

	return result, nil
}

func (st *Storm) scanSelectedRows[T any](
	table Table,
	rows *sql.Rows,
) ([]T, error) {
	values, valuePtrs := createValueContainers(table)
	var result []T

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrScanningRow.Fmt(i).Wraps(e)
		}

		object := constructObject[T](values)
		result = append(result, object)
	}

	e := rows.Err()
	if e != nil {
		return nil, e
	}

	return result, nil
}

func createValueContainers(table Table) ([]any, []any) {
	colCount := table.ColumnCount()
	values := make([]any, colCount, colCount)
	valuePtrs := make([]any, colCount, colCount)

	for i, col := range table.Columns {
		values[i] = col.Zero()
		valuePtrs[i] = &values[i]
	}

	return values, valuePtrs
}

func constructObject[T any](values []any) T {
	var result T

	objTyp := reflect.TypeOf(result)
	objVal := reflect.ValueOf(&result).Elem()

	valueIdx := 0
	fieldIdx := 0

	for ; fieldIdx < objVal.NumField(); fieldIdx++ {
		fieldTyp := objTyp.Field(fieldIdx)

		if !fieldTyp.IsExported() {
			continue
		}

		v := reflect.ValueOf(values[valueIdx])
		fieldVal := objVal.Field(valueIdx)
		fieldVal.Set(v)

		valueIdx++
	}

	return result
}

// SelectById returns the record with the given id from
// the table associated with the passed model. If no
// record is found then an error is returned. The model's
// type must match a registered type or an error is
// returned.
//
//	object, err := SelectById(Model{}, 123)
func (st *Storm) SelectById[T, ID any](
	model T,
	id ID,
) (T, error) {
	var empty T

	if !st.IsOpen() {
		return empty, ErrNotOpen
	}

	err := ErrTableRequest.Fmt(typeName(model), id)
	table, e := st.findTableForModel(model)
	if e != nil {
		return empty, err.Wraps(e)
	}

	e = validateIdType(table, id)
	if e != nil {
		return empty, err.Wraps(e)
	}

	query := nidoking.Given(`
		SELECT
			{{col.GoName}}
		FROM
			{{table.GoName}}
		WHERE
			{{id_col.GoName}} = ?
	`).
		ListMap("col", ",", table.Columns...).
		InlineMap("table", table).
		InlineMap("id_col", table.IdColumn()).
		String()

	rows, e := st.db.Query(query, id)
	if e != nil {
		return empty, err.Wraps(e)
	}

	result, e := st.scanSelectedRows[T](table, rows)
	if e != nil {
		return empty, err.Wraps(e)
	}

	object, ok := getFirstItemIfArray[T](result)
	if !ok {
		return empty, ErrObjectNotFound.Fmt(table.GoName, id)
	}

	return object, nil
}

func getFirstItemIfArray[T any](v any) (T, bool) {
	var empty T
	rv := reflect.ValueOf(v)

	isArray := rv.Kind() == reflect.Array
	isSlice := rv.Kind() == reflect.Slice

	if !isArray && !isSlice {
		return empty, false
	}

	if rv.Len() == 0 {
		return empty, false
	}

	return rv.Index(0).Interface().(T), true
}

// DeleteById removes the records with the given ids from
// the table associated with the passed model. If no
// record is found then nothing happens. The model's
// type must match a registered type or an error is
// returned.
//
//	e := DeleteById(Model{}, 123)
func (st *Storm) DeleteById[T, ID any](
	model T,
	ids ...ID,
) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	table, e := st.findTableForModel(model)
	if e != nil {
		return ErrTableRequest.Fmt(typeName(model)).Wraps(e)
	}

	for _, id := range ids {
		e = st.deleteById(table, id)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(model)).Wraps(e)
		}
	}

	return nil
}

func (st *Storm) deleteById(table Table, id any) error {
	e := validateIdType(table, id)
	if e != nil {
		return e
	}

	query := nidoking.Given(`
		DELETE FROM
			{{table.GoName}}
		WHERE
			{{id_col.GoName}} = ?
	`).
		InlineMap("table", table).
		InlineMap("id_col", table.IdColumn()).
		String()

	_, e = st.db.Exec(query, id)
	return e
}

// Drop removes a table from the database, deleting all
// records in the process. Passing a model for a table that
// doesn't exists does nothing.
//
//	type Player struct {
//		Id int64
//		Name string
//	}
//
//	type Role struct {
//		Id int64
//		Name string
//	}
//
//	err := db.Create(Person{}, Role{})
//	// YUDO: Handle error.
//
//	err := db.Drop(Person{}, Role{})
func (st *Storm) Drop(models ...any) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, m := range models {
		e := st.dropTable(m)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(m)).Wraps(e)
		}
	}

	return nil
}

func (st *Storm) dropTable(model any) error {
	table, e := st.findTableForModel(model)
	if errors.Is(e, ErrNoSuchTable) {
		return nil
	}

	if e != nil {
		return e
	}

	query := nidoking.Given(`
		DROP TABLE IF EXISTS {{table.GoName}}
	`).
		InlineMap("table", table).
		String()

	_, e = st.db.Exec(query)
	return e
}

func (st *Storm) findTableForModel(
	object any,
) (Table, error) {
	typ := reflect.TypeOf(object)

	for _, table := range st.tables {
		if table.GoType == typ {
			return table, nil
		}
	}

	return Table{}, ErrNoSuchTable.Fmt(typ.Name())
}

func extractColumnValues(
	columns []Column,
	object any,
) []any {
	value := reflect.ValueOf(object)
	result := make([]any, len(columns))

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		field := value.FieldByName(col.GoName)
		result[i] = field.Interface()
	}

	return result
}

func validateIdType[ID any](table Table, id ID) error {
	want := table.IdColumn().GoType
	have := reflect.TypeOf(id)

	if want != have {
		return ErrBadIdType.
			Fmt(table.GoName, want.Name(), have.Name())
	}

	return nil
}
