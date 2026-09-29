package storm

import (
	"database/sql"
	"errors"
	"reflect"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/mapper"
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
)

// TODO: Make Storm thread safe. Lock on function entry
//       and defer the unlock.
// TODO: Check the ordering of fields when scanning
//       database rows. Ensure the order the rows appear in
//       the select statement is the order they are
//       scanned to avoid the wrong info appearing in the
//       wrong struct fields.

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
	path   string
	tables []Table
	db     *sql.DB
	cache  modelCache
}

// New returns a new [Storm] for the database represented
// by path.
func New(path string) *Storm {
	return &Storm{
		path:  path,
		cache: modelCache{},
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
		return sin.Err("Unable to open SQLite database").
			Wrap(e).
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

	return sin.Err("Unable to close SQLite database").
		Wrap(e).
		WrapIn(ErrDatabaseFile).
		Fmt(st.path)
}

// Table returns the full table details the passed model
// represents. All columns in the table are included, not
// just those that map to the passed model type.
func (st *Storm) Table(model any) (schema.SqlTable, error) {
	return schema.QueryTable(
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
	table, exists, e := mapper.MapModel(st.db, model)
	if e != nil {
		goto EnhanceError
	}

	if exists {
		return nil
	}

	st.cache.set(model, table)

	e = st.createTable(table)
	if e != nil {
		goto EnhanceError
	}

	// TEMP
	_, e = st.registerTable(model)
	if e != nil {
		goto EnhanceError
	}

	return nil

EnhanceError:
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

// TODO: Remove
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

func (st *Storm) dropTable(model any) (e error) {
	table, found := st.cache.get(model)
	if found {
		goto DoDrop
	}

	table, found, e = mapper.MapModel(st.db, model)
	if e != nil || !found {
		return e
	}

DoDrop:
	st.cache.clearTable(table.SqlName)

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

	values := extractColumnValues2(table.Columns, object)
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
		return nil, err.Wrap(e)
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
		return nil, err.Wrap(e)
	}

	result, e := st.scanSelectedRows[T](table, rows)
	if e != nil {
		return nil, err.Wrap(e)
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
			return nil, ErrScanningRow.Fmt(i).Wrap(e)
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
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

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
		return empty, err.Wrap(e)
	}

	e = validateIdType(table, id)
	if e != nil {
		return empty, err.Wrap(e)
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
		return empty, err.Wrap(e)
	}

	result, e := st.scanSelectedRows[T](table, rows)
	if e != nil {
		return empty, err.Wrap(e)
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
		return ErrTableRequest.Fmt(typeName(model)).Wrap(e)
	}

	for _, id := range ids {
		e = st.deleteById(table, id)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(model)).Wrap(e)
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

func (st *Storm) getOrCreateTable(model any) (mapper.ModelTable, error) {
	var zero mapper.ModelTable

	cachedTable, found := st.cache.get(model)
	if found {
		return cachedTable, nil
	}

	table, exists, e := mapper.MapModel(st.db, model)
	if e != nil {
		return zero, e
	}
	st.cache.set(model, table)

	if exists {
		return table, nil
	}

	e = st.createTable(table)
	if e != nil {
		return zero, e
	}

	return table, nil
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

func extractColumnValues2(
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

func validateIdType[ID any](table Table, id ID) error {
	want := table.IdColumn().GoType
	have := reflect.TypeOf(id)

	if want != have {
		return ErrBadIdType.
			Fmt(table.GoName, want.Name(), have.Name())
	}

	return nil
}
