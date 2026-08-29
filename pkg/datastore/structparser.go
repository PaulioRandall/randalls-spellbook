package datastore

import (
	"fmt"
	"reflect"
	"unicode"
)

func parseDbTable(object any) (DbTable, error) {
	empty := DbTable{}
	t := reflect.TypeOf(object)

	if t.Kind() != reflect.Struct {
		return empty, fmt.Errorf("Datastore entity must be a struct")
	}

	columns, e := parsePublicStructFieldsAsDbColumns(t)
	if e != nil {
		return empty, e
	}

	table := DbTable{
		Name:    t.Name(),
		Columns: columns,
	}

	return table, nil
}

func isPublicField(field reflect.StructField) bool {
	firstLetter := []rune(field.Name)[0]
	return unicode.IsUpper(firstLetter)
}

func parsePublicStructFieldsAsDbColumns(
	structType reflect.Type,
) ([]DbColumn, error) {
	var columns []DbColumn

	for field := range structType.Fields() {
		if !isPublicField(field) {
			continue
		}

		col := DbColumn{
			Name:     field.Name,
			DataType: mapStructFieldTypeToString(field),
		}

		e := validateDbColumn(structType, field, col)
		if e != nil {
			return nil, e
		}

		columns = append(columns, col)
	}

	return columns, nil
}

func mapStructFieldTypeToString(
	field reflect.StructField,
) string {
	switch field.Type.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int:
		return "int"
	case reflect.Float64:
		return "float64"
	default:
		return ""
	}
}

func validateDbColumn(
	structType reflect.Type,
	field reflect.StructField,
	col DbColumn,
) error {
	if col.DataType == "" {
		return fmt.Errorf(
			"Datastore struct '%s' has field '%s' with unsupported type '%s'",
			structType.Name(),
			field.Name,
			field.Type.Kind(),
		)
	}

	return nil
}
