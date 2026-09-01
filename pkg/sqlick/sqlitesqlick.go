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

// AddStructTable satisfies the Sqlick interface.
func (ss *sqliteSqlick) AddStructTable(object any) error {
	table, e := Parse(object)

	if e != nil {
		return fmt.Errorf("Failed to add struct table: %w", e)
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

	fieldValues := ss.listOrderedFieldValues(table, value)
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

func (ss *sqliteSqlick) listOrderedFieldValues(
	table Table,
	value reflect.Value,
) []any {
	result := make([]any, len(table.Columns))

	for i := 0; i < len(table.Columns); i++ {
		col := table.Columns[i]
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
			"Could not generate SQL query for table '%s': %w",
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
