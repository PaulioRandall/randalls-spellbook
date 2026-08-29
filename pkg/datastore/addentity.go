package datastore

import (
	"fmt"
	"reflect"
	"unicode"
)

func (ds *SqliteDatastore) AddEntity(object any) error {
	table, e := parseSqlitePropsFromObject(object)
	if e != nil {
		return e
	}

	ds.tables = append(ds.tables, table)
	return nil
}

func parseSqlitePropsFromObject(
	object any,
) (DbTable, error) {
	empty := DbTable{}
	t := reflect.TypeOf(object)

	if t.Kind() != reflect.Struct {
		return empty, fmt.Errorf("Datastore entity must be a struct")
	}

	columns, e := parsePublicFieldsToSqliteColumns(t)
	if e != nil {
		return empty, e
	}

	table := DbTable{
		name:    t.Name(),
		columns: columns,
	}

	return table, nil
}

func isPublicField(field reflect.StructField) bool {
	firstLetter := []rune(field.Name)[0]
	return unicode.IsUpper(firstLetter)
}

func parsePublicFieldsToSqliteColumns(
	t reflect.Type,
) ([]DbColumn, error) {
	var columns []DbColumn

	for field := range t.Fields() {
		if !isPublicField(field) {
			continue
		}

		col := DbColumn{
			name:    field.Name,
			sqlType: mapStructFieldTypeToSqliteType(field),
		}

		e := validateDbColumn(t, field, col)
		if e != nil {
			return nil, e
		}

		columns = append(columns, col)
	}

	return columns, nil
}

func mapStructFieldTypeToSqliteType(
	field reflect.StructField,
) string {
	switch field.Type.Kind() {
	case reflect.String:
		return "TEXT"
	case reflect.Int:
		return "INTEGER"
	case reflect.Float64:
		return "REAL"
	default:
		return ""
	}
}

func validateDbColumn(
	object reflect.Type,
	field reflect.StructField,
	col DbColumn,
) error {
	if col.sqlType == "" {
		return fmt.Errorf(
			"Datastore struct '%s' has field '%s' with unsupported type '%s'",
			object.Name(),
			field.Name,
			field.Type.Kind(),
		)
	}

	return nil
}
