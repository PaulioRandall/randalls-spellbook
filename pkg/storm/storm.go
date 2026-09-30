package storm

import (
	"database/sql"
	"reflect"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sprints"
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/mapper"
)

// TODO: Make Storm thread safe. Lock on function entry
//       and defer the unlock.
// TODO: Check the ordering of fields when scanning
//       database rows. Ensure the order the rows appear in
//       the select statement is the order they are
//       scanned to avoid the wrong info appearing in the
//       wrong struct fields.
// TODO: API call cache, create a ModelTable cache that
//       lives and dies in a single API call so bulk
//       inserts of the same kind don't fetch metadata we
//       already know will be the same.

var (
	// ErrNotOpen occurs when trying to perform an operation
	// before opening the database.
	ErrNotOpen = sin.Template(
		"Database not open",
	)

	// ErrTableRequest occurs within an error chain when any
	// error occurs involving a specific table/struct, except
	// for 'not found' errors.
	ErrTableRequest = sin.Template(
		"Request failed for table '%s'",
	)

	// ErrScanningRow is returned when an error occurs
	// scanning database results.
	ErrScanningRow = sin.Template(
		"Scanning row '%d'",
	)

	// ErrObjectNotFound is returned when a search for a
	// specific object/row failed.
	ErrObjectNotFound = sin.Template(
		"Object '%s' not found with ID '%v'",
	)

	// ErrDatabaseFile is returned when an error occurs with
	// or while opening or closing the database.
	ErrDatabaseFile = sin.Template(
		"Database IO error '%s'",
	)

	// ErrNoSuchTable is returned when an object is passed
	// to a function which does not have a registered table
	// for its type.
	ErrNoSuchTable = sin.Template(
		"No matching table for object type '%s'",
	)

	// ErrBadIdType is returned when an ID passed to a
	// function, e.g. SelectById, is not of the same type as
	// the ID field of the associated model type. This may
	// be returned even for compatible types like int when
	// int64 is expected.
	ErrBadIdType = sin.Template(
		"ID type mismatch for '%s', got %s, want %s",
	)
)

// Storm is the core type and the interface to the SQLite
// database.
type Storm struct {
	path string
	db   *sql.DB
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
		return sin.Err("Unable to open SQLite database").
			Wrap(e).
			WrapIn(ErrDatabaseFile).
			Fmt(st.path)
	}

	st.db = db
	return nil
}

// IsOpen returns true if the database is open.
func (st *Storm) IsOpen() bool {
	return st.db != nil
}

// Close closes the database. Use with defer as usual.
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

	return sin.Err("Unable to close SQLite database").
		Wrap(e).
		WrapIn(ErrDatabaseFile).
		Fmt(st.path)
}

// Table returns the full table details the passed model
// represents. All columns in the table are included, not
// just those that map to the passed model type.
func (st *Storm) Table(model any) (scumble.SqlTable, error) {
	return scumble.QueryTable(
		st.db,
		reflect.TypeOf(model).Name(),
	)
}

// Create creates tables, represented by the passed models,
// within the database. If a table already exists then
// the model is ignored.
//
// It's safe to call Create at anytime, but doing it all
// upfront is also fine. SQLite works best with integers as
// primary keys but the benefits aren't noticable in most
// use cases.
//
// # Table creation rules
//
//   - Model (struct) name becomes the table name.
//   - The exported model fields become columns.
//   - Field name is the column name.
//   - Field type is mapped to a SQLite column type.
//   - Only primitive types may be used as field types.
//   - All columns have NOT NULL constraint.
//   - All columns have DEFAULT set to the zero value of
//     the field type.
//   - The first field in the model is designated the
//     PRIMARY KEY.
//
// # Type mappings
//
//	INTEGER:
//		int, int8, int16, int32, int64,
//		uint, uint8, uint16, uint32, uint64
//	REAL:
//		float32, float64
//	TEXT:
//		string
func (st *Storm) Create(models ...any) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, m := range models {
		e := st.createTableFromModel(m)
		if e != nil {
			return e
		}
	}

	return nil
}

