package sqlick

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

type sqliteSqlick struct {
	path   string
	gen    SqlGenerator
	tables []Table
	db     *sql.DB
}

func NewSqliteDB(path string) *sqliteSqlick {
	return &sqliteSqlick{
		path: path,
		gen:  sqliteGenerator{},
	}
}

// Open satisfies the Sqlick interface.
func (ss *sqliteSqlick) Open() error {
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

// IsOpen satisfies the Sqlick interface.
func (ss *sqliteSqlick) IsOpen() bool {
	return ss.db != nil
}

// Close satisfies the Sqlick interface.
func (ss *sqliteSqlick) Close() error {
	if !ss.IsOpen() {
		return nil
	}

	defer func() {
		ss.db = nil
	}()

	return ss.db.Close()
}

// Register satisfies the Sqlick interface.
func (ss *sqliteSqlick) Register(object any) error {
	_, found := ss.findTableFor(reflect.TypeOf(object))
	if found {
		return nil
	}

	table, e := Parse(object)

	if e != nil {
		return fmt.Errorf(
			"Failed to register struct/table: %w",
			e,
		)
	}

	ss.tables = append(ss.tables, table)
	return nil
}

// CreateTables satisfies the Sqlick interface.
func (ss *sqliteSqlick) CreateTables() error {
	for _, table := range ss.tables {
		e := ss.createTable(table)
		if e == nil {
			continue
		}

		return fmt.Errorf(
			"Failed to create table '%s': %w",
			table.GoName,
			e,
		)
	}

	return nil
}

func (ss *sqliteSqlick) createTable(table Table) error {
	query, e := ss.gen.TableCreate(table)
	if e != nil {
		return e
	}

	_, e = ss.db.Exec(query)
	return e
}

// Insert satisfies the Sqlick interface.
func (ss *sqliteSqlick) Insert(object any) error {
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

func (ss *sqliteSqlick) insertValue(
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

func (ss *sqliteSqlick) findTableFor(
	typ reflect.Type,
) (Table, bool) {
	for _, table := range ss.tables {
		if table.GoType == typ {
			return table, true
		}
	}
	return Table{}, false
}

func (ss *sqliteSqlick) findTableForOrError(
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

func (ss *sqliteSqlick) listOrderedFieldValues(
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

func (ss *sqliteSqlick) execInsert(
	table Table,
	fieldValues []any,
) error {
	query, e := ss.gen.TableInsert(table, 1)
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

// Update satisfies the Sqlick interface.
func (ss *sqliteSqlick) Update(object any) error {
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

func (ss *sqliteSqlick) updateValue(
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

func (ss *sqliteSqlick) execUpdate(
	table Table,
	fieldValues []any,
) error {
	query, e := ss.gen.TableUpdateById(table)
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

// SelectAll satisfies the Sqlick interface.
func (ss *sqliteSqlick) SelectAll(
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

func (ss *sqliteSqlick) selectAllOfType(
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

func (ss *sqliteSqlick) querySelectAll(
	table Table,
	fieldValues []any,
) (any, error) {
	query, e := ss.gen.TableSelectAll(table)
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

func (ss *sqliteSqlick) scanSelectedRows(
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

// SelectById satisfies the Sqlick interface.
func (ss *sqliteSqlick) SelectById(
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

func (ss *sqliteSqlick) selectByIdOfType(
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

func (ss *sqliteSqlick) validateModelIdType(
	tbl Table,
	id any,
) error {
	idTyp := reflect.TypeOf(id)

	if tbl.IdColumn().GoType.Kind() != idTyp.Kind() {
		return ErrBadIdType
	}

	return nil
}

func (ss *sqliteSqlick) querySelectById(
	tbl Table,
	id any,
) (any, error) {
	query, e := ss.gen.TableSelectById(tbl)
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
