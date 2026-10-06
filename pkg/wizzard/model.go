package wizzard

import (
	"fmt"
	"reflect"
	"strings"
)

// Model maps a Go struct, or part of it, to a database
// table structure, or part of it. A Model is always
// derived from a struct, so there's a one-to-one
// relationship between a Go type and a model but
// one-to-many relationship between tables and models and
// Go types. However, it's possible for some models to
// differ by SqlName only.
type Model struct {
	// GoType is the struct type as returned by
	// reflect.TypeOf.
	GoType reflect.Type

	// GoName is the struct's name, i.e. GoType.Name().
	GoName string

	// SqlName is the name of the table the struct represents
	// within a database.
	SqlName string

	// Props is an ordered list of exported field-to-column
	// mappings.
	Props []Property
}

// Property maps a Go struct's field to a database table
// column.
type Property struct {
	// GoType is the field type as returned by
	// reflect.TypeOf.
	GoType reflect.Type

	// GoIndex is the field index, i.e. it's 0-based position
	// within the struct definition.
	GoIndex int

	// GoName is the field's name, i.e. GoType.Name().
	GoName string

	// SqlType is the SQLite type that maps to the GoType.
	// Note that several Go types may map to a single SQLite
	// type.
	SqlType string

	// SqlName is the name of the column within a database.
	SqlName string

	// SqlDefault is the default value given to the column if
	// a value is not provided during row insertion.
	// This is always the zero value of the GoType.
	SqlDefault any

	// IsKey is true if the property represents an ID or
	// primary key for the model. By default this is the
	// first property in the model, however, this may differ
	// if adjusting a property to align with an existing
	// table column as with the [Map] function.
	IsKey bool
}

// Parse creates a complete [Model] of the passed object's
// type. One should not assume the result will map directly
// to a table within a database; use [Map] to create a
// model that is usable with a specific database. Type
// mappings:
//
//	INTEGER:
//		int, int8, int16, int32, int64,
//		uint, uint8, uint16, uint32, uint64
//	REAL:
//		float32, float64
//	TEXT:
//		string
//
// You are responsible for choosing appropriate data types
// for your structs. If you choose to use an unsigned
// integer type (unit etc) or shorter integer type
// (int8 etc) then you are responsible for ensuring or
// managing number polarity and overflow. I recommend
// sticking to int, int64, float64, and string as these are
// the least likely to cause problems.
func Parse(object any) (Model, error) {
	objectType := derefObjectType(object)
	return parseTypeAsModel(objectType.Name(), objectType)
}

// ParseAs is the same as [Parse] except the table name is
// provided explicitly.
func ParseAs(table string, object any) (Model, error) {
	objectType := derefObjectType(object)
	return parseTypeAsModel(table, objectType)
}

func parseTypeAsModel(
	table string,
	objectType reflect.Type,
) (model Model, e error) {
	if objectType.Kind() != reflect.Struct {
		e = ErrNotStruct.Fmt(objectType.Kind().String())
		goto Err
	}

	model.GoType = objectType
	model.GoName = objectType.Name()
	model.SqlName = table

	model.Props, e = parseProps(objectType)
	if e != nil {
		goto Err
	}

	return model, nil

Err:
	return Model{}, ErrForType.
		Fmt(objectType.Name()).
		Wrap(e).
		WrapIn(ErrForTable).
		Fmt(table)
}

func derefObjectType(object any) reflect.Type {
	typ := reflect.TypeOf(object)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

func parseProps(objectType reflect.Type) ([]Property, error) {
	isIdField := true

	var props []Property

	for i := 0; i < objectType.NumField(); i++ {
		f := objectType.Field(i)

		if !f.IsExported() {
			continue
		}

		prop, e := parseProp(f, i, isIdField)
		if e != nil {
			return nil, ErrForField.Fmt(f.Name).Wrap(e)
		}

		props = append(props, prop)
		isIdField = false
	}

	return props, nil
}

func parseProp(
	f reflect.StructField,
	i int,
	isIdField bool,
) (Property, error) {
	var prop Property
	var defaultValue any

	if f.Type == reflect.TypeOf("") {
		defaultValue = "''"
	} else {
		defaultValue = reflect.Zero(f.Type).Interface()
	}

	sqlType, e := mapGoToSqlType(f.Type.Kind())
	if e != nil {
		return Property{}, e
	}

	prop.GoType = f.Type
	prop.GoIndex = i
	prop.GoName = f.Name
	prop.SqlType = sqlType
	prop.SqlName = f.Name
	prop.SqlDefault = defaultValue
	prop.IsKey = isIdField

	return prop, nil
}

func mapGoToSqlType(fieldKind reflect.Kind) (string, error) {
	switch fieldKind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return "INTEGER", nil
	case reflect.Float32, reflect.Float64:
		return "REAL", nil
	case reflect.String:
		return "TEXT", nil
	default:
		return "", ErrUnsupportedType.Fmt(fieldKind.String())
	}
}

// String returns a developer friendly representation of
// the model.
func (m Model) String() string {
	return m.GoTypeString() + "\n" + m.SqlTableString()
}

// GoTypeString returns a string representation of the
// model's Go type.
func (m Model) GoTypeString() string {
	sb := strings.Builder{}

	sb.WriteString("Go Type: ")
	sb.WriteString(m.GoName)

	for _, p := range m.Props {
		s := fmt.Sprintf(
			"\n\t[%d] %s %s",
			p.GoIndex,
			p.GoName,
			p.GoType.Name(),
		)
		sb.WriteString(s)
	}

	return sb.String()
}

// SqlTableString returns a string representation of the
// model's table details.
func (m Model) SqlTableString() string {
	sb := strings.Builder{}

	sb.WriteString("Sql Table: ")
	sb.WriteString(m.SqlName)

	for _, p := range m.Props {
		s := fmt.Sprintf(
			"\n\t%s %s",
			p.SqlName,
			p.SqlType,
		)
		sb.WriteString(s)
	}

	return sb.String()
}

// Represents returns true if the model represents the
// type of the passed object, and thus can be used with
// the model's methods without causing a type mismatch
// error.
func (m Model) Represents(object any) bool {
	return m.GoType == reflect.TypeOf(object)
}

// KeyProp returns the property representing the model's
// key. Zero value is returned if not present.
func (m Model) KeyProp() Property {
	for _, prop := range m.Props {
		if prop.IsKey {
			return prop
		}
	}
	return Property{}
}

// NonKeyProps returns all properties except the key
// property.
func (m Model) NonKeyProps() []Property {
	result := make([]Property, 0, len(m.Props))

	for _, prop := range m.Props {
		if !prop.IsKey {
			result = append(result, prop)
		}
	}

	return result
}

// New creates a instance of the columns GoType. It will
// contain the type's zero value. The value returned is
// explicitly cast to T.
func (p Property) New[T any]() T {
	return reflect.New(p.GoType).Interface().(T)
}