func (st *Storm) createTableFromModel(model any) (e error) {
	table, exists, e := st.getTable(model)
	if e != nil {
		goto Err
	}

	if exists {
		return nil
	}

	e = st.createTable(table)
	if e != nil {
		goto Err
	}

	return nil

Err:
	return ErrTableRequest.Fmt(typeName(model)).Wrap(e)
}

func (st *Storm) createTable(table mapper.ModelTable) error {
	query := nidoking.Given(`
		CREATE TABLE IF NOT EXISTS {{table.SqlName}} (
			{{col.SqlName}} {{col.SqlType}} NOT NULL DEFAULT {{col.SqlDefault}},
		  PRIMARY KEY ({{pk_col.SqlName}})
		)
	`).
		InlineMap("table", table).
		ListMap("col", "", table.Columns...).
		InlineMap("pk_col", table.PrimaryKeyColumn()).
		String()

	_, e := st.db.Exec(query)
	return e
}

// Drop removes a table from the database. If the target
// table doesn't exist then nothing happens and no error is
// returned. Related model cache entries are also removed.
//
// All data is deleted in the process and there's no way to
// restore it. To protect data, create backups of the
// database file.
func (st *Storm) Drop(models ...any) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, m := range models {
		e := st.dropTable(m)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(m)).Wrap(e)
		}
	}

	return nil
}

func (st *Storm) dropTable(model any) error {
	table, found, e := st.getTable(model)
	if e != nil || !found {
		return e
	}

	if !found {
		return nil
	}

	query := nidoking.Given(`
		DROP TABLE IF EXISTS {{table.SqlName}}
	`).
		InlineMap("table", table).
		String()

	_, e = st.db.Exec(query)
	return e
}

// Insert inserts all passed objects into the database. If
// a table doesn't exist for an object, it will be created
// as if passed to [Storm.Create].
func (st *Storm) Insert[T any](objects ...T) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, o := range objects {
		e := st.insertObject(o)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(o)).Wrap(e)
		}
	}

	return nil
}

func (st *Storm) insertObject(object any) error {
	table, e := st.getOrCreateTable(object)
	if e != nil {
		return e
	}

	query := nidoking.Given(`
		INSERT INTO {{table.SqlName}} (
		  {{col.SqlName}}
		)
		VALUES (
			{{q_marks}}
		)
	`).
		InlineMap("table", table).
		ListMap("col", ",", table.Columns...).
		ListRepeat("q_marks", ",", "?", len(table.Columns)).
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
			return ErrTableRequest.Fmt(typeName(obj)).Wrap(e)
		}
	}

	return nil
}

func (st *Storm) updateObject(object any) error {
	table, found, e := st.getTable(object)
	if e != nil {
		return e
	}

	if !found {
		return ErrNoSuchTable.Fmt(typeName(object))
	}

	pkCol := table.PrimaryKeyColumn()
	nonPkCols := table.NonPrimaryKeyColumns()

	query := nidoking.Given(`
			UPDATE
				{{table.SqlName}}
			SET
				{{non_pk_col.SqlName}} = ?
			WHERE
				{{pk_col.SqlName}} = ?
		`).
		InlineMap("table", table).
		ListMap("non_pk_col", ",", nonPkCols...).
		InlineMap("pk_col", pkCol).
		String()

	// Because PK is the last SQL parameter.
	cols := append(nonPkCols, pkCol)
	values := extractColumnValues(cols, object)
	_, e = st.db.Exec(query, values...)
	return e
}

// Select returns all records for the table associated
// with the passed model.
func (st *Storm) Select[T any](model T) (result []T, e error) {
	var query string
	var rows *sql.Rows

	if !st.IsOpen() {
		return nil, ErrNotOpen
	}

	table, found, e := st.getTable(model)
	if e != nil {
		return nil, e
	}

	sprints.Println(table)
	for _, col := range table.Columns {
		sprints.Println(col)
	}

	if !found {
		return nil, nil
	}

	query = nidoking.Given(`
		SELECT
			{{col.SqlName}}
		FROM
			{{table.SqlName}}
	`).
		InlineMap("table", table).
		ListMap("col", ",", table.Columns...).
		String()

	rows, e = st.db.Query(query)
	if e != nil {
		goto Err
	}

	result, e = st.scanSelectedRows[T](table, rows)
	if e != nil {
		goto Err
	}

	return result, nil

Err:
	return nil, ErrTableRequest.Fmt(typeName(model)).Wrap(e)
}

