package sqlick

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
)

var (
	// ErrNotStruct is used when a non-struct object is
	// passed when a struct is required.
	ErrNotStruct = errors.New("Entity must be a struct")

	// ErrBadFieldType is used when a struct field's type
	// has no mappable SQL type, i.e. cannot be used.
	ErrBadFieldType = errors.New(
		"Struct has unsupported field type",
	)

	// ErrMissingFields is used when a struct has no public
	// fields. Every table must have at least one column.
	ErrMissingFields = errors.New(
		"Struct must have at least one public field",
	)
)

func Parse(object any) (SqlickTable, error) {
	empty := SqlickTable{}
	table := SqlickTable{}

	structType := reflect.TypeOf(object)

	if structType.Kind() != reflect.Struct {
		return empty, ErrNotStruct
	}

	table.GoName = structType.Name()
	e := parseColumns(&table, structType)

	if e != nil {
		return empty, fmt.Errorf(
			"Parse error with struct/table '%s': %w",
			structType.Name(),
			e,
		)
	}

	if len(table.Columns) == 0 {
		return empty, fmt.Errorf(
			"Parse error with struct/table '%s': %w",
			structType.Name(),
			ErrMissingFields,
		)
	}

	return table, nil
}

func parseColumns(
	table *SqlickTable,
	structType reflect.Type,
) error {
	for field := range structType.Fields() {
		if !isPublicField(field) {
			continue
		}

		if e := parseColumn(table, field); e != nil {
			return fmt.Errorf(
				"Parse error with field/column '%s': %w",
				field.Name,
				e,
			)
		}
	}

	return nil
}

func isPublicField(field reflect.StructField) bool {
	firstLetter := []rune(field.Name)[0]
	return unicode.IsUpper(firstLetter)
}

func parseColumn(
	table *SqlickTable,
	field reflect.StructField,
) error {
	col := SqlickColumn{
		GoName: field.Name,
	}

	e := parseColumnType(&col, field)
	if e != nil {
		return e
	}

	table.Columns = append(table.Columns, col)
	return nil
}

func parseColumnType(
	col *SqlickColumn,
	field reflect.StructField,
) error {
	switch field.Type.Kind() {
	case reflect.String:
		col.GoType = "string"
	case reflect.Int:
		col.GoType = "int"
	case reflect.Float64:
		col.GoType = "float64"
	default:
		return ErrBadFieldType
	}

	return nil
}
