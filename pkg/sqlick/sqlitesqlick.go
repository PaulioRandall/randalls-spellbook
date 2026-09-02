package sqlick

import (
	"database/sql"
	"fmt"
	"reflect"

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

/*
// SelectAll satisfies the Sqlick interface.
func (ss *sqliteSqlick) SelectAll(
	object any,
) (any, error) {
	value := reflect.ValueOf(object)
	result, e := ss.selectAllOfValue(value)

	if e == nil {
		return result, nil
	}

	return nil, fmt.Errorf(
		"Failed to select all from database: %w",
		e,
	)
}

func (ss *sqliteSqlick) selectAllOfValue(
	value reflect.Value,
) (any, error) {
	table, e := ss.findTableForOrError(value.Type())
	if e != nil {
		return nil, e
	}

	fieldValues := ss.listOrderedFieldValues(
		table.Columns,
		value,
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
	colNames, e := rows.Columns()
	if e != nil {
		return nil, e
	}

	valuePtrs := createValueContainers(
		table,
		colNames,
	)

	var result []any

	for rows.Next() {
		rows.Scan(valuePtrs...)

		objectValue := reflect.Zero(table.GoType)
		populateObjectWithValues(
			colNames,
			valuePtrs,
			&objectValue,
		)

		object := objectValue.Interface()
		result = append(result, object)
	}

	return result, nil
}

func createValueContainers(
	table Table,
	colNames []string,
) []any {
	colCount := len(colNames)
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

	for i, name := range colNames {
		// TODO: Find column by name and create a zero value
		//       using the type. Will require changing GoType
		//       to reflect.Type.
		valuePtrs[i] = &values[i]
	}

	return valuePtrs
}

func populateObjectWithValues(
	colNames []string,
	valuePtrs []any,
	objectValue *reflect.Value,
) {
	for i, name := range colNames {
		field := objectValue.FieldByName(name)
		value := valuePtrs[i].(*any)
		field.Set(reflect.ValueOf(*value))
	}
}
*/