func (st *Storm) scanSelectedRows[T any](
	table mapper.ModelTable,
	rows *sql.Rows,
) ([]T, error) {
	values, valuePtrs := createValueContainers(table)
	var result []T

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrScanningRow.Fmt(i).Wrap(e)
		}

		object := constructObject[T](table, values)
		result = append(result, object)
	}

	e := rows.Err()
	if e != nil {
		return nil, e
	}

	return result, nil
}

func createValueContainers(table mapper.ModelTable) ([]any, []any) {
	colCount := len(table.Columns)
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

	for i, col := range table.Columns {
		values[i] = col.New[any]()
		valuePtrs[i] = &values[i]
	}

	return values, valuePtrs
}

func constructObject[T any](
	table mapper.ModelTable,
	values []any,
) T {
	var result T

	objVal := reflect.ValueOf(&result).Elem()

	for valueIdx, col := range table.Columns {
		v := reflect.ValueOf(values[valueIdx])
		fieldVal := objVal.Field(col.GoIndex)

		if v.CanConvert(fieldVal.Type()) {
			v = v.Convert(fieldVal.Type())
		} else {
			// TODO: Rethink and tidy
			panic("Can't convert from " + v.Type().Name() + " to " + fieldVal.Type().Name())
		}

		fieldVal.Set(v)
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
) (result T, e error) {
	var empty T
	var query string
	var rows *sql.Rows
	var resultSet []T
	var ok bool

	if !st.IsOpen() {
		return empty, ErrNotOpen
	}

	table, found, e := st.getTable(model)
	if e != nil {
		goto Err
	}

	if !found {
		e = ErrObjectNotFound.Fmt(table.GoName, id)
		goto Err
	}

	query = nidoking.Given(`
		SELECT
			{{col.SqlName}}
		FROM
			{{table.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		ListMap("col", ",", table.Columns...).
		InlineMap("table", table).
		InlineMap("pk_col", table.PrimaryKeyColumn()).
		String()

	rows, e = st.db.Query(query, id)
	if e != nil {
		goto Err
	}

	resultSet, e = st.scanSelectedRows[T](table, rows)
	if e != nil {
		goto Err
	}

	result, ok = getFirstItemIfArray[T](resultSet)
	if !ok {
		e = ErrObjectNotFound.Fmt(table.GoName, id)
		goto Err
	}

	return result, nil

Err:
	return empty, sin.Fmt("For object with ID '%v'", id).
		Wrap(e).
		WrapIn(ErrTableRequest).
		Fmt(typeName(model))
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
func (st *Storm) DeleteById[T, ID any](
	model T,
	ids ...ID,
) (e error) {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, id := range ids {
		e = st.deleteById(model, id)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(model)).Wrap(e)
		}
	}

	return nil
}

func (st *Storm) deleteById(model any, id any) error {
	table, found, e := st.getTable(model)
	if e != nil {
		return e
	}

	if !found {
		return nil
	}

	query := nidoking.Given(`
		DELETE FROM
			{{table.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		InlineMap("table", table).
		InlineMap("pk_col", table.PrimaryKeyColumn()).
		String()

	_, e = st.db.Exec(query, id)
	return e
}

func (st *Storm) getTable(model any) (mapper.ModelTable, bool, error) {
	return mapper.MapModel(st.db, model)
}

func (st *Storm) getOrCreateTable(model any) (mapper.ModelTable, error) {
	var zero mapper.ModelTable

	table, exists, e := mapper.MapModel(st.db, model)
	if e != nil {
		return zero, e
	}

	if exists {
		return table, nil
	}

	e = st.createTable(table)
	if e != nil {
		return zero, e
	}

	return table, nil
}

func extractColumnValues(
	columns []mapper.ModelColumn,
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
