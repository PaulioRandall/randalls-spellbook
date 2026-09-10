package storm

import (
	"database/sql"
	"reflect"

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
	e := makeParentDirs(st.path)
	if e != nil {
		return e
	}

	db, e := sql.Open("sqlite", st.path)
	if e != nil {
		return stormy("Unable to open SQLite database").
			Wrap(e)
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
func (st *Storm) Create[T any](model T) error {
	table, e := st.registerTable(model)
	if e != nil {
		return ErrCreatingTable.Wrap(e)
	}

	e = st.createTable(table)
	if e != nil {
		return ErrCreatingTable.Wrap(e)
	}

	return nil
}

func (st *Storm) registerTable(model any) (Table, error) {
	table, e := st.findTableForModel(model)
	if e == nil {
		return table, nil
	}

	table, e = Parse(model)
	if e != nil {
		return Table{},
			stormy("Failed to register struct/table").Wrap(e)
	}

	st.tables = append(st.tables, table)
	return table, nil
}

func (st *Storm) createTable(table Table) error {
	query := nidoking.Given(`
		CREATE TABLE IF NOT EXISTS {{table}} (
			{{columns}},
		  PRIMARY KEY ({{id_column}})
		)
	`).
		Fmt("table", table.GoName).
		Join("columns", "", table.ColumnNames()...).
		Fmt("id_column", table.IdColumn().GoName).
		String()

	_, e := st.db.Exec(query)
	if e != nil {
		return ErrExeSqlCreate.
			Table(table.GoName).
			Wrap(e)
	}

	return nil
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
func (st *Storm) Insert[T any](object T) error {
	table, e := st.findTableForModel(object)
	if e != nil {
		return ErrInsertingObject.Wrap(e)
	}

	values := extractColumnValues(table.Columns, object)
	e = st.execInsert(table, values)
	if e != nil {
		return ErrInsertingObject.Wrap(e)
	}

	return nil
}

func (st *Storm) execInsert(
	table Table,
	values []any,
) error {
	query := nidoking.Given(`
		INSERT INTO {{table}} (
		  {{columns}}
		)
		VALUES (
			{{q_marks}}
		)
	`).
		Fmt("table", table.GoName).
		Join("columns", ",", table.ColumnNames()...).
		Repeat("q_marks", ",", "?", table.ColumnCount()).
		String()

	_, e := st.db.Exec(query, values...)
	if e != nil {
		return ErrExeSqlInsert.
			Table(table.GoName).
			Wrap(e)
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
func (st *Storm) Update[T any](object T) error {
	table, e := st.findTableForModel(object)
	if e != nil {
		return ErrUpdatingObject.Wrap(e)
	}

	values := extractColumnValues(table.Columns, object)
	// Move ID value to end (for the WHERE clause)
	values = append(values[1:], values[0])

	e = st.execUpdate(table, values)
	if e != nil {
		return ErrUpdatingObject.Wrap(e)
	}

	return nil
}

func (st *Storm) execUpdate(
	table Table,
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
		Fmt("table", table.GoName).
		Join("columns", ",", table.ColumnNames()[1:]...).
		Fmt("id_column", table.IdColumn().GoName).
		String()

	_, e := st.db.Exec(query, fieldValues...)
	if e != nil {
		return ErrExeSqlUpdate.Table(table.GoName).Wrap(e)
	}

	return nil
}

// SelectAll returns all records for the table associated
// with the passed model. The model's type must be
// registered (Register()) and table created
// (CreateTables()) for the select to return without
// error.
//
//	slice, err := SelectAll(Model{})
func (st *Storm) SelectAll[T any](model T) ([]T, error) {
	var empty []T

	table, e := st.findTableForModel(model)
	if e != nil {
		return empty, ErrSelectingAllObjects.Wrap(e)
	}

	results, e := st.querySelectAll[T](table)
	if e != nil {
		return empty, ErrSelectingAllObjects.Wrap(e)
	}

	return results, nil
}

func (st *Storm) querySelectAll[T any](
	table Table,
) ([]T, error) {
	query := nidoking.Given(`
		SELECT
			{{columns}}
		FROM
			{{table}}
	`).
		Join("columns", ",", table.ColumnNames()...).
		Fmt("table", table.GoName).
		String()

	rows, e := st.db.Query(query)
	if e != nil {
		return nil, ErrExeSqlSelect.Table(table.GoName).Wrap(e)
	}

	result, e := st.scanSelectedRows[T](table, rows)
	if e != nil {
		return nil, ErrScanningRows.Table(table.GoName).Wrap(e)
	}

	return result, nil
}

func (st *Storm) scanSelectedRows[T any](
	table Table,
	rows *sql.Rows,
) ([]T, error) {
	values, valuePtrs := createValueContainers(table)
	var result []T

	for rowIdx := 0; rows.Next(); rowIdx++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrScanningRow.Row(rowIdx).Wrap(e)
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

	for fieldIdx < objVal.NumField() {
		fieldTyp := objTyp.Field(fieldIdx)

		if !fieldTyp.IsExported() {
			continue
		}

		v := reflect.ValueOf(values[valueIdx])
		fieldVal := objVal.Field(valueIdx)
		fieldVal.Set(v)

		fieldIdx++
		valueIdx++
	}

	return result
}

// SelectById returns the record with the given id from
// the table associated with the passed model. If no
// record is found then an error is returned. The
// model's type must be registered (Register()) and table
// created (CreateTables()) for the select to return
// without error.
//
//	slice, err := SelectById(Model{}, 123)
func (st *Storm) SelectById[T, ID any](
	model T,
	id ID,
) (T, error) {
	var empty T

	table, e := st.findTableForModel(model)
	if e != nil {
		tableName := reflect.TypeOf(model).Name()
		return empty, ErrSelectingObjectsById.
			Table(tableName).Id(id).Wrap(e)
	}

	e = validateIdType(table, id)
	if e != nil {
		return empty, ErrSelectingObjectsById.
			Table(table.GoName).Id(id).Wrap(e)
	}

	result, e := st.querySelectById[T](table, id)
	if e != nil {
		return empty, ErrSelectingObjectsById.
			Table(table.GoName).Id(id).Wrap(e)
	}

	return result, nil
}

func (st *Storm) querySelectById[T any](
	table Table,
	id any,
) (T, error) {
	var empty T

	query := nidoking.Given(`
		SELECT
			{{columns}}
		FROM
			{{table}}
		WHERE
			{{id_column}} = ?
	`).
		Join("columns", ",", table.ColumnNames()...).
		Fmt("table", table.GoName).
		Fmt("id_column", table.IdColumn().GoName).
		String()

	rows, e := st.db.Query(query, id)
	if e != nil {
		return empty, ErrExeSqlSelect.Table(table.GoName).Wrap(e)
	}

	result, e := st.scanSelectedRows[T](table, rows)
	if e != nil {
		return empty, ErrScanningRows.Table(table.GoName).Wrap(e)
	}

	object, ok := getFirstItemIfArray[T](result)
	if !ok {
		return empty, ErrObjectNotFound.Table(table.GoName)
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

// DeleteById removes the record with the given id from
// the table associated with the passed model. If no
// record is found then nothing happens.
//
//	e := DeleteById(Model{}, 123)
func (st *Storm) DeleteById[T, ID any](
	model T,
	id ID,
) error {
	table, e := st.findTableForModel(model)
	if e != nil {
		tableName := reflect.TypeOf(model).Name()
		return ErrDeletingObject.
			Table(tableName).Id(id).Wrap(e)
	}

	e = validateIdType(table, id)
	if e != nil {
		return ErrDeletingObject.
			Table(table.GoName).Id(id).Wrap(e)
	}

	e = st.execDeleteById(table, id)
	if e != nil {
		return ErrDeletingObject.
			Table(table.GoName).Id(id).Wrap(e)
	}

	return nil
}

func (st *Storm) execDeleteById(
	table Table,
	id any,
) error {
	query := nidoking.Given(`
		DELETE FROM
			{{table}}
		WHERE
			{{id_column}} = ?
	`).
		Fmt("table", table.GoName).
		Fmt("id_column", table.IdColumn().GoName).
		String()

	_, e := st.db.Exec(query, id)
	if e != nil {
		return ErrExeSqlDelete.
			Table(table.GoName).Id(id).Wrap(e)
	}

	return nil
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
func (st *Storm) Drop[T any](model T) error {
	typ := reflect.TypeOf(model)
	table, found := st.findTableForType(typ)
	if !found {
		return nil
	}

	e := st.execDropQuery(table)
	if e != nil {
		return ErrDroppingTable.Table(table.GoName).Wrap(e)
	}

	return nil
}

func (st *Storm) execDropQuery(table Table) error {
	query := nidoking.Given(`
		DROP TABLE IF EXISTS {{table}}
	`).
		Fmt("table", table.GoName).
		String()

	_, e := st.db.Exec(query)
	if e != nil {
		return ErrExeSqlDrop.Table(table.GoName).Wrap(e)
	}

	return nil
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

	return Table{}, ErrNoSuchTable.Table(typ.Name())
}

func (st *Storm) findTableForType(
	typ reflect.Type,
) (Table, bool) {
	for _, table := range st.tables {
		if table.GoType == typ {
			return table, true
		}
	}
	return Table{}, false
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
	if table.IdColumn().GoType != reflect.TypeOf(id) {
		return ErrBadIdType.Table(table.GoName).Id(id)
	}
	return nil
}

func isArrayOrSlice(model any) bool {
	k := reflect.TypeOf(model).Kind()
	return k == reflect.Slice || k == reflect.Array
}
