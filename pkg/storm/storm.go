package storm

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "github.com/glebarez/go-sqlite"
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

// Open opens the database. The directory path is created
// if it doesn't already exist.
//
//	err := db.Open()
//	// YUDO: Handle error.
//	defer db.Close()
func (ss *Storm) Open() error {
	parent := filepath.Dir(ss.path)
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		return fmt.Errorf(
			"Unable to check or create directory path to SQLite database: %w",
			e,
		)
	}

	db, e := sql.Open("sqlite", ss.path)
	if e != nil {
		return fmt.Errorf(
			"Unable to open SQLite database: %w",
			e,
		)
	}

	ss.db = db
	return nil
}

// IsOpen returns true if the database is open.
//
//	if db.IsOpen() {
//		// YUDO.
//	}
func (ss *Storm) IsOpen() bool {
	return ss.db != nil
}

// Close closes the database. Use with defer as usual.
//
//	err := db.Open()
//	// YUDO: Handle error.
//	defer db.Close()
func (ss *Storm) Close() error {
	if !ss.IsOpen() {
		return nil
	}

	defer func() {
		ss.db = nil
	}()

	return ss.db.Close()
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
func (ss *Storm) Create(model any) error {
	e := ss.register(model)
	if e != nil {
		return fmt.Errorf(
			"Failed to create table: %w",
			e,
		)
	}

	typ := reflect.TypeOf(model)
	table, found := ss.findTableFor(typ)
	if !found {
		return fmt.Errorf(
			"No table exists for struct '%s': %w",
			typ.Name(),
			ErrNoTableForType,
		)
	}

	e = ss.createTable(table)
	if e != nil {
		return fmt.Errorf(
			"Failed to create table '%s': %w",
			table.GoName,
			e,
		)
	}

	return nil
}

func (ss *Storm) register(model any) error {
	_, found := ss.findTableFor(reflect.TypeOf(model))
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

	ss.tables = append(ss.tables, table)
	return nil
}

func (ss *Storm) createTable(table Table) error {
	query, e := ss.generateCreateTableSql(table)
	if e != nil {
		return e
	}

	_, e = ss.db.Exec(query)
	return e
}

// TableCreate satisfies the SqlGenerator interface. The
// table name is the GoName, all columns are NOT NULL, and
// the first column is designated the PRIMARY KEY.
//
// Type mappings
//
//	int64     => INTEGER
//	float64 => REAL
//	string  => TEXT
func (ss *Storm) generateCreateTableSql(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

	s := joinLines(
		"CREATE TABLE IF NOT EXISTS %s (",
		"%s,",
		"  PRIMARY KEY (%s)",
		")",
	)

	strCols, e := generateList(tbl.Columns, ss.genColumnDef)
	if e != nil {
		return "", fmt.Errorf(
			"Failed to generate CREATE TABLE SQL for table '%s': %w",
			tbl.GoName,
			e,
		)
	}

	fb.WriteFmt(s, tbl.GoName, strCols, tbl.IdColumn().GoName)
	return fb.String(), nil
}

func (ss *Storm) genColumnDef(
	_ int,
	col Column,
) (string, error) {
	sqlType := goKindToSqliteTypeMappings[col.GoType.Kind()]

	if sqlType == "" {
		return "", fmt.Errorf(
			"Failed to generate SQL for column '%s': %w",
			col.GoName,
			ErrNoSqlType,
		)
	}

	s := fmt.Sprintf("  %s %s NOT NULL", col.GoName, sqlType)
	return s, nil
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
func (ss *Storm) Insert(object any) error {
	value := reflect.ValueOf(object)
	e := ss.insertValue(value)

	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to insert object into database: %w",
		e,
	)
}

func (ss *Storm) insertValue(
	value reflect.Value,
) error {
	table, found := ss.findTableFor(value.Type())

	if !found {
		return fmt.Errorf(
			"No table exists for struct '%s': %w",
			value.Type().Name(),
			ErrNoTableForType,
		)
	}

	fieldValues := ss.listOrderedFieldValues(
		table.Columns,
		value,
	)
	return ss.execInsert(table, fieldValues)
}

func (ss *Storm) findTableFor(
	typ reflect.Type,
) (Table, bool) {
	for _, table := range ss.tables {
		if table.GoType == typ {
			return table, true
		}
	}
	return Table{}, false
}

func (ss *Storm) findTableForOrError(
	typ reflect.Type,
) (Table, error) {
	table, found := ss.findTableFor(typ)

	if found {
		return table, nil
	}

	return Table{}, fmt.Errorf(
		"No table exists for struct '%s': %w",
		typ.Name(),
		ErrNoTableForType,
	)
}

func (ss *Storm) listOrderedFieldValues(
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

func (ss *Storm) execInsert(
	table Table,
	fieldValues []any,
) error {
	query, e := ss.generateInsertRecordSql(table, 1)
	if e != nil {
		return fmt.Errorf(
			"Could not generate insert query for table '%s': %w",
			table.GoName,
			e,
		)
	}

	_, e = ss.db.Exec(query, fieldValues...)
	if e != nil {
		return fmt.Errorf(
			"Failed to insert into table '%s': %w",
			table.GoName,
			e,
		)
	}

	return nil
}

func (ss *Storm) generateInsertRecordSql(
	tbl Table,
	rows int,
) (string, error) {
	if rows < 1 {
		return "", fmt.Errorf(
			"Failed to generate SQL for INSERT into table '%s': %w",
			tbl.GoName,
			ErrTooFewRows,
		)
	}

	fb := fmtBuilder{}

	columns, _ := generateList(tbl.Columns, ss.genColumn)
	values, _ := generateList(tbl.Columns, ss.genQuestionMark)

	s := joinLines(
		"INSERT INTO %s (",
		"%s",
		")%s",
	)

	sv := joinLines(
		" VALUES (",
		"%s",
		")",
	)
	sv = fmt.Sprintf(sv, values)
	sv = strings.Repeat(sv, rows)

	fb.WriteFmt(s, tbl.GoName, columns, sv)
	return fb.String(), nil
}

func (ss *Storm) genColumn(
	_ int,
	col Column,
) (string, error) {
	return fmt.Sprintf("  %s", col.GoName), nil
}

func (ss *Storm) genQuestionMark(
	_ int,
	col Column,
) (string, error) {
	return "  ?", nil
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
func (ss *Storm) Update(object any) error {
	value := reflect.ValueOf(object)
	e := ss.updateValue(value)

	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to update object in database: %w",
		e,
	)
}

func (ss *Storm) updateValue(
	value reflect.Value,
) error {
	table, e := ss.findTableForOrError(value.Type())
	if e != nil {
		return e
	}

	fieldValues := ss.listOrderedFieldValues(
		table.Columns,
		value,
	)

	// Move ID value to end (for the WHERE clause)
	fieldValues = append(fieldValues[1:], fieldValues[0])
	return ss.execUpdate(table, fieldValues)
}

func (ss *Storm) execUpdate(
	table Table,
	fieldValues []any,
) error {
	query, e := ss.generateUpdateRecordSql(table)
	if e != nil {
		return fmt.Errorf(
			"Could not generate update query for table '%s': %w",
			table.GoName,
			e,
		)
	}

	_, e = ss.db.Exec(query, fieldValues...)
	if e != nil {
		return fmt.Errorf(
			"Failed to update row in table '%s': %w",
			table.GoName,
			e,
		)
	}

	return nil
}

func (ss *Storm) generateUpdateRecordSql(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

	nonIdColumns := tbl.NonIdColumns()
	setters, _ := generateList(
		nonIdColumns,
		ss.genColumnSetter,
	)

	s := joinLines(
		"UPDATE",
		"  %s",
		"SET",
		"%s",
		"WHERE",
		"  %s = ?",
	)

	fb.WriteFmt(s, tbl.GoName, setters, tbl.IdColumn().GoName)
	return fb.String(), nil
}

func (ss *Storm) genColumnSetter(
	_ int,
	col Column,
) (string, error) {
	return fmt.Sprintf("  %s = ?", col.GoName), nil
}

// SelectAll returns all records for the table associated
// with the passed model. The model's type must be
// registered (Register()) and table created
// (CreateTables()) for the select to return without
// error.
//
//	slice, err := SelectAll(Model{})
func (ss *Storm) SelectAll(
	model any,
) (any, error) {
	typ := reflect.TypeOf(model)
	result, e := ss.selectAllOfType(typ)

	if e == nil {
		return result, nil
	}

	return nil, fmt.Errorf(
		"Failed to select all from database: %w",
		e,
	)
}

func (ss *Storm) selectAllOfType(
	typ reflect.Type,
) (any, error) {
	table, e := ss.findTableForOrError(typ)
	if e != nil {
		return nil, e
	}

	fieldValues := ss.listOrderedFieldValues(
		table.Columns,
		reflect.Zero(typ),
	)

	return ss.querySelectAll(table, fieldValues)
}

func (ss *Storm) querySelectAll(
	table Table,
	fieldValues []any,
) (any, error) {
	query, e := ss.generateSelectAllRecordsSql(table)
	if e != nil {
		return nil, fmt.Errorf(
			"Could not generate select all query for table '%s': %w",
			table.GoName,
			e,
		)
	}

	rows, e := ss.db.Query(query, fieldValues...)
	if e != nil {
		return nil, fmt.Errorf(
			"Failed to select all from table '%s': %w",
			table.GoName,
			e,
		)
	}

	result, e := ss.scanSelectedRows(table, rows)
	if e != nil {
		return nil, e
	}

	return result, nil
}

func (ss *Storm) generateSelectAllRecordsSql(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

	columns, _ := generateList(tbl.Columns, ss.genColumn)

	s := joinLines(
		"SELECT",
		"%s",
		"FROM",
		"  %s",
	)

	fb.WriteFmt(s, columns, tbl.GoName)
	return fb.String(), nil
}

func (ss *Storm) scanSelectedRows(
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
	values := make([]any, tbl.NumColumn())
	valuePtrs := make([]any, tbl.NumColumn())

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
func (ss *Storm) SelectById(
	model any,
	id any,
) (any, error) {
	typ := reflect.TypeOf(model)
	result, e := ss.selectByIdOfType(typ, id)
	if e == nil {
		return result, nil
	}

	return nil, fmt.Errorf(
		"Failed to select by ID from database: %w",
		e,
	)
}

func (ss *Storm) selectByIdOfType(
	typ reflect.Type,
	id any,
) (any, error) {
	tbl, e := ss.findTableForOrError(typ)
	if e != nil {
		return nil, e
	}

	e = ss.validateModelIdType(tbl, id)
	if e != nil {
		return nil, e
	}

	return ss.querySelectById(tbl, id)
}

func (ss *Storm) validateModelIdType(
	tbl Table,
	id any,
) error {
	idTyp := reflect.TypeOf(id)

	if tbl.IdColumn().GoType.Kind() != idTyp.Kind() {
		return ErrBadIdType
	}

	return nil
}

func (ss *Storm) querySelectById(
	tbl Table,
	id any,
) (any, error) {
	query, e := ss.generateSelectRecordByIdSql(tbl)
	if e != nil {
		return nil, fmt.Errorf(
			"Could not generate select by ID query for table '%s': %w",
			tbl.GoName,
			e,
		)
	}

	rows, e := ss.db.Query(query, id)
	if e != nil {
		return nil, fmt.Errorf(
			"Failed to select by ID from table '%s': %w",
			tbl.GoName,
			e,
		)
	}

	result, e := ss.scanSelectedRows(tbl, rows)
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

func (ss *Storm) generateSelectRecordByIdSql(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

	columns, _ := generateList(tbl.Columns, ss.genColumn)

	s := joinLines(
		"SELECT",
		"%s",
		"FROM",
		"  %s",
		"WHERE",
		"  %s = ?",
	)

	fb.WriteFmt(s, columns, tbl.GoName, tbl.IdColumn().GoName)
	return fb.String(), nil
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
func (ss *Storm) DeleteById(
	model any,
	id any,
) error {
	typ := reflect.TypeOf(model)
	e := ss.deleteByIdOfType(typ, id)
	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Failed to delete by ID from database: %w",
		e,
	)
}
func (ss *Storm) deleteByIdOfType(
	typ reflect.Type,
	id any,
) error {
	tbl, e := ss.findTableForOrError(typ)
	if e != nil {
		return e
	}

	e = ss.validateModelIdType(tbl, id)
	if e != nil {
		return e
	}

	return ss.execDeleteById(tbl, id)
}

func (ss *Storm) execDeleteById(
	tbl Table,
	id any,
) error {
	query, e := ss.generateDeleteRecordByIdSql(tbl)
	if e != nil {
		return fmt.Errorf(
			"Could not generate delete by ID '%v' query for table '%s': %w",
			id,
			tbl.GoName,
			e,
		)
	}

	_, e = ss.db.Exec(query, id)
	if e != nil {
		return fmt.Errorf(
			"Failed to delete by ID '%v' from table '%s': %w",
			id,
			tbl.GoName,
			e,
		)
	}

	return nil
}

func (ss *Storm) generateDeleteRecordByIdSql(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

	s := joinLines(
		"DELETE FROM",
		"  %s",
		"WHERE",
		"  %s = ?",
	)

	fb.WriteFmt(s, tbl.GoName, tbl.IdColumn().GoName)
	return fb.String(), nil
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
func (ss *Storm) Drop(model any) error {
	typ := reflect.TypeOf(model)
	table, found := ss.findTableFor(typ)
	if !found {
		return nil
	}

	e := ss.dropTable(table)
	if e != nil {
		return fmt.Errorf(
			"Failed to drop table '%s': %w",
			table.GoName,
			e,
		)
	}

	return nil
}

func (ss *Storm) dropTable(table Table) error {
	query := fmt.Sprintf(
		"DROP TABLE IF EXISTS %s",
		table.GoName,
	)
	_, e := ss.db.Exec(query, table.GoName)
	return e
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
